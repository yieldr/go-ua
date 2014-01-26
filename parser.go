package ua

import (
	"flag"
	"io/ioutil"
	"launchpad.net/goyaml"
)

type Parser struct {
	UserAgentPatterns []UserAgentPattern `yaml:"user_agent_parsers"`
	OsPatterns        []OsPattern        `yaml:"os_parsers"`
	DevicePatterns    []DevicePattern    `yaml:"device_parsers"`
}

type Client struct {
	UserAgent *UserAgent
	Os        *Os
	Device    *Device
}

var args struct {
	file string
}

func init() {
	flag.StringVar(&args.file, "regex", "regexes.yaml", "File containing regular expressions able to match user agents.")
}

func New() *Parser {
	return NewFromFile(args.file)
}

func NewFromFile(file string) *Parser {
	p := new(Parser)

	b, err := ioutil.ReadFile(args.file)
	if nil != err {
		panic(err)
	}

	err = goyaml.Unmarshal(b, &p)
	if err != nil {
		panic(err)
	}

	return p
}

func (parser *Parser) ParseUserAgent(line string) *UserAgent {
	ua := new(UserAgent)
	found := false
	for _, uaPattern := range parser.UserAgentPatterns {
		uaPattern.Match(line, ua)
		if len(ua.Family) > 0 {
			found = true
			break
		}
	}
	if !found {
		ua.Family = "Other"
	}
	return ua
}

func (parser *Parser) ParseOs(line string) *Os {
	os := new(Os)
	found := false
	for _, osPattern := range parser.OsPatterns {
		osPattern.Match(line, os)
		if len(os.Family) > 0 {
			found = true
			break
		}
	}
	if !found {
		os.Family = "Other"
	}
	return os
}

func (parser *Parser) ParseDevice(line string) *Device {
	for _, pattern := range parser.DevicePatterns {
		if device, ok := pattern.Match(line); ok {
			return device
		}
	}
	return unknownDevice()
}

func (parser *Parser) Parse(line string) *Client {
	cli := new(Client)
	cli.UserAgent = parser.ParseUserAgent(line)
	cli.Os = parser.ParseOs(line)
	cli.Device = parser.ParseDevice(line)
	return cli
}
