//JS "singleton" pattern
const MtApi = (function() {
  //// PRIVATE
  /**
  * Basic request code.
  * @param path relative URL
  * @param data
  * @param component context for responseHandler
  * @param responseHandler function to call after response
  */
  function ApiRequest(path, data, component, responseHandler) {
    fetch(
      '/api/' + path,
      {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data ?? {})
      }
    )
      .then(response => response.json())
      .then(function (response) {
        if (response.status == "ok") {
          MtData.status.status = "success";
        }
        else {
          MtData.status.status = response.status;
        }

        MtData.status.message = response.message ?? "";

        if (typeof (responseHandler) == 'function') {
          responseHandler.call(component, response);
        }
      });
  }


  //// PUBLIC
  return {
    AdminGetChecklists: function (component)
    {
      ApiRequest('admin/checklists', null, component, function (response) {
        //console.log(response);
        if (response.status == "ok")
        {
          component.checklists = response.checklists
        }
      });
    }
  };
})();
