package ua

import (
	"fmt"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v1"
)

func uaInitTesting(file string) []map[string]string {
	fmt.Print(file + ": ")
	testFile, _ := ioutil.ReadFile(file)
	testMap := make(map[string][]map[string]string)
	_ = yaml.Unmarshal(testFile, &testMap)
	return testMap["test_cases"]
}

var uaDefaultRegexFile string = "../../regexes.yaml"
var uaParser *Parser = nil

func uaInitParser(regexFile string) {
	if uaParser == nil {
		uaParser, _ = NewParser(regexFile)
	}
}

func uaInitParserNew(regexFile string) {

	if uaParser == nil {
		uaParser, err := NewParser(regexFile)
		if uaParser != nil {
			fmt.Println(err.Error())
		}

	}
}

func uaHelperTest(file string) bool {
	uaInitParserNew(uaDefaultRegexFile)
	uaInitParser(uaDefaultRegexFile)
	tests := uaInitTesting(file)
	for _, test := range tests {

		// Other language ports of ua_parser skips js_ua in testing
		if test["js_ua"] != "" {
			continue
		}

		testingString := test["user_agent_string"]
		ua := uaParser.ParseUserAgent(testingString)

		if ua.Family != test["family"] || ua.Major != test["major"] ||
			ua.Minor != test["minor"] || ua.Patch != test["patch"] {
			fmt.Println("FAIL")
			fmt.Printf("Expected: %v\nActual: %v\n", test, ua)
			return false
		}
	}
	return true
}

func TestUserAgent(t *testing.T) {
	if !uaHelperTest("../../test_resources/test_user_agent_parser.yaml") {
		t.Fail()
	} else {
		fmt.Println("PASS")
	}
}

func TestFirefoxUserAgents(t *testing.T) {
	if !uaHelperTest("../../test_resources/firefox_user_agent_strings.yaml") {
		t.Fail()
	} else {
		fmt.Println("PASS")
	}
}

func TestPgtsBrowsersList(t *testing.T) {
	if !uaHelperTest("../../test_resources/pgts_browser_list.yaml") {
		t.Fail()
	} else {
		fmt.Println("PASS")
	}
}
