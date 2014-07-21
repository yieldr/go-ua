package ua

import (
	"regexp"
	"strings"
)

type Device struct {
	Family string
}

type DevicePattern struct {
	regexp      *regexp.Regexp
	Regex       string `yaml:"regex"`
	Replacement string `yaml:"device_replacement"`
}

func (d *DevicePattern) Regexp() *regexp.Regexp {
	if d.regexp == nil {
		d.regexp = regexp.MustCompile(d.Regex)
	}
	return d.regexp
}

func (d *DevicePattern) HasSubexp() bool {
	return d.Regexp().NumSubexp() >= 1
}

func (d *DevicePattern) Match(line string) (*Device, bool) {
	matches := d.Regexp().FindStringSubmatch(line)
	if len(matches) > 0 {
		device := new(Device)
		if d.Replacement != "" {
			if strings.Contains(d.Replacement, "$1") && d.HasSubexp() && len(matches) >= 2 {
				device.Family = strings.Replace(d.Replacement, "$1", matches[1], 1)
			} else {
				device.Family = d.Replacement
			}
		} else if d.HasSubexp() {
			device.Family = matches[1]
		}
		return device, true
	}
	return nil, false
}

func (d *Device) String() string {
	return d.Family
}

func unknownDevice() *Device {
	d := new(Device)
	d.Family = "Other"
	return d
}
