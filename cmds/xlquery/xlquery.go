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
	"os"
	"path"

	// Caltech Library packages
	"github.com/caltechlibrary/xlquery"
)

var (
	showHelp    bool
	showVersion bool
	showLicense bool

	eprintsSearchURL = "http://authors.library.caltech.edu/cgi/search/advanced/"
	sheetName        = "Sheet1"
	resultSheetName  = "Result1"
	overwriteResult  = false
	skipFirstRow     = true
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
 + The environment varaible EPRINTS_SEARCH_URL can overwrite the default
   CaltechAUTHORS search URL.

 OPTIONS

`, appName)

	flag.VisitAll(func(f *flag.Flag) {
		if len(f.Name) > 1 {
			fmt.Printf("    -%s, -%s\t%s\n", f.Name[0:1], f.Name, f.Usage)
		}
	})

	fmt.Fprintf(fp, `
	
 Version: %s

`, xlquery.Version)
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
	flag.BoolVar(&skipFirstRow, "S", skipFirstRow, "set boolean for skipping first row of sheet (default true)")
	flag.BoolVar(&skipFirstRow, "Skip", skipFirstRow, "set boolean for skipping first row of spreadsheet (default true)")
	flag.StringVar(&sheetName, "s", sheetName, "set the sheet name, e.g. \"Sheet1\"")
	flag.StringVar(&sheetName, "sheet", sheetName, "set the sheet name, e.g. \"Sheet1\"")
	flag.StringVar(&sheetName, "r", resultSheetName, "set the result sheet name, e.g. \"Result1\"")
	flag.StringVar(&sheetName, "result-sheet", resultSheetName, "set the result sheet name, e.g. \"Result1\"")

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
		fmt.Fprintf(os.Stderr, "USAGE: %s XLXS_FILENAME SHEET_NAME QUERY_COLUMN [RESULT_SHEET_NAME]\n", appname)
		os.Exit(1)
	}
	fname, sheetName, queryColumn := args[0], args[1], args[2]
	if len(args) >= 4 {
		resultSheetName = args[3]
	}

	fmt.Printf("Workbook name: %s, queryColumn: %s, resultSheet: %s\n", fname, queryColumn, resultSheetName)
	xlq := new(xlquery.XLQuery)
	xlq.Init()
	xlq.EPrintsSearchURL = eprintsSearchURL
	xlq.WorkbookName = fname
	xlq.SheetName = sheetName
	xlq.QueryColumn = queryColumn
	xlq.ResultSheetName = resultSheetName
	xlq.OverwriteResult = overwriteResult
	xlq.SkipFirstRow = skipFirstRow

	err := xlquery.CliRunner(xlq, func(msg string) {
		fmt.Fprintf(os.Stdout, "%s\n", msg)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
