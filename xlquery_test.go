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
	"log"
	"path"
	"testing"

	// 3rd Party packages
	"github.com/tealeg/xlsx"
)

func TestQueryEPrints(t *testing.T) {
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
	log.Println("This is a test")

	for i, sheet := range xldocTest.Sheets {
		for j, _ := range sheet.Rows {
			q := sheet.Cell(j, 0)
			r := sheet.Cell(j, 2)
			log.Printf("sheet: %d, row: %d, q: %s, r: %s\n", i, j, q.Value, r.Value)
			err = QueryEPrints(sheet, j, 0, j, 2)
			if err != nil {
				t.Errorf("QueryEPrints(sheet, %d, 0, %d, 2) failed, %s", j, j, err)
				t.FailNow()
			}
		}
	}
}
