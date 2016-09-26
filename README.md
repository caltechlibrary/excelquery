
# xlquery

xlquery is an experimental utility designed to update workbook files (.xlsx) taking a column as search query strings and 
putting the results in another column of the same spreadsheet. This proof of concept works with EPrints repository 
software relying on advanced search and results returned in RSS2 format. By default the search if performanced against
[CaltechAUTHORS](https://authors.library.caltech.edu) repository. You can point at a different EPrints repository by
setting the environment variable *EPRINTS_SEARCH_URL*.

## USAGE

```shell
    xlquery [OPTIONS] WORKBOOK_NAME QUERY_SHEET_NAME QUERY_COLUMN
```

The command line program *xlquery* takes the name of a xlsx file along with a sheet name (or number) and the column name for the query string. It also accepts an optional sheeting name of the results (otherwise a new sheet w
and one for the result. By default it searches on "Sheet1". This can be changed with the "-s" command line option.

The simple form where column *A* in *Sheet 1" holds the query string and results will be put in a new sheet called "Results" 

```shell
    xlquery titlelist.xlsx "Sheet 1" A 
```

*xlquery* will display console message describing the processing on stdout. If there are errors they will be sent to 
stderr with catastrophic errors exiting with a value 1. If the program is successful it will exit with the value 0.

+ The sheet name can be the textual number of the sheet or its index (the first sheet's index is zero)
+ Query column should correspond to the sheet you want to run through (e.g. "Sheet1")
+ Column names are in Excel's letter format (e.g. "A", "FX", "BBC").
+ data path is the part of the result you want to use (e.g. ".item[].link" is the RSS item link field)

## OPTIONS

    -d, -datapath	set the data path to return for results, e.g. ".item[].link"
    -h, -help	show help information
    -l, -license	show license information
    -o, -overwrite	overwrite the results column
    -s, -sheet	set the sheet name, e.g. "Sheet1"
    -v, -version	show version information


## Example

```shell
    xlquery -sheet "Title List" -overwrite ./testdata/demo2.xlsx A
```

This opens demo2.xlsx and uses the sheet named "Title List". It populates fresh results in in a new sheet called "Results" based on the 
query string in column *A* of "Title List". The results are taken from the data path of ".item[].link" from the RSS2 
response in the search request.


