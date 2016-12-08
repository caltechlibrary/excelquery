
+ What additional features do I need from epgo?
    + Can I create an search interface?
    + Do I need to maintain my own index?
    + Do I create a sheet that I can use as an auto-complete source in Excel?
+ JSON Pointer, JSONPath, jq: implementations, references and prior art
    + [RFC-6901](https://tools.ietf.org/html/rfc6901), 2013-04, JavaScript Object Notation (JSON) Pointer
        + https://github.com/dustin/go-jsonpointer, 2014-present, Golang implementation of JSON Pointer
        + https://github.com/xeipuuv/gojsonpointer, 2014-2015, Golang implenentation of JSON Pointer
    + [Goessner's article on JSONPath](http://goessner.net/articles/JsonPath/), 2007-02-21, Stefan Goessner's idea about XPath4JSON
        + https://github.com/FlowCommunications/JSONPath - Flow's versions looks like it is based on origin 2007 and rewritten in modern PHP
            + https://jsonpath.curiousconcept.com/ - claims to be a JSON Path Expression Tester, based on Flow's version
    + https://github.com/dchester/jsonpath, 2015-present?, NodeJS implementaton of Stefan Goessner's JSONPath
    + https://github.com/kennknowles/python-jsonpath-rw, 2013-2015, Python implementation of Stefan Goessner's JSONPath
    + [Tidwall/GJSON](https://github.com/tidwall/gjson) - Tidwall's GJSON provides a varient appproach to jq like dot path language for output
        + doesn't appear to be aware of Stefan Goessner's work
    + [MongoDB's Filter docs](https://docs.mongodb.com/getting-started/shell/query/), 2007-2009?, uses a dotNotation idea for describing query filter/result
    + [stedolan/jq](https://stedolan.github.io/jq) - current work in JSON query/filter on the command line
        + Probably the best of current implementation,
        + Written in C/C++

