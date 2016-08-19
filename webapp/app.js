(function (doc, win) {
    "use strict";
    var eprintsSearchURL = document.getElementById("eprintsSearchURL"),
        dataPath = document.getElementById("dataPath"),
        overwriteResult = document.getElementById("overwriteResult"),
        sheetName = document.getElementById("sheetName"),
        workbook = document.getElementById("workbook"),
        queryColumn = document.getElementById("queryColumn"),
        resultColumn = document.getElementById("resultColumn"),
        data;

  workbook.addEventListener("change", function (evt)  {
      var fp = workbook.files[0],
          xlq = xlquery.New(),
          xlr = xlresponse.New(),
          files = this.files,
          reader = new FileReader();

      xlq.EPrintsSearchURL = eprintsSearchURL.value;
      xlq.DataPath = dataPath.value;
      xlq.SheetName = sheetName.value;
      xlq.Overwrite = overwriteResult.value;
      
      console.log("DEBUG evt.target.files", evt.target.files[0]);
      reader.onload = function (eFile) {
          console.log("DEBUG eFile.target.result (data url)", eFile.target.result);
          //FIXME: Pass the data url to xlq.Run()
      };
      reader.readAsDataURL(files[0]);

      /*
      console.log("DEBUG queryColumn.value", queryColumn.value);
      console.log("DEBUG resultColumn.value", resultColumn.value);
      xlr = xlq.Run(data, queryColumn.value, resultColumn.value);
      if (xlr.Errors.length > 0) {
          console.log("DEBUG xlr.Errors", xlr.Errors);
      }
      console.log("DEBUG xlr.Data", xlr.Data);
      */
  }, false)

}(document, window))
