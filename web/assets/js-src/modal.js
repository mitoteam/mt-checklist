//Bootstrap's modal dialog helper

//JS "singleton" pattern
const MtModal = (function() {
  //// PRIVATE
  const bsModal = new bootstrap.Modal('#modal');
  const titleElement = document.querySelector("#modal .modal-title");
  const bodyElement = document.querySelector("#modal .modal-body");

  //// PUBLIC
  return {
    ShowHtml: function (title, body_html)
    {
      titleElement.textContent = title;
      bodyElement.innerHTML = body_html;

      bsModal.show();
    },

    Hide: function ()
    {
      bsModal.hide();
    },

    GetBodyEl: function()
    {
      return bodyElement
    }
  };
})();
