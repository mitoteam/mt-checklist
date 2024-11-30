//#region Admin area Components
let ComponentAdminChecklistsList = {
  data: function() {
    return {
      checklists: [],
    }
  },

  mounted: function() {
    //console.log("the component is now mounted");

    MtApi.AdminGetChecklists(this);
  },

  template: `
<div class="mb-3 p-3 border">
  <a href="#" class="btn btn-primary" onclick="MtModal.ShowGetHtml('dbg', '/experiment');">
    <i class="far fa-plus"></i> New checklist
  </a>
  <a href="#" class="btn btn-secondary ms-1" onclick="MtModal.ShowHtml('test title', 'test body <b>html</b>');">
    <i class="far fa-flask-vial"></i> test btn
  </a>
</div>
<div>
  <ul>
    <li v-for="checklist in checklists">
      {{ checklist.name }}
    </li>
  </ul>
</div>`,
}
//#endregion
