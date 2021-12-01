package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	exporterVersion = "0.5"
)

// SMARTOptions is a inner representation of a options
type SMARTOptions struct {
	BindTo                string
	URLPath               string
	FakeJSON              bool
	SMARTctlLocation      string
	CollectPeriod         string
	CollectPeriodDuration time.Duration
	Devices               []string
}

// Options is a representation of a options
type Options struct {
	SMARTctl SMARTOptions
}

// Parse options from command line arguments file
func loadOptions() *Options {
	opts := &Options{}

	flag.StringVar(&opts.SMARTctl.BindTo, "bindTo", ":9633", "address and port to bind to")
	flag.StringVar(&opts.SMARTctl.URLPath, "urlPath", "/metrics", "metrics endpoint path")
	flag.BoolVar(&opts.SMARTctl.FakeJSON, "fakeJson", false, "use fake json (only for debugging)")
	flag.StringVar(&opts.SMARTctl.SMARTctlLocation, "smartCtlLocation", "/usr/sbin/smartctl", "smartctl binary version >7.0 required")
	flag.StringVar(&opts.SMARTctl.CollectPeriod, "collectPeriod", "60s", "minimal time interval between two smartctl runs")
	verbose := flag.Bool("verbose", false, "Verbose log output")
	debug := flag.Bool("debug", false, "Debug log output")
	version := flag.Bool("version", false, "Show application version and exit")
	flag.Parse()
	deviceGlobs := flag.Args()

	if *version {
		fmt.Printf("smartctl_exporter version: %s\n", exporterVersion)
		os.Exit(0)
	}

	logger = newLogger(*verbose, *debug)
	d, err := time.ParseDuration(opts.SMARTctl.CollectPeriod)
	if err != nil {
		logger.Panic("Failed read collect_not_more_than_period (%s): %s", opts.SMARTctl.CollectPeriod, err)
	}
	opts.SMARTctl.CollectPeriodDuration = d

	if len(deviceGlobs) > 0 {
		opts.SMARTctl.Devices = make([]string, 0)
		for _, deviceGlob := range deviceGlobs {
			glob, err := filepath.Glob(deviceGlob)
			if err != nil {
				logger.Panic("error reading device list: %s", err)
			}
			opts.SMARTctl.Devices = append(opts.SMARTctl.Devices, glob...)
		}
		logger.Debug("Parsed options: %s", opts)
	} else {
		opts.SMARTctl.Devices, err = findSmartDevices(opts)
		if err != nil {
			logger.Panic("scanning for S.M.A.R.T enabled devices failed: %s", err)
		}
	}
	if len(opts.SMARTctl.Devices) == 0 {
		logger.Panic("at least one device is required")
	}
	return opts
}

// use smartctl --scan to find disks
func findSmartDevices(opts *Options) ([]string, error) {
	logger.Info("search for smart enabled devices...")
	out, err := exec.Command(opts.SMARTctl.SMARTctlLocation, "--json", "--scan").Output()
	if err != nil {
		logger.Warning("S.M.A.R.T. output reading error: %s", err)
		return nil, err
	}
	json := parseJSON(string(out))
	rcOk := resultCodeIsOk(json.Get("smartctl.exit_status").Int())
	jsonOk := jsonIsOk(json)

	if !jsonOk || !rcOk {
		return nil, fmt.Errorf("error parsing S.M.A.R.T. output")
	}

	devices := make([]string, 0)
	for _, dev := range json.Get("devices").Array() {
		devName := dev.Get("name").String()
		protocol := dev.Get("protocol").String()
		logger.Info("...found device %s using protocol %s", devName, protocol)
		devices = append(devices, devName)
	}
	return devices, nil
}
