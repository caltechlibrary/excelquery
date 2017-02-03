
# USAGE

    excelquery XLSX_FILENAME SHEET_NO ROW_NO COLUMN_NO

## SYNOPSIS

excelquery query our repositories for matching information.

## OPTIONS

```
	-h	show help information
	-help	show help information
	-l	show license information
	-license	show license information
	-s	set boolean for skipping first row of sheet (default true)
	-skip	set boolean for skipping first row of spreadsheet (default true)
	-v	show version information
	-version	show version information
```

## EXAMPLE

```
	excelquery inventory.xlsx 0 20 20
```

Query sheet number 0 (the first sheet) based on row 20,
column 20.

