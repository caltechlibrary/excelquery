
# xlquery

xlquery is an experiment with workbook files (.xlsx) taking a column as search query strings and 
putting the results in a related column of the spreadsheet.

Inputs

+ Workbook filename
+ column for query string
+ column for results
+ data path to result value
+ URL for service or filename to scan

Outputs

+ revised Workbook filename with the results column populated where data was found

## Demo

```shell
    make
    ./bin/xlquery -overwrite ./testdata/demo2.xlsx Sheet1 A C ".item[].link"
```

This opens demo2.xlsx and populates fresh results in column C based on the data path of ".item[].link" from the RSS2 response
in the search request.



