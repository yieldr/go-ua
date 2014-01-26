package ua

import (
	"regexp"
)

type Os struct {
	Family     string
	Major      string
	Minor      string
	Patch      string
	PatchMinor string
}

type OsPattern struct {
	regexp          *regexp.Regexp
	Regex           string `yaml:"regex"`
	OsReplacement   string `yaml:"os_replacement"`
	OsV1Replacement string `yaml:"os_v1_replacement"`
	OsV2Replacement string `yaml:"os_v2_replacement"`
}

func (os *OsPattern) Regexp() *regexp.Regexp {
	if os.regexp == nil {
		os.regexp = regexp.MustCompile(os.Regex)
	}
	return os.regexp
}

func (osPattern *OsPattern) Match(line string, os *Os) {
	bytes := osPattern.Regexp().FindStringSubmatch(line)
	if len(bytes) > 0 {
		groupCount := osPattern.Regexp().NumSubexp()

		if len(osPattern.OsReplacement) > 0 {
			os.Family = osPattern.OsReplacement
		} else if groupCount >= 1 {
			os.Family = bytes[1]
		}

		if len(osPattern.OsV1Replacement) > 0 {
			os.Major = osPattern.OsV1Replacement
		} else if groupCount >= 2 {
			os.Major = bytes[2]
			if len(osPattern.OsV2Replacement) > 0 {
				os.Minor = osPattern.OsV2Replacement
			} else if groupCount >= 3 {
				os.Minor = bytes[3]
				if groupCount >= 4 {
					os.Patch = bytes[4]
					if groupCount >= 5 {
						os.PatchMinor = bytes[5]
					}
				}
			}
		}
	}
}

func (os *Os) ToString() string {
	var str string
	if os.Family != "" {
		str += os.Family
	}
	version := os.ToVersionString()
	if version != "" {
		str += " " + version
	}
	return str
}

func (os *Os) ToVersionString() string {
	var version string
	if os.Major != "" {
		version += os.Major
	}
	if os.Minor != "" {
		version += "." + os.Minor
	}
	if os.Patch != "" {
		version += "." + os.Patch
	}
	if os.PatchMinor != "" {
		version += "." + os.PatchMinor
	}
	return version
}
