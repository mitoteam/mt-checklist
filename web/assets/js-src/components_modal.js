let ComponentMtModal = {
  props: {
    title: String,
  },

  template: `
<div class="modal fade" id="modal" data-bs-backdrop="static" data-bs-keyboard="true" tabindex="-1" aria-hidden="true">
  <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable modal-lg modal-fullscreen-md-down">
    <div class="modal-content">
      <div class="modal-header">
        <h1 class="modal-title fs-5">{{ title }}</h1>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
      </div>
      <div class="modal-body text-prewrap">
        Yo-Ho!
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
      </div>
    </div>
  </div>
</div>
`
}

//Bootstrap's modal dialog helper
//JS "singleton" pattern
const MtModal = (function() {
  //// PRIVATE
  let initialized = false;
  let bsModal = null;
  let titleElement = document.querySelector("#modal .modal-title");
  let bodyElement = document.querySelector("#modal .modal-body");

  function init()
  {
    if(!initialized)
    {
      initialized = true;

      bsModal = bootstrap.Modal.getOrCreateInstance('#modal');
      titleElement = document.querySelector("#modal .modal-title");
      bodyElement = document.querySelector("#modal .modal-body");
    }
  };

  //// PUBLIC
  return {
    ShowHtml: function (title, body_html)
    {
      init();

      titleElement.textContent = title;
      bodyElement.innerHTML = body_html;

      bsModal.show();
    },

    ShowGetHtml: function (title, url)
    {
      init();

      titleElement.textContent = title;
      bsModal.show();

      MtTools.SimpleGet(bodyElement, url);
    },

    Hide: function ()
    {
      bsModal.hide();
    },

    GetBodyEl: function()
    {
      return bodyElement;
    }
  };
})();
