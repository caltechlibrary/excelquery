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
package xlquery

import (
	"net/url"
	"path"
	"testing"

	// 3rd Party packages
	"github.com/caltechlibrary/xlquery/rss2"
	"github.com/tealeg/xlsx"
)

func TestColumnNameToIndex(t *testing.T) {
	testVals := map[string]int{
		"A":   0,
		"R":   17,
		"Z":   25,
		"AA":  26,
		"AB":  27,
		"AM":  38,
		"AS":  44,
		"AX":  49,
		"AZ":  51,
		"BA":  52,
		"BE":  56,
		"BY":  76,
		"CA":  78,
		"CL":  89,
		"DO":  118,
		"EG":  136,
		"EZ":  155,
		"FX":  179,
		"GZ":  207,
		"IE":  238,
		"IT":  253,
		"LS":  330,
		"MA":  338,
		"MT":  357,
		"OK":  400,
		"PQ":  432,
		"RD":  471,
		"TJ":  529,
		"ZZ":  701,
		"AAA": 702,
		"AAB": 703,
		"AAZ": 727,
		"ABA": 728,
		"ADQ": 796,
		"ARC": 1146,
		"ARG": 1150,
		"ARM": 1156,
		"ASK": 1180,
		"ASM": 1182,
		"ATM": 1208,
		"AUX": 1245,
		"AVE": 1252,
		"AVI": 1256,
		"AWE": 1278,
		"AWK": 1284,
		"AZZ": 1377,
		"BAA": 1378,
		"BAD": 1381,
		"BAM": 1390,
		"BAT": 1397,
		"BBC": 1406,
		"BED": 1485,
	}

	for s, i := range testVals {
		r, err := columnNameToIndex(s)
		if err != nil {
			t.Errorf("Couldn't convert %s to int, %s", s, err)
		}
		if r != i {
			t.Errorf("ColumnNameToIndex(%q) != %d, returned %d", s, i, r)
		}
	}
}

func TestSheetHandling(t *testing.T) {
	xldoc := xlsx.NewFile()
	sheet, err := xldoc.AddSheet("Sheet1")
	if err != nil {
		t.Errorf("Can't add sheet: %s", err)
		t.FailNow()
	}
	row := sheet.AddRow()
	A := row.AddCell()
	B := row.AddCell()
	A.Value = "Query"
	B.Value = "Results"

	queryTerms := map[string]string{
		"flood characteristics of alluvial": "Flood Characteristics of Alluvial Streams Important to Pipeline Crossings.",
		"gravitational waves in a":          "Gravitational Waves in a Shallow Compressible Liquid",
		"experimental design of low":        "",
	}
	for query, val := range queryTerms {
		row = sheet.AddRow()
		A = row.AddCell()
		A.Value = query
		B = row.AddCell()
		B.Value = val
	}
	// Make our test directory if needed
	fname := path.Join("testdata", "test-0.xlsx")
	err = xldoc.Save(fname)
	if err != nil {
		t.Errorf("Can't save %s, %s", fname, err)
	}

	xldocTest, err := xlsx.OpenFile(fname)
	if err != nil {
		t.Errorf("Can't open %s, %s", fname, err)
		t.FailNow()
	}

	for _, sheet := range xldocTest.Sheets {
		for j, _ := range sheet.Rows {
			q := sheet.Cell(j, 0)
			r := sheet.Cell(j, 2)
			qTest := getCell(sheet, j, 0)
			if q.Value != qTest {
				t.Errorf("GetCell(sheet, %d, 0) expected %s, got %s", j, q.Value, qTest)
			}
			if r.Value != "" {
				err := updateCell(sheet, j, 2, "This is a test", false)
				if err != nil {
					t.Errorf("Expected an err on update to cell %d,2", j)
				}
			}
			err := updateCell(sheet, j, 2, "This is a test 2", true)
			if err != nil {
				t.Errorf("Expected err to be nil on update to cell %d,2, %s", j, err)
			}
			r2 := sheet.Cell(j, 2)
			if r2.Value != "This is a test 2" {
				t.Errorf("Expected %q, got %s", "This is a test 2", r2.Value)
			}
		}
	}
}

func TestQuerySupport(t *testing.T) {
	eprintsAPI, err := url.Parse("http://authors.library.caltech.edu/cgi/search/advanced/")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	eprintsAPI = updateParameters(eprintsAPI, map[string]string{
		"title":  "Molecules in solution",
		"output": "RSS2",
	})
	if eprintsAPI == nil {
		t.Errorf("Something went wrong updating eprintsAPI query")
		t.FailNow()
	}
	buf, err := request(eprintsAPI, map[string]string{})
	if err != nil {
		t.Errorf("Failed to run %s, %s", eprintsAPI.String(), err)
		t.FailNow()
	}
	r, err := rss2.Parse(buf)
	if err != nil {
		t.Errorf("Failed to parse response, buf[0:24] %q, err %q", buf[0:24], err)
		t.FailNow()
	}
	_, err = r.Filter(".item[].title")
	if err != nil {
		t.Errorf("Failed to filter for titles, %s", err)
		t.FailNow()
	}
}
