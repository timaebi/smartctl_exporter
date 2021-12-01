[![CircleCI](https://circleci.com/gh/Sheridan/smartctl_exporter.svg?style=svg)](https://circleci.com/gh/Sheridan/smartctl_exporter)

# smartctl_exporter
Export smartctl statistics to prometheus

Example output you can show in [EXAMPLE.md](EXAMPLE.md)

## Need more?
**If you need additional metrics - contact me :)**
**Create a feature request, describe the metric that you would like to have and attach exported from smartctl json file**

# Requirements
smartmontools >= 7.0, because export to json [released in 7.0](https://www.smartmontools.org/browser/tags/RELEASE_7_0/smartmontools/NEWS#L11)

# Usage

```
Usage of smartctl_exporter:
  -bindTo string
        address and port to bind to (default ":9633")
  -collectPeriod string
        minimal time interval between two smartctl runs (default "60s")
  -debug
        Debug log output
  -fakeJson
        use fake json (only for debugging)
  -smartCtlLocation string
        smartctl binary version >7.0 required (default "/usr/sbin/smartctl")
  -urlPath string
        metrics endpoint path (default "/metrics")
  -verbose
        Verbose log output
  -version
        Show application version and exit

```