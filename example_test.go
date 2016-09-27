//
// exampl_test.go provides example code for demoing xlquery
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
package xlquery_test

import (
	"fmt"
	"path"

	// Caltech Library package
	"github.com/caltechlibrary/xlquery"
)

// ExampleXLQuery create a XLQuery object for running latter
func ExampleXLQuery() {
	// Create a new empty XLQuery
	xlq := new(xlquery.XLQuery)
	// Set some sane defaults
	xlq.Init()
	// Now you can overwrite them as necessary...
	xlq.WorkbookName = path.Join("testdata", "test-1.xlsx")
	xlq.SheetName = "Sheet1"
	xlq.QueryColumn = "A"
	xlq.ResultSheetName = "Result1"
	xlq.OverwriteResult = true
	xlq.SkipFirstRow = true
	// At this point you can run xlquery with the CliRunner() or WebRunner() depending
	// on your environment.
}

// ExampleCliRunner uses an XLQuery structure for sain settings and the CliRunner() function to process
func ExampleCliRunner() {
	// Creates a new
	xlq := new(xlquery.XLQuery)
	// Set some sane defaults
	xlq.Init()
	xlq.WorkbookName = path.Join("testdata", "test-1.xlsx")
	xlq.SheetName = "Sheet1"
	xlq.QueryColumn = "A"
	xlq.ResultSheetName = "Result1"
	xlq.OverwriteResult = true
	xlq.SkipFirstRow = true
	err := xlquery.CliRunner(xlq, func(msg string) {
		fmt.Println(msg)
	})
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
}
