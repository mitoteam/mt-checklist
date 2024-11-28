// Main application code

Vue.createApp({
  // delimiters are set only for this component. Each component has it own delimiters.
  delimiters: ['[[', ']]'],

  components: {
    ComponentStatus: ComponentStatus,
    ComponentMtModal: ComponentMtModal,
    ComponentAdminChecklistsList: ComponentAdminChecklistsList,
  },

  data() {
    MtData = Vue.reactive(MtData);
    return MtData;
  },
})
.mount('#app')
