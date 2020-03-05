package main

import (
	"bytes"
	"fmt"

	"github.com/dgo-/gopml"
)

func main() {

	// parse opml from web
	url := "https://raw.githubusercontent.com/dgo-/gopml/master/test_data/valid.opml"
	opml, err := gopml.ParseHTTP(url)
	if err != nil {
		panic(err)
	}

	// write opml data to file
	file := "/tmp/test.opml"
	if opml.WriteFile(file) != nil {
		panic(err)
	}

	// read opml from file
	opml, err = gopml.ParseFile(file)
	if err != nil {
		panic(err)
	}

	// return a opml string
	str, err := opml.String()
	if err != nil {
		panic(err)
	}

	// read opml from string
	opml, err = gopml.ParseString(str)
	if err != nil {
		panic(err)
	}

	// return a opml byte string
	byteData, err := opml.Byte()
	if err != nil {
		panic(err)
	}

	// read opml from io.reader
	reader := bytes.NewReader(byteData)
	opml, err = gopml.Parse(reader)
	if err != nil {
		panic(err)
	}

	// use opml data struct
	fmt.Println(opml.Head.Title)

	// use gopml.oTime struct to handle time objects
	date := opml.Head.DateCreated.Time
	fmt.Printf("Type of Date: %T\n", date)
	fmt.Println(date)

	// Print all outlines with isComment enabled
	for _, out := range opml.Body.Outlines {
		if out.IsComment.Bool() {
			fmt.Println(out.Title)
		}
	}

	// print title of all outines that have 4 sub outlines
	for _, out := range opml.Body.Outlines {
		if len(out.Outlines) == 4 {
			fmt.Println(out.Title)
		}
	}
}
