<template>
  <v-layout column>
    <files-browser
      ref="filesBrowser"
      :src-files="src.files"
      :src-path="store.projectDirectory"
      :src-loading="fetchingLocalFiles"
      :dest-path="projectPath"
      :dest-loading="fetchingServerFiles"
      :dest-files="dest"
      :files-progress="uploadProgress"
    >
      <template v-slot:dest-toolbar>
        <v-btn icon @click="fetchServerFiles">
          <v-icon>refresh</v-icon>
        </v-btn>
      </template>
    </files-browser>
    <v-layout class="row toolbar mx-1 my-1">
      <v-spacer/>
      <v-btn rounded @click="uploadFiles">
        <v-icon class="mr-2">cloud_upload</v-icon>
        <span>Upload</span>
      </v-btn>
      <v-spacer/>
    </v-layout>
  </v-layout>
</template>

<script>
import FilesBrowser from '@/components/FilesBrowser'

export default {
  name: 'Upload',
  components: { FilesBrowser },
  data () {
    return {
      fetchingLocalFiles: false,
      fetchingServerFiles: false,
      src: {},
      dest: [],
      uploadProgress: null
    }
  },
  computed: {
    store () {
      return this.$parent
    },
    config () {
      return this.store.projectInfo
    },
    projectPath () {
      return this.store.projectPath
    }
  },
  watch: {
    projectPath: {
      immediate: true,
      handler (path) {
        if (path) {
          this.fetchServerFiles()
        } else {
          this.dest = []
        }
      }
    }
  },
  created () {
    this.$ws.bind('UploadProgress', this.onProgressMessage)
    // this.fetchFiles()
  },
  beforeDestroy () {
    this.$ws.unbind('UploadProgress', this.onProgressMessage)
  },
  activated () {
    this.fetchFiles()
    if (this.projectPath) {
      this.fetchServerFiles()
    }
  },
  methods: {
    fetchFiles () {
      this.fetchingLocalFiles = true
      this.$ws.request('Files')
        .then(resp => {
          this.fetchingLocalFiles = false
          this.src = JSON.parse(resp)
        })
    },
    fetchServerFiles () {
      this.fetchingServerFiles = true
      this.$http.get(`/api/project/files/${this.projectPath}`)
        .then(resp => {
          this.dest = resp.data
          this.fetchingServerFiles = false
        })
        .catch(err => {
          if (err.response && err.response.status !== 404) {
            console.error(err)
          }
          this.dest = []
          this.fetchingServerFiles = false
        })
    },
    uploadFiles () {
      const { newFiles, modifiedFiles } = this.$refs.filesBrowser
      const files = this.src.files.filter(f => newFiles[f.path] || modifiedFiles[f.path])
      const upload = {}
      files.forEach(f => {
        upload[f.path] = {
          size: f.size,
          progress: 0
        }
      })
      this.uploadProgress = upload
      this.$ws.sendJSON('SendFiles', { files, project: this.projectPath })
    },
    onProgressMessage(e, msg) {
      const data = JSON.parse(msg)
      if (!this.uploadProgress) {
        return
      }
      Object.entries(data).forEach(([file, progress]) => {
        const fileUpload = this.uploadProgress[file]
        if (fileUpload) {
          fileUpload.progress = 100 * (progress / fileUpload.size)
        }
      })

      const running = Object.values(this.uploadProgress).some(u => u.progress !== 100)
      if (!running) {
        this.uploadProgress = null
        // this.$ws.sendJSON('SaveProjectConfig', { project: this.projectPath, config: this.config })
        this.$http.post(`/api/project/config/${this.projectPath}`, this.config)
          .then(() => {
            this.fetchServerFiles()
            this.store.projectConfig = {
              ...this.config,
              topics: []
            }
          })
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.browser {
  flex: 1 1;
  overflow: hidden;
}
.toolbar {
  flex: 0 0 auto;
}
</style>
