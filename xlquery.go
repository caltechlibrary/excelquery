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
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/xlquery/rss2"

	// 3rd Party packages
	"github.com/tealeg/xlsx"
)

const (
	// Version of this package
	Version = `v0.0.2`
)

var (
	resultLabels = []string{
		"Title",
		"Description",
		"Link",
		"GUID",
	}
	resultMap = map[string]string{
		"Title":       ".item[].title",
		"Description": ".item[].description",
		"Link":        ".item[].link",
		"GUID":        ".item.[].guid",
	}
)

// XLQuery holds the settings to run the XLQuery process over a spreadsheet contacting the
// EPrints repository search CGI script.
type XLQuery struct {
	EPrintsSearchURL string
	ResultDataPaths  []string
	WorkbookName     string
	SheetName        string
	QueryColumn      string
	ResultSheetName  string
	SkipFirstRow     bool
	OverwriteResult  bool
	DataURL          string
	ErrorList        []string
}

// ColumnNameToIndex turns a column reference e.g. 'A', 'BF' into a zero-based array position
func ColumnNameToIndex(colName string) (int, error) {
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
		return -1, errors.New("No column letter provided")
	}
	sum := 0
	letters := strings.Split(strings.ToUpper(colName), "")
	for i, l := range letters {
		if v, ok := m[l]; ok == true {
			sum = sum * 26
			sum += v
		} else {
			return -1, errors.New(`Can't find value for "` + letters[i] + `" in "` + colName + `"`)
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
	cell := sheet.Cell(row, col)
	if overwrite == false && cell.Value != "" {
		return errors.New(`Cell already has a value ` + cell.Value)
	}
	cell.Value = value
	// Update the style to use TextWrap = true
	style := cell.GetStyle()
	style.Alignment.WrapText = true
	cell.SetStyle(style)
	return nil
}

// UpdateParameters adds/overwrites any mapped values to the URL object passed in.
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
func UpdateParameters(api *url.URL, queryTerms map[string]string) *url.URL {
	q := api.Query()
	for key, val := range queryTerms {
		q.Set(key, val)
	}
	api.RawQuery = q.Encode()
	return api
}

// Request executes an HTTP request to the service returning a Query structure
// and error value.
func Request(api *url.URL, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", api.String(), nil)
	if err != nil {
		return nil, err
	}

	for ky, val := range headers {
		req.Header.Add(ky, val)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// given an RSS2 document return all the entries matching so we can apply some sort of data path
// e.g. .version, .channel.title, .channel.link, .item[].link, .item[].guid, .item[].title, .item[].description

// CliRunner is the run method for a command line tool
func CliRunner(xlq *XLQuery, println func(string)) error {
	var (
		resultSheet  *xlsx.Sheet
		saveWorkbook bool
		start        int
		err          error
		ok           bool
	)
	workbook, err := xlsx.OpenFile(xlq.WorkbookName)
	if err != nil {
		return errors.New("Can't open " + xlq.WorkbookName + ", " + err.Error())
	}
	sheet, ok := workbook.Sheet[xlq.SheetName]
	if ok == false {
		return errors.New("Can't read " + xlq.WorkbookName + "." + xlq.SheetName + ", " + err.Error())
	}
	qIndex, err := ColumnNameToIndex(xlq.QueryColumn)
	if err != nil {
		return errors.New("Can't find column " + xlq.QueryColumn + ", in " + xlq.WorkbookName + "." + xlq.SheetName + ", " + err.Error())
	}

	// Use an existing sheet or create a new one to save results in.
	if xlq.OverwriteResult == false {
		resultSheet, err = workbook.AddSheet(xlq.ResultSheetName)
		if err != nil {
			return errors.New("Can't create " + xlq.WorkbookName + "." + xlq.SheetName + ", " + err.Error())
		}
	} else {
		resultSheet, ok = workbook.Sheet[xlq.ResultSheetName]
		if ok == false {
			return errors.New("Can't find " + xlq.WorkbookName + "." + xlq.SheetName)
		}
	}

	// This defaults to CaltechAUTHORs advanced search, can be overwritten in the environment.
	eprintsAPI, err := url.Parse(xlq.EPrintsSearchURL)
	if err != nil {
		return errors.New("Can't parse CaltechAUTHORS URL " + xlq.EPrintsSearchURL + ", " + err.Error())
	}
	if xlq.SkipFirstRow == true {
		start = 1
	}
	for i := range sheet.Rows {
		if i >= start {
			// Update the search paraters
			searchString := GetCell(sheet, i, qIndex)
			eprintsAPI = UpdateParameters(eprintsAPI, map[string]string{
				"title":  searchString,
				"output": "RSS2",
			})
			buf, err := Request(eprintsAPI, map[string]string{})
			if err != nil {
				xlq.Error(eprintsAPI.String() + " request failed, " + err.Error())
			} else {
				feed, err := rss2.Parse(buf)
				if err != nil {
					xlq.Error("Can't parse response " + eprintsAPI.String() + ", " + err.Error())
				} else {
					results, err := feed.Filter([]string{
						resultMap["Title"],
						resultMap["Description"],
						resultMap["Link"],
						resultMap["GUID"],
					})
					if err != nil {
						xlq.Error("filter on link error, " + err.Error())
						saveWorkbook = false
					} else {
						// Write new Row (Iterate through columns using UpdateCell) for results sheet
						for col, key := range resultLabels {
							//FIXME: Is "i" the row we want to write to?
							err = UpdateCell(resultSheet, i, col, results[key].(string), true)
							if err != nil {
								xlq.Error("Can't update " + xlq.WorkbookName + "." + xlq.SheetName + ", " + err.Error())
							} else {
								saveWorkbook = false
							}
						}
					}
					results = nil
				}
				feed = nil
			}
			buf = nil
		}
	}
	if saveWorkbook == true {
		err := workbook.Save(xlq.WorkbookName)
		if err != nil {
			xlq.Error("Can't save " + xlq.WorkbookName + ", " + err.Error())
			return errors.New(xlq.Errors())
		}
		println("Wrote " + xlq.WorkbookName)
	}
	if len(xlq.ErrorList) > 0 {
		return errors.New(xlq.Errors())
	}
	return nil
}

//
// Web code, the following functions are for using with GopherJS
//

// Init initializes a XLQuery object with reasonable values.
func (xlq *XLQuery) Init() {
	xlq.EPrintsSearchURL = `http://authors.library.caltech.edu/cgi/search/advanced/`
	xlq.ResultDataPaths = []string{
		resultMap["Title"],
		resultMap["Description"],
		resultMap["Link"],
		resultMap["GUID"],
	}
	xlq.WorkbookName = `Untitled.xlsx`
	xlq.SheetName = `Sheet1`
	xlq.QueryColumn = ``
	xlq.ResultSheetName = `Result1`
	xlq.SkipFirstRow = true
	xlq.OverwriteResult = false
	xlq.DataURL = ``
	xlq.ErrorList = []string{}
}

func (xlq *XLQuery) Error(e interface{}) {
	switch e.(type) {
	case error:
		xlq.ErrorList = append(xlq.ErrorList, e.(error).Error())
	case string:
		xlq.ErrorList = append(xlq.ErrorList, e.(string))
	}
}

func (xlq *XLQuery) Errors() string {
	return strings.Join(xlq.ErrorList, "\n")
}

// dataURLPrefix gets the data URL's previs up through and including ';base64,'
func dataURLPrefix(src string) string {
	// Notes: "data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,"
	var (
		pre = `data:`
		b64 = `;base64,`
	)
	if strings.HasPrefix(src, pre) && strings.Contains(src, b64) {
		i := strings.Index(src, b64)
		if i > -1 {
			return src[9 : i+8]
		}
	}
	// An emptry string means no prefix found
	return ""
}

// dataURLToByteArray converts a data URL to byte array or returns an error
func dataURLToByteArray(pre, src string) ([]byte, error) {
	if strings.HasPrefix(src, pre) {
		return base64.StdEncoding.DecodeString(strings.TrimPrefix(src, pre))
	}
	return []byte(src), errors.New("Not a data URL " + pre)
}

// byteArrayToDataURL takes a prefix and byte array truning a string formatted as a data URL
func byteArrayToDataURL(pre string, buf []byte) string {
	// Notes: "data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,"
	return pre + base64.StdEncoding.EncodeToString(buf)
}

// WebRunner takes simple JS types as parameters and returns a dataURL string
func (xlq *XLQuery) WebRunner(dataURL string) string {
	//FIXME: need a real implementation
	return `data:text/plain;base64,` + base64.StdEncoding.EncodeToString([]byte("Hello World!!!"))
}
