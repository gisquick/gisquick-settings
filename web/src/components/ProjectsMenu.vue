<template>
  <v-layout align-end>
    <input
      ref="input"
      type="file"
      @input="onChange"
      @change="onChange"
      hidden
    >
    <v-btn
      small rounded text
      class="mx-1"
      color="lime darken-2"
      @click="$refs.input.click()"
    >
      <v-icon class="mr-2">cloud_upload</v-icon>
      <span>New project</span>
    </v-btn>
    <v-progress-circular v-if="uploading" :value="progress" color="lime darken-2" size="28"/>
  </v-layout>
</template>

<script>
export default {
  data () {
    return {
      uploading: false,
      progress: 0
    }
  },
  methods: {
    onChange (evt) {
      const file = evt.target.files[0]
      evt.target.value = ''
      if (!file) {
        return
      }
      const form = new FormData()
      form.append(file.name, file, file.name)
      this.progress = 0
      this.uploading = true
      this.$http.post('/api/project/upload', form, {
        onUploadProgress: evt => {
          this.progress = 100 * (evt.loaded / evt.total)
        }
      }).then(resp => {
        this.uploading = false
        this.$http.get('/projects.json').then(resp => {
          this.$root.projects = resp.data.projects
        })
      }).catch(err => {
        this.uploading = false
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
</style>
