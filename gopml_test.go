package gopml

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const (
	// Files for testing
	FileDir             = "test_data/"
	FileTiny            = "tiny.opml"
	FileValid           = "valid.opml"
	FileInvalid         = "invalid.opml"
	FileInvalidDate     = "invalid-date.opml"
	FileInvalidDateAttr = "invalid-date-attr.opml"
	FileInvalidBool     = "invalid-bool.opml"

	//
	TestWrite = "write.opml"

	// result strings
	ResultTiny = "<OPML version=\"1.0\">\n\t<head>\n\t\t<title>tiny</title>\n\t</head>\n\t<body>\n\t\t<outline text=\"text\" title=\"tiny\"></outline>\n\t</body>\n</OPML>"
)

func TestParseTimeOPML(t *testing.T) {
	now := time.Now()

	// test rfc822
	rfc822, err := ParseTimeOPML(now.Format(time.RFC822))
	if err != nil {
		t.Error("Unable to parse RFC822: " + err.Error())
	}
	if now.Format(time.RFC822) != rfc822.Format(time.RFC822) {
		t.Error("wrong result at parse RFC822")
	}

	// test opmltime (four digits year)
	opmltime, err := ParseTimeOPML(now.Format(OpmlTime))
	if err != nil {
		t.Error("Unable to parse opmltime: " + err.Error())
	}
	if now.Format(OpmlTime) != opmltime.Format(OpmlTime) {
		t.Error("wrong result at parse opmltime")
	}

	// test fail
	_, err = ParseTimeOPML("no valid time")
	if err == nil {
		t.Error("parse invalid timestamp")
	}
}

func TestParseFile(t *testing.T) {
	filename := FileDir + FileTiny
	result, err := ParseFile(filename)

	if err != nil {
		t.Error(err.Error())
	}

	str, err := result.String()
	if err != nil {
		t.Error(err.Error())
	}

	if str != ResultTiny {
		t.Error("Result match not with wanted result")
	}
}

func TestParseString(t *testing.T) {
	result, err := ParseString(ResultTiny)
	if err != nil {
		t.Error("failed to parse string")
	}

	str, err := result.String()
	if err != nil {
		t.Error(err.Error())
	}

	if str != ResultTiny {
		t.Error("Result match not with wanted result")
	}
}

func TestParseHTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer ts.Close()

	result, err := ParseHTTP(ts.URL)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = result.String()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestWriteFile(t *testing.T) {
	inputfile := FileDir + FileTiny
	opml, err := ParseFile(inputfile)
	if err != nil {
		t.Error(err.Error())
	}

	err = opml.WriteFile("/proc/no_access")
	if err == nil {
		t.Error(err.Error())
	}
	err = opml.WriteFile(TestWrite)
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(TestWrite)

	input, err := ioutil.ReadFile(inputfile)
	if err != nil {
		t.Error(err.Error())
	}

	output, err := ioutil.ReadFile(TestWrite)

	if err != nil {
		t.Error(err.Error())
	}

	if bytes.Equal(input, output) != true {
		t.Error("input and output file not equal")
	}
}

func TestParse_invalid(t *testing.T) {
	// parse invalid file
	invalidFiles := []string{FileInvalid, FileInvalidDate, FileInvalidDateAttr, FileInvalidBool}
	for _, f := range invalidFiles {
		inputfile := FileDir + f
		file, err := os.Open(inputfile)
		if err != nil {
			t.Error(err.Error())
		}

		_, err = Parse(file)
		if err == nil {
			t.Error("parse somehow invalid file")
		}
	}
}

func TestParse_valid(t *testing.T) {
	// parse working file
	inputfile := FileDir + FileValid
	file, err := os.Open(inputfile)
	if err != nil {
		t.Error(err.Error())
	}
	defer file.Close()

	opml, err := Parse(file)
	if err != nil {
		t.Error(err.Error())
	}

	_, err = opml.Byte()
	if err != nil {
		t.Error(err.Error())
	}
}
