
Results Fields:

+ required
    + title, author, year, 
+ optional
    + eprints id, persistant url, 


Feature adds:

+ xlquery options
    + set start row
    + set end row
+ xlquery result objects
    + populate more than one column with results (still split by line if multiple results returned)
    + specify the fields you want returned as a compound object

Performance: 
    + send Betsy a script performing the queries and she will measure the performance hit

Someday:
    + add mult, column response (e.g. '{.item[].link, .item[].title}')

