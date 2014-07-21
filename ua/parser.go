package ua

import (
	"flag"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v1"
)

type Parser struct {
	file     string
	patterns Patterns
	mux      sync.Mutex
}

type Patterns struct {
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

func New() (*Parser, error) {
	return NewParser(args.file)
}

func NewParser(file string) (*Parser, error) {
	patterns, err := LoadPatterns(file)
	if err != nil {
		return nil, err
	}
	parser := &Parser{}
	parser.file = file
	parser.patterns = patterns
	return parser, nil
}

func (p *Parser) Patterns() Patterns {
	p.mux.Lock()
	defer p.mux.Unlock()
	return p.patterns
}

func (p *Parser) SetPatterns(patterns Patterns) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.patterns = patterns
}

func (p *Parser) ParseUserAgent(line string) *UserAgent {
	ua := new(UserAgent)
	found := false
	for _, pattern := range p.Patterns().UserAgentPatterns {
		pattern.Match(line, ua)
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

func (p *Parser) ParseOs(line string) *Os {
	os := new(Os)
	found := false
	for _, pattern := range p.Patterns().OsPatterns {
		pattern.Match(line, os)
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

func (p *Parser) ParseDevice(line string) *Device {
	for _, pattern := range p.Patterns().DevicePatterns {
		if device, ok := pattern.Match(line); ok {
			return device
		}
	}
	return unknownDevice()
}

func (p *Parser) Parse(line string) *Client {
	cli := new(Client)
	cli.UserAgent = p.ParseUserAgent(line)
	cli.Os = p.ParseOs(line)
	cli.Device = p.ParseDevice(line)
	return cli
}

func (p *Parser) Reload() error {
	patterns, err := LoadPatterns(p.file)
	if err != nil {
		return err
	}
	p.SetPatterns(patterns)
	return nil
}

func LoadPatterns(file string) (p Patterns, err error) {
	b, err := ioutil.ReadFile(file)
	if nil != err {
		return
	}
	err = yaml.Unmarshal(b, &p)
	if err != nil {
		return
	}
	return
}
