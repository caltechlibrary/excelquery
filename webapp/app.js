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
        resultBlock = document.getElementById("xlrBlock"),
        dataURL = "",
        xlq = xlquery.New();

    
    workbook.addEventListener("change", function (evt)  {
      var fp = evt.target.files[0] || {},
          files = this.files,
          reader = new FileReader();

      xlq.WorkbookName = fp.fileName || ""; 
      console.log("DEBUG evt.target.files", evt.target.files[0]);
      reader.onload = function (eFile) {
          console.log("DEBUG eFile.target.result (data url)", eFile.target.result.substring(0, 64));
          dataURL = eFile.target.result;
      };
      reader.readAsDataURL(files[0]);
    }, false)

    runButton.addEventListener("click", function (evt) {
        var isAlpha = new RegExp("^[a-z,A-Z]+$", "g"),
            isDataURL = new RegExp("^data\:application/","i");

        console.log("DEBUG runButton clicked");

        /* Validate queryColumn */
        if (!queryColumn.value.match(isAlpha)) {
            console.log("ERROR: queryColumn show be in the for A, AA, ABC", queryColumn);
            return;
        }
        xlq.QueryColumn = queryColumn.value.trim();

        /* Validate resultColumn */
        if (!resultColumn.value.match(isAlpha)) {
            console.log("ERROR: resultColumn show be in the for A, AA, ABC", resultColumn);
            return;
        }
        xlq.ResultColumn = resultColumn.value.trim();

        /* Validate overwriteResult */
        if (overwriteResult.checked === true) {
            xlq.OverwriteResult = true;
        } else {
            xlq.OverwriteResult = false;
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
        console.log("DEBUG got all the way to call xlquery.Run()");
        dataURL = xlq.WebRunner(dataURL);
        console.log("DEBUG output dataURL: ", dataURL);
        xlrBlock.innerHTML = "<pre>"+dataURL+"</pre>";
    }, false);
}(document, window))
