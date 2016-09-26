
# Ideas about workflow

+ xlquery prototype
    + Enter titles and partial titles to query EPrints
    + use *xlquery* to return a list of links per query term
    + work that spreadsheet down until you have very few links per row
    + For each row with a sink link produce a BibTeX entry 
    + Export the BibTeX entries into a spreadsheet (fields become columns)
+ theoretical workflow 1 (multi-tool)
    + type query strings into a column of a spreadsheet
    + run through xlquery
    + cleanup the url lists produced by xlquery
    + run through xlurl2bib to produce a BibTeX file
    + currate as needed
    + run through xlbib2xl to produce a revised spreadsheet with Bib info
+ theoretical workflow 2 (single tools, uses separate sheet for results)
    + type query strings into a column in a spreadsheet
    + run through xlquery 
        + this could store multi-value results as a separate worksheet, on column per key field
        + results would need to suppor expressioning a complex object result (like [jq](https://github.com/stedolan/jq))
+ theoretical workflow 3 (two tools)
    + type in query strings into column in a spreadsheet
    + run through xlquery to generate a url results column
    + run through xlquery with the -bib-sheet option generating multi-valued results in a new sheet



