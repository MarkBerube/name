package name

import (
	"regexp"
	"strconv"
	"testing"
)

func TestRandomName(t *testing.T) {
	var conf *SheetConfig
	conf = NewConfig("https://test.url",
		"123foo",
		"456bar",
		"NameList",
		1,
		"testing:\n")

	l := GetRandomNameList(1, conf, TestingAPIClient{})

	if len(l) != 1 {
		t.Error("Expected 1, got: " + strconv.Itoa(len(l)))
	}

	m, _ := regexp.MatchString("testing:\\n[A-Za-z ]+", l[0])

	if !m {
		t.Error("Name not found, got: " + l[0])
	}
}

func TestMultipleRandomName(t *testing.T) {

	var conf *SheetConfig
	conf = NewConfig("https://test.url",
		"123foo",
		"456bar",
		"NameList",
		1,
		"testing:\n")

	l := GetRandomNameList(3, conf, TestingAPIClient{})

	if len(l) != 1 {
		t.Error("Expected 1, got: " + strconv.Itoa(len(l)))
	}

	m, _ := regexp.MatchString("testing:\\n[A-Za-z ]+\\n[A-Za-z ]+\\n[A-Za-z ]+\\n", l[0])

	if !m {
		t.Error("Name not found, got: " + l[0])
	}
}

func TestTitlesOff(t *testing.T) {
	var conf *SheetConfig
	conf = NewConfig("https://test.url",
		"123foo",
		"456bar",
		"NameList",
		1,
		"testing:\n")

	conf.IgnoreTitles()

	l := GetRandomNameList(1, conf, TestingAPIClient{})

	if len(l) != 1 {
		t.Error("Expected 1, got: " + strconv.Itoa(len(l)))
	}

	m, _ := regexp.MatchString("testing:\\n[A-Za-z ]+(the Bard|the Knight|the Wizard)\\n", l[0])

	if m {
		t.Error("Found a title on name, got: " + l[0])
	}
}

func TestSecondaryAdditionOff(t *testing.T) {
	var conf *SheetConfig
	conf = NewConfig("https://test.url",
		"123foo",
		"456bar",
		"NameList",
		1,
		"testing:\n")

	conf.IgnoreSecondNameAppends()

	l := GetRandomNameList(1, conf, TestingAPIClient{})

	if len(l) != 1 {
		t.Error("Expected 1, got: " + strconv.Itoa(len(l)))
	}

	m, _ := regexp.MatchString("testing:\\n[A-Za-z ]+(foo|bar)[A-Za-z ]+\\n", l[0])

	if m {
		t.Error("Found a secondary addition on name, got: " + l[0])
	}
}

type TestingAPIClient struct{}

func (t TestingAPIClient) getNameList(listRange string, config *SheetConfig) []string {

	switch listRange {
	case "A2:A":
		// first name
		return []string{"Billy", "Joe", "Lenny"}
	case "B2:B":
		// last name
		return []string{"Berube", "Smith", "Anderson"}
	case "C2:C":
		// title
		return []string{"Bard", "Knight", "Wizard"}
	case "D2:D":
		// secondary addition to last name
		return []string{"foo", "bar"}
	}

	return []string{"can't get here"}

}
