
# xlquery

xlquery is an experimental utility designed to update workbook files (.xlsx) taking a column as search query strings and 
putting the results in another column of the same spreadsheet. This proof of concept works with EPrints repository 
software relying on advanced search and results returned in RSS2 format. By default the search if performanced against
[CaltechAUTHORS](https://authors.library.caltech.edu) repository. You can point at a different EPrints repository by
setting the environment variable *EPRINTS_SEARCH_URL*.

## USAGE

```shell
    xlquery [OPTIONS] WORKBOOK_NAME QUERY_COLUMN RESULT_COLUMN
```

The command line program *xlquery* takes the name of a xlsx file along with the column name for the query string
and one for the result. By default it searches on "Sheet1". This can be changed with the "-s" command line option.

The simple form where column *A* holds the query string and results will be put in column *B*

```shell
    xlquery titlelist.xlsx A B
```

*xlquery* will display console message describing the processing on stdout. If there are errors they will be sent to 
stderr with catastrophic errors exiting with a value 1. If the program is successful it will exit with the value 0.

+ Sheet name should correspond to the sheet you want to run through (e.g. "Sheet1")
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
    xlquery -sheet "Title List" -overwrite ./testdata/demo2.xlsx A C 
```

This opens demo2.xlsx and uses the sheet named "Title List". It populates fresh results in column *C* based on the 
query string in column *A*. The results are taken from the data path of ".item[].link" from the RSS2 
response in the search request.


