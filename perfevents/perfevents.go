/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package perfevents

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
)

const (
	// Name of plugin
	name = "perfevents"
	// Version of plugin
	version = 8
	// Type of plugin
	pluginType = plugin.CollectorPluginType
	// Namespace definition
	ns_vendor     = "intel"
	ns_class      = "linux"
	ns_type       = "perfevents"
	ns_subtype    = "cgroup"
	not_supported = "<not supported>"
	not_counted   = "<not counted>"
)

type event struct {
	etype string
	id    string
	value interface{}
}

type Perfevents struct {
	cgroup_events map[string]event
	Init          func() error

	// the map of perf events which are unsupported by kernel
	// <not supported> - platform does not support this processors's performance monitoring unit (PMU) when kernel does not support perf event, but the kernel has perf module enabled
	// <not counted> - when there is no kernel support for perf event (disabled module)
	unsupportedEvents map[string]string
}

var CGROUP_EVENTS = []string{"cycles", "instructions", "cache-references", "cache-misses",
	"branch-instructions", "branch-misses", "stalled-cycles-frontend",
	"stalled-cycles-backend", "ref-cycles"}

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

// CollectMetrics returns HW metrics from perf events subsystem
// for Cgroups present on the host.
func (p *Perfevents) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	if len(mts) == 0 {
		return nil, nil
	}

	events := []string{}
	cgroups := []string{}

	// Get list of events and cgroups from Namespace
	for _, m := range mts {
		event, cgroup, err := getEventAndCgroupFromNamespace(m.Namespace().Strings())
		if err != nil {
			return nil, err
		}

		if _, isNotSupported := p.unsupportedEvents[event]; !isNotSupported {
			// append if supported (not exist in map of unsupported events)
			events = append(events, event)
			cgroups = append(cgroups, cgroup)
		}
	}

	if len(cgroups) != len(events) {
		return nil, fmt.Errorf("Invalid args for perf command, the number of events=%d, cgroups%d (expected to be equal)", len(events), len(cgroups))
	}

	// in case when all requested metrics have an event which is unsupported
	if len(events) == 0 {
		return nil, fmt.Errorf("There is no supported perf events for requested metrics")
	}

	// Prepare events (-e) and Cgroups (-G) switches for "perf stat"
	cgroupsSwitch := "-G" + strings.Join(cgroups, ",")
	eventsSwitch := "-e" + strings.Join(events, ",")

	// Prepare "perf stat" command
	cmd := exec.Command("perf", "stat", "--log-fd", "1", `-x;`, "-a", eventsSwitch, cgroupsSwitch, "--", "sleep", "1")
	output, err := cmd.Output()

	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Execution of perf command returns err=%v, command=%v", err, cmd.Args))
		return nil, err
	}

	// Parse "perf stat" output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// skip empty lines
		if len(line) == 0 {
			continue
		}

		data := strings.Split(line, ";")
		if len(data) < 3 {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Invalid output format %v (expected at least 3 elements separated by `;`)", line))
			continue
		}
		e := event{id: data[3], etype: data[2]}
		ekey := getEventKey(e.etype, e.id)

		// check kernel support for perf event
		switch data[0] {
		case not_supported:
			fmt.Fprintln(os.Stderr, fmt.Sprintf("There is no support for perf event `%s`", e.etype))
			// add this event to unsupported_events; it will be omitted in collection next time
			p.unsupportedEvents[e.etype] = data[0]
			// set value to "<not supported>"
			e.value = data[0]

		case not_counted:
			// only log it, the value of event is `nil`
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Perf event %s is not counted for %s", e.etype, e.id))

		default:
			// numeric value is expected
			if e.value, err = strconv.ParseUint(data[0], 10, 64); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Sprintf("Invalid metric value for %s:%s, err=%v", e.etype, e.id, err))
			}
		}
		p.cgroup_events[ekey] = e
	}

	// Populate metrics
	metrics := []plugin.MetricType{}

	for _, m := range mts {
		var val interface{}
		// skip error (because it's handle in the beginning of CollectMetrics())
		event, cgroup, _ := getEventAndCgroupFromNamespace(m.Namespace().Strings())

		// retrieve value based on eventKey
		if event, ok := p.cgroup_events[getEventKey(event, cgroup)]; ok {
			val = event.value
		}

		metric := plugin.MetricType{
			Namespace_: m.Namespace(),
			Data_:      val,
			Timestamp_: time.Now(),
			Tags_:      m.Tags(),
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// GetMetricTypes returns the metric types exposed by perf events subsystem
func (p *Perfevents) GetMetricTypes(_ plugin.ConfigType) ([]plugin.MetricType, error) {
	err := p.Init()
	if err != nil {
		return nil, err
	}
	cgroups, err := list_cgroups()
	if err != nil {
		return nil, err
	}
	if len(cgroups) == 0 {
		return nil, nil
	}

	return get_supported_metrics(ns_subtype, cgroups, CGROUP_EVENTS), nil
}

// GetConfigPolicy returns a ConfigPolicy
func (p *Perfevents) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	return c, nil
}

// New initializes Perfevents plugin
func NewPerfeventsCollector() *Perfevents {
	return &Perfevents{cgroup_events: make(map[string]event), unsupportedEvents: make(map[string]string), Init: initialize}
}

func initialize() error {
	file, err := os.Open("/proc/sys/kernel/perf_event_paranoid")
	if err != nil {
		if os.IsExist(err) {
			return errors.New("perf_event_paranoid file exists but couldn't be opened")
		}
		return errors.New("perf event system not enabled")
	}

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return errors.New("cannot read from perf_event_paranoid")
	}

	paranoid, err := strconv.ParseInt(scanner.Text(), 10, 64)
	if err != nil {
		return errors.New("invalid value in perf_event_paranoid file")
	}

	if paranoid >= 1 {
		fmt.Fprintf(os.Stderr, "Per event paranoia level is %v (see `/proc/sys/kernel.perf_event_paranoid`). There is no permission to collect some stats. List of perf metrics can be limited", paranoid)
	}

	return nil
}

func getEventAndCgroupFromNamespace(ns []string) (event string, cgroup string, err error) {
	if err = validateNamespace(ns); err != nil {
		return
	}
	flatcgroup := strings.Replace(ns[5], "_", ".", -1)
	cgroup = strings.Replace(flatcgroup, ":", "/", -1)
	event = ns[4]
	return
}

func getEventKey(etype, eid string) string {
	return fmt.Sprintf("%s:%s", etype, eid)
}

func get_supported_metrics(source string, cgroups []string, events []string) []plugin.MetricType {
	mts := []plugin.MetricType{}
	for _, e := range events {
		for _, c := range flatten_cg_name(cgroups) {
			mts = append(mts, plugin.MetricType{Namespace_: core.NewNamespace(ns_vendor, ns_class, ns_type, source, e, c)})
		}
	}
	return mts
}
func flatten_cg_name(cg []string) []string {
	flat_cg := []string{}
	for _, c := range cg {
		flat_cg = append(flat_cg, strings.Replace(c, "/", ":", -1))
	}
	return flat_cg
}

func list_cgroups() ([]string, error) {
	cgroups := []string{}
	base_path := "/sys/fs/cgroup/perf_event/"
	err := filepath.Walk(base_path, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			cgroup_name := strings.TrimPrefix(path, base_path)
			if len(cgroup_name) > 0 {
				cgroups = append(cgroups, removeNotAllowChars(cgroup_name))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cgroups, nil
}

func removeNotAllowChars(str string) string {
	notAllowedChars := []string{"(", ")", "[", "]", "{", "}", " ", ".", ",", ";", "?", "!"}

	for _, chr := range notAllowedChars {
		str = strings.Replace(str, chr, "_", -1)
	}
	return str
}
func validateNamespace(namespace []string) error {
	if len(namespace) != 6 {
		return errors.New(fmt.Sprintf("unknown metricType %s (should containt exactly 6 segments)", namespace))
	}
	if namespace[0] != ns_vendor {
		return errors.New(fmt.Sprintf("unknown metricType %s (expected 1st segment %s)", namespace, ns_vendor))
	}

	if namespace[1] != ns_class {
		return errors.New(fmt.Sprintf("unknown metricType %s (expected 2nd segment %s)", namespace, ns_class))
	}
	if namespace[2] != ns_type {
		return errors.New(fmt.Sprintf("unknown metricType %s (expected 3rd segment %s)", namespace, ns_type))
	}
	if namespace[3] != ns_subtype {
		return errors.New(fmt.Sprintf("unknown metricType %s (expected 4th segment %s)", namespace, ns_subtype))
	}
	if !namespaceContains(namespace[4], CGROUP_EVENTS) {
		return errors.New(fmt.Sprintf("unknown metricType %s (expected 5th segment %v)", namespace, CGROUP_EVENTS))
	}
	return nil
}

func namespaceContains(element string, slice []string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}
