# snap collector plugin - Linux perf events

This plugin collects following hardware metrics for Cgroups from "perf" (Performance Counters for Linux):
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

This plugin is used in the [snap framework] (http://github.com/intelsdi-x/snap).

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

* root privileges
* "perf" installed on a host
* /proc/sys/kernel/perf_event_paranoid set to 0 (in other case list of metrics can be limited)
* Linux kernel version at least 2.6.31
* [golang 1.5+](https://golang.org/dl/) (needed only for building)



### Installation

#### Download perfevents plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [GitHub Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-perfevents
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-perfevents.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported  
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`
 
## Documentation

To learn more about this plugin and Linux perf counters, visit:

* [Linux perf events wiki] (https://perf.wiki.kernel.org/index.php/Main_Page)
* [snap perfevents unit test](https://github.com/intelsdi-x/snap-plugin-collector-perfevents/blob/master/perfevents/perfevents_test.go)

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

This is done from the snap directory.

In one terminal window, open the snap daemon (in this case with logging set to 1 and trust disabled):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0
```

In another terminal window:
Load perfevents collector plugin
```
$ $SNAP_PATH/bin/snapctl plugin load build/rootfs/snap-plugin-collector-perfevents

Plugin loaded
Name: perfevents
Version: 8
Type: collector
Signed: false
Loaded Time: Wed, 13 Jul 2016 12:43:19 CEST
```
See available metrics for your system
```
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file (see examples/tasks/perfevents-file.json)[examples/tasks/perfevents-file.json]:    
```json
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "3s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/linux/perfevents/*" :{}
            },            
           "publish": [
                        {
                            "plugin_name": "file",
                            "config": {
                                "file": "/tmp/published_perfevents"
                            }
                        }
                ]
        }
    }
}
```

Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load build/plugin/snap-publisher-file
Plugin loaded
Name: file
Version: 3
Type: publisher
Signed: false
Loaded Time: Wed, 13 Jul 2016 12:44:11 CEST
```

Create task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/perfevents-file.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

See sample output from snapctl task watch <task_id>

```
$ $SNAP_PATH/bin/snapctl task watch 02dd7ff4-8106-47e9-8b86-70067cd0a850

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
(Keys `ctrl+c` terminate task watcher)

These data are published to file and stored there (in this example in /tmp/published_perfevents).

Notice, that if perf stat command returns:
 - **not counted**, the value of metric will be `nil`
 - **not supported**, the value of metric will be `<not supported>` and it will be omitted in the next cycle of collection

Stop task:
```
$ $SNAP_PATH/bin/snapctl task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-perfevents/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-perfevents/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [@andrzej-k](https://github.com/andrzej-k)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
