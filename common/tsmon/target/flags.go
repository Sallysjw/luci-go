// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package target

import (
	"flag"
	"net"
	"os"
	"regexp"
	"strings"
)

// SysInfo overrides system's hostname and region for tests.
type SysInfo struct {
	Hostname string
	Region   string
}

// Flags defines command line flags related to tsmon targets.  Use NewFlags()
// to get a Flags struct with sensible default values.
type Flags struct {
	TargetType      string
	DeviceHostname  string
	DeviceRegion    string
	DeviceRole      string
	DeviceNetwork   string
	TaskServiceName string
	TaskJobName     string
	TaskRegion      string
	TaskHostname    string
	TaskNumber      int
	AutoGenHostname bool

	// If nil, system info is computed from the actual host. Used
	// in tests.
	SysInfo *SysInfo
}

// NewFlags returns a Flags struct with sensible default values.  Hostname,
// region and network flags are expensive to compute, so get assigned default
// values later in SetDefaultsFromHostname.
func NewFlags() Flags {
	return Flags{
		TargetType:      "device",
		DeviceHostname:  "",
		DeviceRegion:    "",
		DeviceRole:      "default",
		DeviceNetwork:   "",
		TaskServiceName: "",
		TaskJobName:     "",
		TaskRegion:      "",
		TaskHostname:    "",
		TaskNumber:      0,
		AutoGenHostname: false,
	}
}

// SetDefaultsFromHostname computes the expensive default values for hostname,
// region and network fields.
func (fl *Flags) SetDefaultsFromHostname() {
	if fl.SysInfo == nil {
		hostname, region := getFQDN()
		fl.SysInfo = &SysInfo{Hostname: hostname, Region: region}
	}
	network := getNetwork(fl.SysInfo.Hostname)
	hostname := fl.SysInfo.Hostname

	if fl.AutoGenHostname {
		hostname = "autogen:" + hostname
	}
	if fl.DeviceHostname == "" {
		fl.DeviceHostname = hostname
	}
	if fl.DeviceRegion == "" {
		fl.DeviceRegion = fl.SysInfo.Region
	}
	if fl.DeviceNetwork == "" {
		fl.DeviceNetwork = network
	}
	if fl.TaskRegion == "" {
		fl.TaskRegion = fl.SysInfo.Region
	}
	if fl.TaskHostname == "" {
		fl.TaskHostname = hostname
	}
}

// Register adds tsmon target related flags to a FlagSet.
func (fl *Flags) Register(f *flag.FlagSet) {
	f.StringVar(&fl.TargetType, "ts-mon-target-type", fl.TargetType,
		"the type of target that is being monitored ('device' or 'task')")
	f.StringVar(&fl.DeviceHostname, "ts-mon-device-hostname", fl.DeviceHostname,
		"name of this device")
	f.StringVar(&fl.DeviceRegion, "ts-mon-device-region", fl.DeviceRegion,
		"name of the region this devices lives in")
	f.StringVar(&fl.DeviceRole, "ts-mon-device-role", fl.DeviceRole,
		"role of the device")
	f.StringVar(&fl.DeviceNetwork, "ts-mon-device-network", fl.DeviceNetwork,
		"name of the network this device is connected to")
	f.StringVar(&fl.TaskServiceName, "ts-mon-task-service-name", fl.TaskServiceName,
		"name of the service being monitored")
	f.StringVar(&fl.TaskJobName, "ts-mon-task-job-name", fl.TaskJobName,
		"name of this job instance of the task")
	f.StringVar(&fl.TaskRegion, "ts-mon-task-region", fl.TaskRegion,
		"name of the region in which this task is running")
	f.StringVar(&fl.TaskHostname, "ts-mon-task-hostname", fl.TaskHostname,
		"name of the host on which this task is running")
	f.IntVar(&fl.TaskNumber, "ts-mon-task-number", fl.TaskNumber,
		"number (e.g. for replication) of this instance of this task")
	f.BoolVar(&fl.AutoGenHostname, "ts-mon-autogen-hostname", fl.AutoGenHostname,
		"Indicate that the hostname is autogenerated. "+
			"This option must be set on autoscaled GCE VMs, Kubernetes pods, "+
			"or any other hosts with dynamically generated names.")
}

func getFQDN() (string, string) {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if names, err := net.LookupAddr(ipNet.IP.String()); err == nil {
					for _, name := range names {
						parts := strings.Split(name, ".")
						if len(parts) > 1 {
							return strings.ToLower(parts[0]), strings.ToLower(parts[1])
						}
					}
				}
			}
		}
	}
	if hostname, err := os.Hostname(); err == nil {
		parts := strings.Split(hostname, ".")
		if len(parts) > 1 {
			return strings.ToLower(parts[0]), strings.ToLower(parts[1])
		}
		return strings.ToLower(hostname), "unknown"
	}
	return "unknown", "unknown"
}

func getNetwork(hostname string) string {
	if m := regexp.MustCompile(`^([\w-]*?-[acm]|master)(\d+)a?$`).FindStringSubmatch(hostname); m != nil {
		return m[2]
	}
	return ""
}
