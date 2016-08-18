//
// xlquery/webapp/webapp.go is a wrapper for xlquery.go targetting GopherJS and embedding xlquery functionality as a webapp in a web browser.
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
	"bytes"
	"io"
	"net/url"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/xlquery"
	"github.com/caltechlibrary/xlquery/rss2"

	// 3rd Party Library
	"github.com/gopherjs/gopherjs/js"
	"github.com/tealeg/xlsx"
)

// A structure to hold the Object in JS for accessing Go code
type XLQuery struct {
	Version          string
	EPrintsSearchURL string
	SheetName        string
	DataPath         string
	Overwrite        bool
}

type XLQResponse struct {
	Data   string
	Errors []string
}

// Error appends an error string to the Errors property of XLQResponse
func (xlqr *XLQResponse) Error(s string) {
	xlqr.Errors = append(xlqr.Errors, s)
}

// Init sets the default values of the XLQuery object.
func (xlq *XLQuery) Init() {
	xlq.Version = xlquery.Version
	xlq.EPrintsSearchURL = "http://authors.library.caltech.edu/cgi/search/advanced/"
	xlq.SheetName = "Sheet1"
	xlq.DataPath = ".item[].link"
	xlq.Overwrite = false
}

// Run take the byte array of the raw XLSX source code and
// processes it much like *main* function in the cli.
func (xlq *XLQuery) Run(data, queryColumn, resultColumn string) *XLQResponse {
	xlqResponse := new(XLQResponse)
	workbook, err := xlsx.OpenBinary([]byte(data))
	if err != nil {
		xlqResponse.Error("Can't read the xlsx content " + err.Error())
		return xlqResponse
	}
	if sheet, ok := workbook.Sheet[xlq.SheetName]; ok == true {
		qIndex, err := xlquery.ColumnNameToIndex(queryColumn)
		if err != nil {
			xlqResponse.Error("Can't find column " + queryColumn + ", in " + xlq.SheetName + ", " + err.Error())
			return xlqResponse
		}
		rIndex, err := xlquery.ColumnNameToIndex(resultColumn)
		if err != nil {
			xlqResponse.Error("Can't find column " + resultColumn + ", in " + xlq.SheetName + ", " + err.Error())
			return xlqResponse
		}

		// This defaults to CaltechAUTHORs advanced search, can be overwritten in the environment.
		eprintsAPI, err := url.Parse(xlq.EPrintsSearchURL)
		if err != nil {
			xlqResponse.Error("Can't parse CaltechAUTHORS URL " + xlq.EPrintsSearchURL + ", " + err.Error())
			return xlqResponse
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
				buf, err := xlquery.Request(eprintsAPI, map[string]string{})
				if err != nil {
					xlqResponse.Error(eprintsAPI.String() + " request failed, " + err.Error())
				} else {
					feed, err := rss2.Parse(buf)
					if err != nil {
						xlqResponse.Error("Can't parse response " + eprintsAPI.String() + ", " + err.Error())
					} else {
						links, err := feed.Filter(xlq.DataPath)
						if err != nil {
							xlqResponse.Error("filter on link error, " + err.Error())
						} else if links != nil {
							s := strings.Join(links[xlq.DataPath].([]string), "\r")
							err = xlquery.UpdateCell(sheet, i, rIndex, s, xlq.Overwrite)
							if err != nil {
								xlqResponse.Error("Failed to update cell results for " + searchString + ", " + err.Error())
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
			var buf bytes.Buffer
			wr := io.Writer(&buf)
			err := workbook.Write(wr)
			if err != nil {
				xlqResponse.Error(err.Error())
			}
			xlqResponse.Data = buf.String()
		}
	}
	return xlqResponse
}

func NewXLQuery() *js.Object {
	return js.MakeWrapper(&XLQuery{})
}

func NewXLQResponse() *js.Object {
	return js.MakeWrapper(&XLQResponse{})
}

func main() {
	js.Global.Set("xlquery", map[string]interface{}{
		"New": NewXLQuery,
	})
	js.Global.Set("xlresponse", map[string]interface{}{
		"New": NewXLQResponse,
	})
}
