<template>
  <v-layout align-end>
    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title class="py-2 headline grey lighten-2">Are you sure?</v-card-title>
        <v-divider color="grey"/>
        <v-card-actions>
          <v-spacer/>
          <v-btn text @click="deleteDialog = false">
            Cancel
          </v-btn>
          <v-btn
            text
            color="deep-orange"
            @click="deleteProject"
          >
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-btn
      rounded text
      class="mx-1"
      target="_blank"
      :href="`/?PROJECT=${user}/${folder}/${projectName}`"
    >
      <v-icon class="mr-2">map</v-icon>
      <span>Map</span>
    </v-btn>
    <v-divider vertical class="my-1"/>
    <v-btn
      rounded text
      class="mx-1 deep-orange--text"
      @click="deleteDialog = true"
    >
      <v-icon class="mr-2">delete_forever</v-icon>
      <span>Delete</span>
    </v-btn>
  </v-layout>
</template>

<script>
export default {
  props: {
    user: String,
    folder: String,
    projectName: String
  },
  data () {
    return {
      deleteDialog: false
    }
  },
  methods: {
    deleteProject () {
      this.$http.delete(`/api/project/delete/${this.user}/${this.folder}`)
        .then(resp => {
          this.deleteDialog = false
          this.$router.push({ name: 'projects' })
        })
    }
  }
}
</script>

<style lang="scss" scoped>
.v-btn--text {
  &:before {
    display: none;
  }
}
.theme--dark .icon {
  color: inherit;
}
.r-90 {
  transform: rotate(90deg);
}
</style>
