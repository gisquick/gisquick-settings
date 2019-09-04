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
      active-class="orange--text"
      :to="{name: 'layers'}"
    >
      <v-icon class="mr-2">layers</v-icon>
      <span>Layers</span>
    </v-btn>
    <v-btn
      rounded text
      class="mx-1"
      active-class="orange--text"
      :to="{name: 'topics'}"
    >
      <icon name="topics" class="mr-2"/>
      <span>Topics</span>
    </v-btn>
    <v-btn
      rounded text
      class="mx-1"
      active-class="orange--text"
      :to="{name: 'files'}"
    >
      <v-icon class="mr-2">cloud_upload</v-icon>
      <span>Files</span>
    </v-btn>
    <v-divider vertical class="mt-2"/>
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
