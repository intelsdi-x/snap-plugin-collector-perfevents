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
* Linux kernel version at least 2.6.31
* /proc/sys/kernel/perf_event_paranoid set to 0
* [golang 1.4+](https://golang.org/dl/)

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

Namespace | Data Type | Source | Description
------------------|--------|-----------|-------------------------------
/intel/linux/perfevents/cgroup/cycles/[GROUP_NAME] | float64 | hostname | Total cycles. Be wary of what happens during CPU frequency scaling.
/intel/linux/perfevents/cgroup/instructions/[GROUP_NAME] | float64 | hostname | Retired instructions
/intel/linux/perfevents/cgroup/cache-references/[GROUP_NAME] | float64 | hostname | Cache accesses. Usually this indicates Last Level Cache accesses but this may vary depending on your CPU. This may include prefetches and coherency messages; again this depends on the design of your CPU.
/intel/linux/perfevents/cgroup/cache-misses/[GROUP_NAME] | float64 | hostname | Cache misses. Usually this indicates Last Level Cache misses; this is intended to be used in conjunction with the cache-references event to calculate cache miss rates.
/intel/linux/perfevents/cgroup/branch-instructions/[GROUP_NAME] | float64 | hostname | Retired branch instructions
/intel/linux/perfevents/cgroup/branch-misses/[GROUP_NAME] | float64 | hostname | Mispredicted branch instructions
/intel/linux/perfevents/cgroup/stalled-cycles-frontend/[GROUP_NAME] | float64 | hostname | Stalled cycles during issue
/intel/linux/perfevents/cgroup/stalled-cycles-backend/[GROUP_NAME] | float64 | hostname | Stalled cycles during retirement
/intel/linux/perfevents/cgroup/ref-cycles/[GROUP_NAME] | float64 | hostname | Total cycles; not affected by CPU frequency scaling

### Examples
TBD

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
