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
  props: {
    checklists: [],
  },

  template: `
<div>
  Some content
</div>`
}
//#endregion
