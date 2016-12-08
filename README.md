
# xlquery

xlquery is an experimental utility designed to update workbook files (.xlsx) taking a column as search query strings and 
putting the results in another worksheet of the same workbook. This proof of concept works with EPrints repository 
software relying on advanced search and results returned in RSS2 format. By default the search if performanced against
[CaltechAUTHORS](https://authors.library.caltech.edu) repository. You can point at a different EPrints repository by
setting the environment variable *EPRINTS_SEARCH_URL*.

## USAGE

```shell
    xlquery [OPTIONS] WORKBOOK_NAME QUERY_SHEET_NAME QUERY_COLUMN [RESULT_SHEET_NAME]
```

The command line program *xlquery* takes the name of a xlsx file along with a sheet name (or number) and the column name 
for the query string. A sheet name can optionally be supplied for results.  By default it searches on "Sheet1" and by 
default a new sheet is created called "Result". This can be changed with the "-s" and "-r" command line options.

The simple form where column *A* in *Sheet 1" holds the query string and results will be put in a new sheet called "Result" 

```shell
    xlquery titlelist.xlsx "Sheet 1" A 
```

*xlquery* will display console message describing the processing on stdout. If there are errors they will be sent to 
stderr with catastrophic errors exiting with a value 1. If the program is successful it will exit with the value 0.

+ The sheet name can be the textual number of the sheet or its index (the first sheet's index is zero)
+ Query column should correspond to the sheet you want to run through (e.g. "Sheet1")
+ Column names are in Excel's letter format (e.g. "A", "FX", "BBC").

## OPTIONS

+ -h, -help   show help information
+ -l, -license    show license information
+ -v, -version    show version information

+ -s, -skip   set boolean for skipping first row of spreadsheet (default true)


## Example

```shell
    xlquery ./testdata/demo2.xlsx "Title List" A
```

This opens demo2.xlsx and uses the sheet named "Title List". It populates fresh results in in a new sheet called "Result" based on the 
query string in column *A* of "Title List". The results are taken from the item field of the RSS2 
response to the search request.


