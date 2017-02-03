package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"

	// 3rd Party Go packages
	"github.com/tealeg/xlsx"
)

var (
	usage = `USAGE: %s XLSX_FILENAME SHEET_NO ROW_NO COLUMN_NO`

	description = `

%s check the contents of a cell
`

	examples = `
EXAMPLE

	%s inventory.xlsx 0 20 20

Show the contents of inventory.xlsx, sheet number 0 (the first sheet) 
row 20 and column 20.
`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool
)

func init() {
	// Standard Options
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()
	args := flag.Args()

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
