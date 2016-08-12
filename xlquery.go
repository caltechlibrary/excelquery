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
	"fmt"
	"net/url"
	"strings"

	// 3rd Party packages
	"github.com/tealeg/xlsx"
)

const (
	// Version of this package
	Version = "v0.0.1"
)

// ColumnToInt turns a column reference e.g. 'A', 'BF' into an interger value
func ColumnToInt(colName string) (int, error) {
	m := map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
		"D": 4,
		"E": 5,
		"F": 6,
		"G": 7,
		"H": 8,
		"I": 9,
		"J": 10,
		"K": 11,
		"L": 12,
		"M": 13,
		"N": 14,
		"O": 15,
		"P": 16,
		"Q": 17,
		"R": 18,
		"S": 19,
		"T": 20,
		"U": 21,
		"V": 22,
		"W": 23,
		"X": 24,
		"Y": 25,
		"Z": 26,
	}
	if strings.TrimSpace(colName) == "" {
		return -1, fmt.Errorf("No column letter provided")
	}
	sum := 0
	letters := strings.Split(strings.ToUpper(colName), "")
	for i, l := range letters {
		if v, ok := m[l]; ok == true {
			sum = sum * 26
			sum += v
		} else {
			return -1, fmt.Errorf("Can't find value for %q in %q", letters[i], colName)
		}
	}
	return sum - 1, nil
}

// GetCell given a Spreadsheet, row and col, return the query string or error
func GetCell(sheet *xlsx.Sheet, row int, col int) string {
	cell := sheet.Cell(row, col)
	if cell != nil {
		return cell.Value
	}
	return ""
}

// UpdateCell given a Spreadsheeet, row and col, save the value respecting the overWrite flag or return an error
func UpdateCell(sheet *xlsx.Sheet, row int, col int, value string, overwrite bool) error {
	//FIXME: add support for case of missing column or row.
	cell := sheet.Cell(row, col)
	if overwrite == false && cell.Value != "" {
		return fmt.Errorf("Cell(%d, %d) already has a value %s", row, col, cell.Value)
	}
	cell.Value = value
	return nil
}

// UpdateQuery adds any mapped values to the URL object passed in.
//
// URL attribute for EPrints advanced search (output is Atom):
//  Scheme: http
//  Host: eprint-repository.example.org
//  Path: /cgi/search/advanced
//  Query parameters:
// 		title: Molecules in solutoin
// 		output: Atom
//
// Example usage:
// api, _ := url.Parse("http://eprint-repository.example.org/cgi/search/advanced")
// xlquery.UpdateQuery(api, map[string]string{"title": title, "output":"Atom"})
// data, err := http.Get(api.String())
// ...
func UpdateQuery(api *url.URL, queryTerms map[string]string) *url.URL {
	q := api.Query()
	for key, val := range queryTerms {
		q.Set(key, val)
	}
	api.RawQuery = q.Encode()
	return api
}
