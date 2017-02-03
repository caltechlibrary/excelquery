//
// excelquery - a package for quering Caltech library API (and others) and integrating results into an Excel Workbook.
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
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/excelquery"
)

var (
	usage = `USAGE: %s XLSX_FILENAME SHEET_NO ROW_NO COLUMN_NO`

	description = `

%s query our repositories for matching information.
`

	examples = `
EXAMPLE

	%s inventory.xlsx 0 20 20

Query sheet number 0 (the first sheet) based on row 20,
column 20.
`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	eprintsSearchURL = "http://authors.library.caltech.edu/cgi/search/advanced/"
	sheetName        = "Sheet1"
	resultSheetName  = "Result"
	skipFirstRow     = true
)

func init() {
	// General flags
	flag.BoolVar(&showHelp, "h", false, "show help information")
	flag.BoolVar(&showHelp, "help", false, "show help information")
	flag.BoolVar(&showVersion, "v", false, "show version information")
	flag.BoolVar(&showVersion, "version", false, "show version information")
	flag.BoolVar(&showLicense, "l", false, "show license information")
	flag.BoolVar(&showLicense, "license", false, "show license information")

	// App specific flags
	flag.BoolVar(&skipFirstRow, "s", skipFirstRow, "set boolean for skipping first row of sheet (default true)")
	flag.BoolVar(&skipFirstRow, "skip", skipFirstRow, "set boolean for skipping first row of spreadsheet (default true)")

	// Set from environment
	if val := os.Getenv("EPRINTS_SEARCH_URL"); val != "" {
		eprintsSearchURL = val
	}
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()

	// Configuration and command line interation
	cfg := cli.New(appName, appName, fmt.Sprintf(excelquery.LicenseText, appName, excelquery.Version), excelquery.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName)
	cfg.ExampleText = fmt.Sprintf(examples, appName)

	if showHelp == true {
		fmt.Println(cfg.Usage())
		os.Exit(0)
	}

	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}

	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "USAGE: %s XLXS_FILENAME SHEET_NAME QUERY_COLUMN [RESULT_SHEET_NAME]\n", appName)
		os.Exit(1)
	}
	resultSheetName = "Result"
	fname, sheetName, queryColumn := args[0], args[1], args[2]
	if len(args) >= 4 {
		resultSheetName = args[3]
	}

	fmt.Printf("Workbook name: %s, query sheet %s, query column: %s, result sheet: %s\n", fname, sheetName, queryColumn, resultSheetName)
	xlq := new(excelquery.XLQuery)
	xlq.Init()
	xlq.EPrintsSearchURL = eprintsSearchURL
	xlq.WorkbookName = fname
	xlq.SheetName = sheetName
	xlq.QueryColumn = queryColumn
	xlq.ResultSheetName = resultSheetName
	xlq.OverwriteResult = true
	xlq.SkipFirstRow = skipFirstRow

	err := excelquery.CliRunner(xlq, func(msg string) {
		fmt.Fprintf(os.Stdout, "%s\n", msg)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
