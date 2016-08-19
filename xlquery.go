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
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	// 3rd Party packages
	"github.com/caltechlibrary/xlquery/rss2"
	"github.com/tealeg/xlsx"
)

const (
	// Version of this package
	Version = "v0.0.1"
)

// XLQuery holds the information neccessary to process an Excel Workbook for both
// the command line and in a web browser via GopherJS
type XLQuery struct {
	XMLName          xml.Name `json:"-"`
	Version          string   `xml:"version" json:"version"`
	EPrintsSearchURL string   `xml:"eprintsSearchURL" json:"eprintsSearchURL"`
	ResultDataPath   string   `xml:"resultsDataPath" json:"resultsDataPath"`
	WorkbookName     string   `xml:"workbookName" json:"workbookName"`
	SheetName        string   `xml:"sheetName" json:"sheetName"`
	QueryColumn      string   `xml:"queryColumn" json:"queryColumn"`
	ResultColumn     string   `xml:"resultColumn" json:"resultColumn"`
	SkipFirstRow     bool     `xml:"skipFirstRow" json:"skipFirstRow"`
	OverwriteResult  bool     `xml:"overwriteResult" json:"overwriteResult"`
	DataURL          string   `xml:"dataURL" json:"dataURL"`
	Errors           []string `xml:"errors" json:"errors"`
}

func (xlq *XLQuery) Error(e interface{}) {
	switch e.(type) {
	case error:
		xlq.Errors = append(xlq.Errors, e.(error).Error())
	case string:
		xlq.Errors = append(xlq.Errors, e.(string))
	}
}

// dataURLToByteArray converts string dataURL to byte array or returns an error
func dataURLToByteArray(src string) ([]byte, error) {
	var pre = "data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,"
	if strings.HasPrefix(src, pre) {
		return base64.StdEncoding.DecodeString(strings.TrimPrefix(src, pre))

	}
	return []byte(src), errors.New("Not a data URL for type " + pre)
}

// byteArrayToDataURL converts a byte array back into a DataURL
func byteArrayToDataURL(data []byte) string {
	var pre = "data:application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;base64,"
	src := base64.StdEncoding.EncodeToString(data)
	return pre + src
}

// columnNameToIndex turns a column reference e.g. 'A', 'BF' into a zero-based array position
func columnNameToIndex(colName string) (int, error) {
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

// getCell given a Spreadsheet, row and col, return the query string or error
func getCell(sheet *xlsx.Sheet, row int, col int) string {
	cell := sheet.Cell(row, col)
	if cell != nil {
		return cell.Value
	}
	return ""
}

// updateCell given a Spreadsheeet, row and col, save the value respecting the overWrite flag or return an error
func updateCell(sheet *xlsx.Sheet, row int, col int, value string, overwrite bool) error {
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

// updateParameters adds/overwrites any mapped values to the URL object passed in.
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
func updateParameters(api *url.URL, queryTerms map[string]string) *url.URL {
	q := api.Query()
	for key, val := range queryTerms {
		q.Set(key, val)
	}
	api.RawQuery = q.Encode()
	return api
}

// request executes an HTTP request to the service returning a Query structure
// and error value.
func request(api *url.URL, headers map[string]string) ([]byte, error) {
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

// CliRun is the run method for a command line tool
func CliRun(xlq *XLQuery, println func(string)) error {
	workbook, err := xlsx.OpenFile(xlq.WorkbookName)
	if err != nil {
		return errors.New("Can't open " + xlq.WorkbookName + ", " + err.Error())
	}
	if sheet, ok := workbook.Sheet[xlq.SheetName]; ok == true {
		qIndex, err := columnNameToIndex(xlq.QueryColumn)
		if err != nil {
			return errors.New("Can't find column " + xlq.QueryColumn + ", in " + xlq.WorkbookName + "." + xlq.SheetName + ", " + err.Error())
		}
		rIndex, err := columnNameToIndex(xlq.ResultColumn)
		if err != nil {
			return errors.New("Can't find column " + xlq.ResultColumn + ", in " + xlq.WorkbookName + "." + xlq.SheetName + ", " + err.Error())
		}

		// This defaults to CaltechAUTHORs advanced search, can be overwritten in the environment.
		eprintsAPI, err := url.Parse(xlq.EPrintsSearchURL)
		if err != nil {
			return errors.New("Can't parse CaltechAUTHORS URL " + xlq.EPrintsSearchURL + ", " + err.Error())
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
				searchString := getCell(sheet, i, qIndex)
				eprintsAPI = updateParameters(eprintsAPI, map[string]string{
					"title":  searchString,
					"output": "RSS2",
				})
				buf, err := request(eprintsAPI, map[string]string{})
				if err != nil {
					xlq.Error(eprintsAPI.String() + " request failed, " + err.Error())
				} else {
					feed, err := rss2.Parse(buf)
					if err != nil {
						xlq.Error("Can't parse response " + eprintsAPI.String() + ", " + err.Error())
					} else {
						links, err := feed.Filter(xlq.ResultDataPath)
						if err != nil {
							xlq.Error("filter on link error, " + err.Error())
						} else if links != nil {
							s := strings.Join(links[xlq.ResultDataPath].([]string), "\n")
							if s != "" {
								println(`Searching for "` + searchString + `", found: ` + "\n" + s)
								err = updateCell(sheet, i, rIndex, s, xlq.OverwriteResult)
								if err != nil {
									xlq.Error("Failed to update cell results for " + searchString + ", " + err.Error())
								} else {
									saveWorkbook = true
								}
							} else {
								println(`No results for "` + searchString + `"`)
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
			err := workbook.Save(xlq.WorkbookName)
			if err != nil {
				xlq.Error("Can't save " + xlq.WorkbookName + ", " + err.Error())
				return errors.New(strings.Join(xlq.Errors, "\n"))
			}
			println("Wrote " + xlq.WorkbookName)
		}
	}
	if len(xlq.Errors) > 0 {
		return errors.New(strings.Join(xlq.Errors, "\n"))
	}
	return nil
}

// WebRunner wrapper for running inside a web browser
