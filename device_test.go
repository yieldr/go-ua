package ua

import (
	"io/ioutil"
	"launchpad.net/goyaml"
	"testing"
)

type TestCase struct {
	UA     string `yaml:"user_agent_string"`
	Family string `yaml:"family"`
}

func TestDevice(t *testing.T) {
	b, err := ioutil.ReadFile("./resources/test_device.yaml")
	if err != nil {
		t.Fatal("Unable to locate resource file.")
	}

	var tests map[string][]TestCase

	err = goyaml.Unmarshal(b, &tests)
	if err != nil {
		t.Fatal("Unable to unmarshal yaml.")
	}

	parser, _ := NewParser("regexes.yaml")

	if len(tests) == 0 {
		t.Skip("No test cases found")
	}

	for _, test := range tests["test_cases"] {
		device := parser.ParseDevice(test.UA)
		if device.Family != test.Family {
			t.Error("Expected:", device.Family, "Actual:", test.Family)
		}
	}
}
