# gopml
[![Build Status](https://travis-ci.org/dgo-/gopml.png?branch=master)](https://travis-ci.org/dgo-/gopml)
[![codecov](https://codecov.io/gh/dgo-/gopml/branch/master/graph/badge.svg)](https://codecov.io/gh/dgo-/gopml)
[![goreport](https://goreportcard.com/badge/github.com/dgo-/gopml)](https://goreportcard.com/report/github.com/dgo-/gopml)
[![Go Doc](https://godoc.org/github.com/dgo-/gopml?status.svg)](https://godoc.org/github.com/dgo-/gopml)


gopml is a Go package that parses .opml files.

## Installation
You can install it simple with go get:
```
go get github.com/dgo-/gopml
```

## Usage
There a short example to that use most common functions:
```go
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
```

The plain example source code can be found in the example directory. 


## License
MIT
