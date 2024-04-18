//#region Admin area Components
let ComponentAdminChecklistsList = {
  data: function() {
    return {
      checklists: [],
    }
  },

  mounted: function() {
    console.log("the component is now mounted");

    MtApi.AdminGetChecklists(this);
  },

  template: `
<div class="mb-3 p-3 border">
  <a href="#" class="btn btn-primary" onclick="MtModal.ShowHtml('title', 'htm<b>l</b>');">
    <i class="far fa-plus"></i> New checklist
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
