//
// xlquery - a package for quering Caltech library API (and others) and integrating results into an Excel Workbook.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2016, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/xlquery"
	"github.com/caltechlibrary/xlquery/rss2"
	"github.com/tealeg/xlsx"
)

var (
	showHelp    bool
	showVersion bool
	showLicense bool

	eprintsSearchURL = "http://authors.library.caltech.edu/cgi/search/advanced/"
	sheetName        = "Sheet1"
	dataPath         = ".item[].link"
	overwriteResult  = false
)

const (
	license = `

Copyright (c) 2016, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

`
)

func usage(fp *os.File, appName string) {
	fmt.Fprintf(fp, `
 USAGE: %s [OPTION] WORKBOOK_NAME QUERY_COLUMN RESULT_COLUMN

 Populate a workbook (e.g. ".xlsx" file) by using a query column's value as a query string
 updating the result column's value (by default it does not overwrite existing data).

 + The default sheet name to use in a workbook is "Sheet1"
 + The default datapath ".item[].link" which represents an RSS item's link field
 + Column names are in letter format (e.g. "A" is column 1, "B" column 2, etc.)

 OPTIONS

`, appName)

	flag.VisitAll(func(f *flag.Flag) {
		if len(f.Name) > 1 {
			fmt.Printf("    -%s, -%s\t%s\n", f.Name[0:1], f.Name, f.Usage)
		}
	})

	fmt.Fprintf(fp, `
	
 Version: %s `, xlquery.Version)
}

func init() {
	// General flags
	flag.BoolVar(&showHelp, "h", false, "show help information")
	flag.BoolVar(&showHelp, "help", false, "show help information")
	flag.BoolVar(&showVersion, "v", false, "show version information")
	flag.BoolVar(&showVersion, "version", false, "show version information")
	flag.BoolVar(&showLicense, "l", false, "show license information")
	flag.BoolVar(&showLicense, "license", false, "show license information")

	// App specific flags
	flag.BoolVar(&overwriteResult, "o", overwriteResult, "overwrite the results column")
	flag.BoolVar(&overwriteResult, "overwrite", overwriteResult, "overwrite the results column")
	flag.StringVar(&sheetName, "s", sheetName, "set the sheet name, e.g. \"Sheet1\"")
	flag.StringVar(&sheetName, "sheet", sheetName, "set the sheet name, e.g. \"Sheet1\"")
	flag.StringVar(&dataPath, "d", dataPath, "set the datapath for results, e.g. \".item[].link\"")
	flag.StringVar(&dataPath, "datapath", dataPath, "set the datapath for results, e.g. \".item[].link\"")

	// Set from environment
	if val := os.Getenv("EPRINTS_SEARCH_URL"); val != "" {
		eprintsSearchURL = val
	}
}

func main() {
	appname := path.Base(os.Args[0])
	flag.Parse()
	if showHelp == true {
		usage(os.Stdout, appname)
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Printf(" Version %s\n", xlquery.Version)
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Printf(" %s\n%s\n", appname, license)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "USAGE: %s XLXS_FILENAME QUERY_COLUMN RESULT_COLUMN\n", appname)
		os.Exit(1)
	}
	fname, queryColumn, resultsColumn := args[0], args[1], args[2]

	fmt.Printf("Workbook name: %s, queryColumn: %s, resultColumn: %s\n", fname, queryColumn, resultsColumn)
	workbook, err := xlsx.OpenFile(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open %s, %s", fname, err)
		os.Exit(1)
	}
	if sheet, ok := workbook.Sheet[sheetName]; ok == true {
		qIndex, err := xlquery.ColumnNameToIndex(queryColumn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't find column %s, in %s.%s, %s", queryColumn, fname, sheetName, err)
			os.Exit(1)
		}
		rIndex, err := xlquery.ColumnNameToIndex(resultsColumn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't find column %s, in %s.%s, %s", queryColumn, fname, sheetName, err)
			os.Exit(1)
		}

		// This defaults to CaltechAUTHORs advanced search, can be overwritten in the environment.
		eprintsAPI, err := url.Parse(eprintsSearchURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't parse CaltechAUTHORS URL %s, %s", eprintsSearchURL, err)
			os.Exit(1)
		}
		saveWorkbook := false
		for i, row := range sheet.Rows {
			// Assume first row of the spreadsheet is headings
			if i > 0 {
				// If row is too short for append necessary cells
				if len(row.Cells) <= rIndex {
					for len(row.Cells) <= rIndex {
						row.AddCell()
					}
				}
				// Update the search paraters
				searchString := xlquery.GetCell(sheet, i, qIndex)
				eprintsAPI = xlquery.UpdateParameters(eprintsAPI, map[string]string{
					"title":  searchString,
					"output": "RSS2",
				})
				fmt.Fprintf(os.Stdout, "%s\n", eprintsAPI.String())
				buf, err := xlquery.Request(eprintsAPI, map[string]string{})
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s request failed, %s", eprintsAPI.String(), err)
				} else {
					feed, err := rss2.Parse(buf)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Can't parse response %s, %s", eprintsAPI.String(), err)
					} else {
						links, err := feed.Filter(dataPath)
						if err != nil {
							fmt.Fprintf(os.Stderr, "filter on link error, %s", err)
						} else if links != nil {
							fmt.Printf("\tfound %d links\n", len(links[dataPath].([]string)))
							s := strings.Join(links[dataPath].([]string), "\r")
							err = xlquery.UpdateCell(sheet, i, rIndex, s, overwriteResult)
							if err != nil {
								fmt.Fprintf(os.Stderr, "Failed to update cell results for %s, %s", searchString, err)
							} else {
								saveWorkbook = true
							}
						}
						links = nil
					}
					feed = nil
				}
				buf = nil
			}
		}
		if saveWorkbook == true {
			fmt.Printf("Writing results to %s\n", fname)
			err := workbook.Save(fname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Can't save %s, %s\n", fname, err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "%s saved.", fname)
		}
	}
}
