package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	// 3rd Party Go packages
	"github.com/tealeg/xlsx"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 4 {
		fmt.Fprintf(os.Stderr, "USAGE: t XLSX_FILENAME SHEET_NO ROW_NO COLUMN_NO\n")
		os.Exit(1)
	}
	i := 0
	fname := args[0]
	i++
	sheetNo, err := strconv.Atoi(args[i])
	if err != nil {
		fmt.Fprintf(os.Stderr, "SHEET_NO %s should be a number, %s\n", args[i], err)
		os.Exit(1)
	}
	i++
	row, err := strconv.Atoi(args[i])
	if err != nil {
		fmt.Fprintf(os.Stderr, "row %s should be a number, %s\n", args[i], err)
		os.Exit(1)
	}
	i++
	column, err := strconv.Atoi(args[i])
	if err != nil {
		fmt.Fprintf(os.Stderr, "column %s should be a number, %s\n", args[i], err)
		os.Exit(1)
	}
	i++

	workbook, err := xlsx.OpenFile(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open workbook %s, %s\n", fname, err)
		os.Exit(1)
	}
	if len(workbook.Sheets) <= sheetNo {
		fmt.Fprintf(os.Stderr, "Can't find sheet %s in workbook %s\n", sheetNo, fname)
		os.Exit(1)
	}

	sheet := workbook.Sheets[sheetNo]

	// Now write the cell data out!
	src, err := json.Marshal(sheet.Cell(row, column))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't marshal cell, %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", src)
}
