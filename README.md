
# DISCONTINUATION OF PROJECT 

**This project will no longer be maintained by Intel.  Intel will not provide or guarantee development of or support for this project, including but not limited to, maintenance, bug fixes, new releases or updates.  Patches to this project are no longer accepted by Intel. If you have an ongoing need to use this project, are interested in independently developing it, or would like to maintain patches for the community, please create your own fork of the project.**



# Snap collector plugin - Linux perf events

This plugin collects following hardware metrics for Cgroups from "perf" (Performance Counters for Linux) for [Snap Telemetry Framework](http://github.com/intelsdi-x/snap):
*  cycles
*  instructions
*  cache-references
*  cache-misses
*  branch-instructions
*  branch-misses
*  stalled-cycles-frontend
*  stalled-cycles-backend
*  ref-cycles

To see which of them are supported on your platform, please execute command `perf stat ls`.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

In order to use this plugin you need "perf" to be installed on a Linux target host.  

### System Requirements

- Linux (kernel 2.6.31+)
- root privileges
- "perf" installed on a host
- /proc/sys/kernel/perf_event_paranoid set to 0 (in other case list of metrics can be limited)

### Installation

#### Download perfevents plugin binary:
You can get the pre-built binaries for your OS and architecture at plugins's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-collector-perfevents/releases) page. Download the plugin from the latest release and load it into snapteld (`/opt/snap/plugins` is the default location for Snap packages).

#### To build the plugin binary:
Use https://github.com/intelsdi-x/snap-plugin-collector-perfevents or your fork as repo.

Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-perfevents.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage

* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-perfevents#examples).
 
## Documentation

To learn more about this plugin and Linux perf counters, visit:

* [Linux perf events wiki](https://perf.wiki.kernel.org/index.php/Main_Page)
* [Snap perfevents unit test](https://github.com/intelsdi-x/snap-plugin-collector-perfevents/blob/master/perfevents/perfevents_test.go)

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type |  Description
------------------|--------|-------------------------------
/intel/linux/perfevents/cgroup/cycles/[GROUP_NAME] | float64  | Total cycles. Be wary of what happens during CPU frequency scaling.
/intel/linux/perfevents/cgroup/instructions/[GROUP_NAME] | float64  | Retired instructions
/intel/linux/perfevents/cgroup/cache-references/[GROUP_NAME] | float64  | Cache accesses. Usually this indicates Last Level Cache accesses but this may vary depending on your CPU. This may include prefetches and coherency messages; again this depends on the design of your CPU.
/intel/linux/perfevents/cgroup/cache-misses/[GROUP_NAME] | float64  | Cache misses. Usually this indicates Last Level Cache misses; this is intended to be used in conjunction with the cache-references event to calculate cache miss rates.
/intel/linux/perfevents/cgroup/branch-instructions/[GROUP_NAME] | float64  | Retired branch instructions
/intel/linux/perfevents/cgroup/branch-misses/[GROUP_NAME] | float64  | Mispredicted branch instructions
/intel/linux/perfevents/cgroup/stalled-cycles-frontend/[GROUP_NAME] | float64  | Stalled cycles during issue
/intel/linux/perfevents/cgroup/stalled-cycles-backend/[GROUP_NAME] | float64  | Stalled cycles during retirement
/intel/linux/perfevents/cgroup/ref-cycles/[GROUP_NAME] | float64  | Total cycles; not affected by CPU frequency scaling

### Examples
Example running perfevents collector and writing data to a file.

Ensure [snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `sudo service snap-telemetry start`
* systemd: `sudo systemctl start snap-telemetry`
* command line: `sudo snapteld -l 1 -t 0 &`

Download and load snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-perfevents/latest/linux/x86_64/snap-plugin-collector-perfevents
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ snaptel plugin load snap-plugin-collector-perfevents
$ snaptel plugin load snap-plugin-publisher-file
```

See available metrics for your system
```
$ snaptel metric list
NAMESPACE                                                               VERSIONS
/intel/linux/perfevents/cgroup/branch-instructions/system_slice         8
/intel/linux/perfevents/cgroup/branch-misses/system_slice               8
/intel/linux/perfevents/cgroup/branch-misses/user_slice                 8
/intel/linux/perfevents/cgroup/cache-misses/system_slice                8
/intel/linux/perfevents/cgroup/cache-misses/user_slice                  8
/intel/linux/perfevents/cgroup/branch-instructions/system_slice         8
/intel/linux/perfevents/cgroup/branch-instructions/user_slice           8
/intel/linux/perfevents/cgroup/branch-misses/system_slice               8
/intel/linux/perfevents/cgroup/branch-misses/user_slice                 8
/intel/linux/perfevents/cgroup/cache-misses/system_slice                8
/intel/linux/perfevents/cgroup/cache-misses/user_slice                  8
/intel/linux/perfevents/cgroup/stalled-cycles-backend/system_slice      8
/intel/linux/perfevents/cgroup/stalled-cycles-backend/user_slice        8
/intel/linux/perfevents/cgroup/stalled-cycles-frontend/system_slice     8
/intel/linux/perfevents/cgroup/stalled-cycles-frontend/user_slice       8
```

Download an [example task file](https://github.com/intelsdi-x/snap-plugin-collector-perfevents/blob/master/examples/tasks/perfevents-file.json) and load it:
```
$ curl -sfLO https://raw.githubusercontent.com/intelsdi-x/snap-plugin-collector-perfevents/master/examples/tasks/perfevents-file.json
$ snaptel task create -t perfevents-file.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

See realtime output from `snaptel task watch <task_id>` (CTRL+C to exit)
```
$ snaptel task watch 02dd7ff4-8106-47e9-8b86-70067cd0a850

/intel/linux/perfevents/cgroup/branch-instructions/system_slice                                          1.249881e+06            2016-07-13 12:43:32.273501321 +0200 CEST
/intel/linux/perfevents/cgroup/branch-misses/system_slice                                                62875                   2016-07-13 12:43:32.273557564 +0200 CEST
/intel/linux/perfevents/cgroup/branch-misses/user_slice                                                                          2016-07-13 12:43:32.273561431 +0200 CEST
/intel/linux/perfevents/cgroup/cache-misses/system_slice                                                 166413                  2016-07-13 12:43:32.27344248 +0200 CEST
/intel/linux/perfevents/cgroup/cache-misses/user_slice                                                                           2016-07-13 12:43:32.273446952 +0200 CEST
/intel/linux/perfevents/cgroup/branch-instructions/system_slice                                          1.460461e+06            2016-07-13 12:43:35.268157381 +0200 CEST
/intel/linux/perfevents/cgroup/branch-instructions/user_slice                                                                    2016-07-13 12:43:35.268162628 +0200 CEST
/intel/linux/perfevents/cgroup/branch-misses/system_slice                                                82953                   2016-07-13 12:43:35.268226596 +0200 CEST
/intel/linux/perfevents/cgroup/branch-misses/user_slice                                                                          2016-07-13 12:43:35.268232068 +0200 CEST
/intel/linux/perfevents/cgroup/cache-misses/system_slice                                                 224287                  2016-07-13 12:43:35.268059369 +0200 CEST
/intel/linux/perfevents/cgroup/cache-misses/user_slice                                                                           2016-07-13 12:43:35.268064862 +0200 CEST
/intel/linux/perfevents/cgroup/stalled-cycles-backend/system_slice                                       <not supported>         2016-07-13 12:43:35.268064862 +0200 CEST
/intel/linux/perfevents/cgroup/stalled-cycles-backend/user_slice                                         <not supported>         2016-07-13 12:43:35.268064862 +0200 CEST
/intel/linux/perfevents/cgroup/stalled-cycles-frontend/system_slice                                      <not supported>         2016-07-13 12:44:32.266638857 +0200 CEST
/intel/linux/perfevents/cgroup/stalled-cycles-frontend/user_slice                                        <not supported>         2016-07-13 12:44:32.266644089 +0200 CEST
```

This data is published to a file `/tmp/published_perfevents` per task specification.

Notice, that if perf stat command returns:
 - **not counted**, the value of metric will be `nil`
 - **not supported**, the value of metric will be `<not supported>` and it will be omitted in the next cycle of collection

Stop task:
```
$ snaptel task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Using inside docker container
Plugin can collect perf events inside the docker container. To do that, it is needed to run container in privileged mode and mount host's cgroupfs.

#### Example docker run command for CentOS:
```bash
docker run -ti --privileged -v /sys/fs/cgroup:/sys/fs/cgroup:ro centos:latest bash
```
**--privileged** is needed to allow collecting system-wide stats  
**-v /sys/fs/cgroup:/sys/fs/cgroup:ro** is needed for perf to be able to access host's cgroup hierarchy

Procfs is shared between host and docker container, so you can configure /proc/sys/kernel/perf_event_paranoid once on host side (or from inside the privileged container, result will be the same).

#### Running on Kubernetes
To obtain the same result in Kubernetes, you need to configure your pods the same way:
- For running pod in privileged mode: https://kubernetes.io/docs/user-guide/security-context/
- For mounting cgroupfs volume: https://kubernetes.io/docs/user-guide/volumes/

#### Detecting which cgroup corresponds to current container
To obtain cgroup of container from inside of that container, you need to check /proc/self/cgroup file.

The file contains current container mountpoints for each subsystem and to retrieve current container's perf_event subsystem cgroup mountpoint, you can run such command:

```bash
cat /proc/self/cgroup | grep -oP "perf_event:/\K.*"
```

Example result:
```
docker/ae05db810861f060d1732f4ae508973782ac68ff60dff644c9ff4953eb437627
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-perfevents/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-perfevents/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@andrzej-k](https://github.com/andrzej-k)

**Thank you!** Your contribution, through code and participation, is incredibly important to us.
