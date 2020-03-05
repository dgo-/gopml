package gopml

/*
OPML Spec:
1.0: http://dev.opml.org/spec1.html
2.0: http://dev.opml.org/spec2.html
*/

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// OpmlTime supports RFC822 with four digits for the year
	OpmlTime = "02 Jan 2006 15:04 MST"
)

// OPML data structure
type OPML struct {
	Version string `xml:"version,attr"`
	Head    Head   `xml:"head"`
	Body    Body   `xml:"body"`
}

// Head hold the OPML header
type Head struct {
	Title           string `xml:"title"`
	DateCreated     *oTime `xml:"dateCreated,omitempty"`
	DateModified    *oTime `xml:"dateModified,omitempty"`
	OwnerName       string `xml:"ownerName,omitempty"`
	OwnerEmail      string `xml:"ownerEmail,omitempty"`
	OwnerID         string `xml:"ownerId,omitempty"`
	Docs            string `xml:"docs,omitempty"`
	ExpansionState  string `xml:"expansionState,omitempty"`
	VertScrollState int    `xml:"vertScrollState,omitempty"`
	WindowTop       int    `xml:"windowTop,omitempty"`
	WindowBottom    int    `xml:"windowBottom,omitempty"`
	WindowLeft      int    `xml:"windowLeft,omitempty"`
	WindowRight     int    `xml:"windowRight,omitempty"`
}

// Body hold all outlines objects
type Body struct {
	Outlines []Outline `xml:"outline"`
}

// Outline hold an outline object
type Outline struct {
	Outlines     []Outline `xml:"outline"`
	Text         string    `xml:"text,attr"`
	Type         string    `xml:"type,attr,omitempty"`
	IsComment    *oBool    `xml:"isComment,attr,omitempty"`
	IsBreakpoint *oBool    `xml:"isBreakpoint,attr,omitempty"`
	Created      *oTime    `xml:"created,attr,omitempty"`
	Category     []string  `xml:"category,attr,omitempty"`
	XMLURL       string    `xml:"xmlUrl,attr,omitempty"`
	HTMLURL      string    `xml:"htmlUrl,attr,omitempty"`
	URL          string    `xml:"url,attr,omitempty"`
	Language     string    `xml:"language,attr,omitempty"`
	Title        string    `xml:"title,attr,omitempty"`
	Version      string    `xml:"version,attr,omitempty"`
	Description  string    `xml:"description,attr,omitempty"`
}

/************** custom datatayes ***************/

type oBool bool

func (field *oBool) UnmarshalXMLAttr(attr xml.Attr) error {
	value := attr.Value

	// if no value is given default ist false
	if value == "" {
		*field = oBool(false)
		return nil
	}

	bvalue, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	*field = oBool(bvalue)
	return nil
}

func (field oBool) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var attr xml.Attr
	if field == false {
		return attr, nil
	}

	attr.Name = name
	attr.Value = strconv.FormatBool(bool(field))
	return attr, nil
}

func (value *oBool) Bool() bool {
	if value == nil {
		return false
	}
	return bool(*value)
}

type oTime struct {
	time.Time
}

func (field *oTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var in string
	d.DecodeElement(&in, &start)
	if in == "" {
		return nil
	}

	t, err := ParseTimeOPML(in)
	if err != nil {
		return err
	}
	*field = oTime{t}
	return nil
}

func (field *oTime) UnmarshalXMLAttr(attr xml.Attr) error {
	value := attr.Value
	if value == "" {
		return nil
	}

	t, err := ParseTimeOPML(value)
	if err != nil {
		return err
	}
	*field = oTime{t}
	return nil
}

func (field oTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var t string
	if field.IsZero() {
		return nil
	}
	t = field.Format(OpmlTime)
	return e.EncodeElement(t, start)
}

func (field oTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	var attr xml.Attr
	if field.IsZero() {
		return attr, nil
	}

	attr.Name = name
	attr.Value = field.Format(OpmlTime)
	return attr, nil
}

/************** Functions ************************/

// ParseTimeOPML parse opml supported time strings
func ParseTimeOPML(inputString string) (time.Time, error) {
	// prefer four digiti year
	t, err := time.Parse(OpmlTime, inputString)
	if err != nil {
		t, err = time.Parse(time.RFC822, inputString)
	}
	return t, err
}

// Parse an the given data for opml file
func Parse(data io.Reader) (OPML, error) {
	var opml OPML

	decoder := xml.NewDecoder(data)
	err := decoder.Decode(&opml)
	return opml, err
}

// ParseFile load an OPML from the given file
func ParseFile(filePath string) (OPML, error) {
	var opml OPML
	file, err := os.Open(filePath)
	if err == nil {
		opml, err = Parse(file)
	}
	return opml, err
}

// ParseHTTP GET an OPML file over http and parse it
func ParseHTTP(url string) (OPML, error) {
	var opml OPML
	res, err := http.Get(url)
	if err == nil {
		defer res.Body.Close()
		content, err := ioutil.ReadAll(res.Body)
		if err == nil {
			opml, err = Parse(bytes.NewReader(content))
		}
	}
	return opml, err
}

// ParseString parse str to an OPML
func ParseString(str string) (OPML, error) {
	return Parse(strings.NewReader(str))
}

// Byte return opml in bytes
func (opml OPML) Byte() ([]byte, error) {
	return xml.MarshalIndent(opml, "", "\t")
}

// String return opml as string
func (opml OPML) String() (string, error) {
	b, err := opml.Byte()
	return string(b), err
}

// WriteFile write opml to given file
func (opml OPML) WriteFile(filename string) error {
	data, err := opml.Byte()
	if err != nil {
		return err
	}

	fh, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = fh.Write(data)
	return err
}
