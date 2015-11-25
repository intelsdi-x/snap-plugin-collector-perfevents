# Snap Perfevents Collector Plugin

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

Project link: https://github.com/intelsdi-x/snap-plugin-collector-perfevents

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license-and-authors)
6. [Acknowledgements](#acknowledgements)

## Getting Started
In order to use this plugin you need "perf" to be installed on a Linux target host.

### System Requirements

* "perf" installed on a host
* Linux kernel version at least 2.6.31
* /proc/sys/kernel/perf_event_paranoid set to 0

### Installation

Plugin compilation
```
make
```

### Configuration and Usage

* root previledges are required in order to run this plugin
* this plugin was tested on Ubuntu 14.04

## Documentation
To learn more about metrics exposed by "perf" visit Perf wiki at: https://perf.wiki.kernel.org/index.php/Main_Page

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

## Community Support
This repository is one of **many** plugins in the **Snap Framework**: a powerful telemetry agent framework. To reach out to other uses, reach out to us on:

* Snap Gitter channel (@TODO Link)
* Our Google Group (@TODO Link)

The full project is at http://github.com:intelsdi-x/snap.

## Contributing
We love contributions! :heart_eyes:

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Snap, along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
List authors, co-authors and anyone you'd like to mention

* Author: [Andrzej Kuriata](https://github.com/andrzej-k)
* Author: [Justin Guidroz](https://github.com/geauxvirtual)

**Thank you!** Your contribution is incredibly important to us.
