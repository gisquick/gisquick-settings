<template>
  <v-layout align-end>
    <input
      ref="input"
      type="file"
      @input="onChange"
      @change="onChange"
      hidden
    >
    <v-tooltip bottom>
      <template v-slot:activator="{ on }">
        <v-btn
          v-on="on"
          rounded text
          @click="$refs.input.click()"
        >
          <v-icon class="mr-1">cloud_upload</v-icon>
          <span>Upload</span>
        </v-btn>
      </template>
      <span>Upload Zip archive created with plugin </span>
    </v-tooltip>
    <v-tooltip bottom>
      <template v-slot:activator="{ on }">
        <v-btn
          v-on="on"
          rounded text
          active-class="lime--text text--darken-2"
          :to="{name: 'publish'}"
        >
          <v-icon class="mr-1">$vuetify.icons.qgis</v-icon>
          <span>Publish</span>
        </v-btn>
      </template>
      <span>Publish directly from QGIS</span>
    </v-tooltip>
    <v-progress-circular
      v-if="uploading"
      :value="progress"
      color="lime darken-2"
      size="28"
    />
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
        this.$emit('upload-finished')
        this.$notification.show('Project uploaded')
      }).catch(err => {
        this.uploading = false
        const msg = (err && err.response && err.response.data) || 'Error'
        this.$notification.error(msg)
      })
    }
  }
}
</script>
