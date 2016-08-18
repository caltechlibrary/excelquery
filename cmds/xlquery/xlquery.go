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

	overwriteResult bool
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
 USAGE: %s [OPTION] WORKBOOK_NAME SHEET_NAME QUERY_COLUMN RESULT_COLUMN DATA_PATH

 Populate an spreadsheet in ".xlsx" format by using a query column's value as a query string
 updating the result column's value (without overwriting existing data).

 + Sheet name should correspond to the sheet you want to run through (e.g. "Sheet 1")
 + Column names are in Excel's letter format (e.g. "A", "FX", "BBC").
 + data path is the part of the result you want to use (e.g. url matching the title queried)

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
	flag.BoolVar(&showHelp, "h", false, "show help information")
	flag.BoolVar(&showHelp, "help", false, "show help information")
	flag.BoolVar(&showVersion, "v", false, "show version information")
	flag.BoolVar(&showVersion, "version", false, "show version information")
	flag.BoolVar(&showLicense, "l", false, "show license information")
	flag.BoolVar(&showLicense, "license", false, "show license information")

	flag.BoolVar(&overwriteResult, "o", false, "overwrite the results column")
	flag.BoolVar(&overwriteResult, "overwrite", false, "overwrite the results column")
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
	if len(args) < 5 {
		fmt.Fprintf(os.Stderr, "USAGE: %s XLXS_FILENAME SHEET_NAME QUERY_COLUMN RESULT_COLUMN DATA_PATH\n", appname)
		os.Exit(1)
	}
	fname, sname, queryColumn, resultsColumn, dataPath := args[0], args[1], args[2], args[3], args[4]

	fmt.Printf("Test fname: %s, sheet: %s, queryColumn: %s, resultColumn: %s, dataPath %s\n", fname, sname, queryColumn, resultsColumn, dataPath)
	workbook, err := xlsx.OpenFile(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open %s, %s", fname, err)
		os.Exit(1)
	}
	if sheet, ok := workbook.Sheet[sname]; ok == true {
		qIndex, err := xlquery.ColumnNameToIndex(queryColumn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't find column %s, in %s.%s, %s", queryColumn, fname, sname, err)
			os.Exit(1)
		}
		rIndex, err := xlquery.ColumnNameToIndex(resultsColumn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't find column %s, in %s.%s, %s", queryColumn, fname, sname, err)
			os.Exit(1)
		}
		// FIXME: This should not be hardcoded, setup as a environment var? a commmand line option?
		eprintsAPI, err := url.Parse("http://authors.library.caltech.edu/cgi/search/advanced/")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't parse CaltechAUTHORS URL %s", err)
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
			workbook.Save(fname)
		}
	}
}
