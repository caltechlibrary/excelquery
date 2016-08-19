(function (doc, win) {
    "use strict";
    var eprintsSearchURL = document.getElementById("eprintsSearchURL"),
        dataPath = document.getElementById("dataPath"),
        overwriteResult = document.getElementById("overwriteResult"),
        sheetName = document.getElementById("sheetName"),
        workbook = document.getElementById("workbook"),
        queryColumn = document.getElementById("queryColumn"),
        resultColumn = document.getElementById("resultColumn"),
        runButton = document.getElementById("xlqRun"),
        dataURL = "";

    workbook.addEventListener("change", function (evt)  {
      var fp = workbook.files[0],
          files = this.files,
          reader = new FileReader();

      console.log("DEBUG evt.target.files", evt.target.files[0]);
      reader.onload = function (eFile) {
          console.log("DEBUG eFile.target.result (data url)", eFile.target.result.substring(0, 64));
          dataURL = eFile.target.result;
      };
      reader.readAsDataURL(files[0]);
    }, false)

    runButton.addEventListener("click", function (evt) {
        var isAlpha = new RegExp("^[a-z,A-Z]+$", "g"),
            isDataURL = new RegExp("^data\:application/vnd\.openxmlformats-officedocument\.spreadsheetml\.sheet;base64,","i"),
            xlr = {},
            xlq = {};

        console.log("DEBUG runButton clicked");
        /* Instantiate a XLQuery object */
        xlq = xlquery.New();

        /* Validate queryColumn */
        if (!queryColumn.value.match(isAlpha)) {
            console.log("ERROR: queryColumn show be in the for A, AA, ABC", queryColumn);
            return;
        }

        /* Validate resultColumn */
        if (!resultColumn.value.match(isAlpha)) {
            console.log("ERROR: resultColumn show be in the for A, AA, ABC", resultColumn);
            return;
        }

        /* Validate overwriteResult */
        if (overwriteResult.checked === true) {
            xlq.Overwrite = true;
        } else {
            xlq.Overwrite = false;
        }

        /* Validate sheetName */
        if (sheetName.value.trim() === "") {
            xlq.SheetName = "Sheet1";
        } else {
            xlq.SheetName = sheetName.value.trim();
        }

        /* Validate dataPath */
        if (dataPath.value.trim() === "") {
            xlq.DataPath = ".item[].link";
        } else {
            xlq.DataPath = dataPath.value.trim();
        }

        /* Validate dataURL */
        if (!dataURL.match(isDataURL)) {
            console.log("ERROR: not a dataURL ["+dataURL+"]");
            return;
        }
        console.log("DEBUG got all the way to call xlq.Run()");
        xlr = xlq.Run(dataURL, queryColumn.value, resultColumn.value);
        console.log("DEBUG xlr: ", xlr);
       
        // xlr = xlq.Run()...
        // if not errors create an embedded dataURL for download
    }, false);
}(document, window))
