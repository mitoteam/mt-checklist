let ComponentStatus = {
  props: {
    body: String,
    kind: {
      type: String,
      default: "primary",
    },
  },

  template: `
<div v-if="body" class="mb-3 alert" :class="['alert-' + kind]" role="alert">
  {{body}}
</div>`
}

//#region Admin area Components
let ComponentAdminChecklistsList = {
  data: function() {
    return {
      checklists: [],
    }
  },

  mounted: function() {
    console.log("the component is now mounted");
    this.checklists = [{name: "asdf"}, {name: "BBB"}];
    console.log(this.checklists);
  },

  template: `
<div>
  Some content
  <ul>
    <li v-for="checklist in checklists">
      {{ checklist.name }}
    </li>
  </ul>
</div>`,
}
//#endregion
