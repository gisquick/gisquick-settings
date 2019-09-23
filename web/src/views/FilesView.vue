<template>
  <v-layout class="column files-page fill-height">
    <files-browser
      ref="filesBrowser"
      :src-files="localFiles"
      :src-path="localDirectory"
      :src-loading="loadingLocalFiles"
      :dest-path="projectPath"
      :dest-loading="loadingServerFiles"
      :dest-files="serverFiles"
      :files-progress="upload"
      class="mb-1"
    >
      <template v-slot:src-toolbar>
        <v-btn icon @click="fetchLocalFiles">
          <v-icon>refresh</v-icon>
        </v-btn>
      </template>
      <template v-slot:dest-toolbar>
        <v-btn icon @click="fetchServerFiles">
          <v-icon>refresh</v-icon>
        </v-btn>
      </template>
    </files-browser>
    <div class="toolbar mx-1 my-1">
      <v-spacer/>
      <v-btn
        rounded
        :disabled="!canUpload"
        @click="startUpload"
      >
        <v-icon class="mr-2">cloud_upload</v-icon>
        <span>Upload</span>
      </v-btn>
      <v-btn rounded :href="`/api/project/download/${projectPath}`" class="self-align-end">
        <v-icon class="mr-2">archive</v-icon>
        <span>Download</span>
      </v-btn>
    </div>
  </v-layout>
</template>

<script>
import FilesBrowser from '@/components/FilesBrowser'

export default {
  name: 'Files',
  components: { FilesBrowser },
  props: {
    user: String,
    folder: String,
    projectName: String
  },
  data () {
    return {
      localDirectory: '',
      localFiles: [],
      serverFiles: [],
      loadingServerFiles: false,
      loadingLocalFiles: false,
      upload: null
    }
  },
  computed: {
    projectPath () {
      return this.user && this.folder && `${this.user}/${this.folder}`
    },
    canUpload () {
      return !this.loadingLocalFiles && !this.loadingServerFiles && this.localFiles.length > 0
    }
  },
  watch: {
    projectPath: {
      immediate: true,
      handler (path) {
        if (path) {
          this.fetchServerFiles()
          this.fetchLocalFiles()
        } else {
          this.serverFiles = []
        }
      }
    }
  },
  created () {
    this.$ws.bind('Files', this.onFilesMessage)
    this.$ws.bind('UploadProgress', this.onProgressMessage)
  },
  beforeDestroy () {
    this.$ws.unbind('Files', this.onFilesMessage)
    this.$ws.unbind('UploadProgress', this.onProgressMessage)
  },
  methods: {
    onFilesMessage (e, msg) {
      this.loadingLocalFiles = false
      const data = JSON.parse(msg)
      this.localFiles = data.files
      this.localDirectory = data.directory
    },
    onProgressMessage(e, msg) {
      const data = JSON.parse(msg)
      if (!this.upload) {
        return
      }
      Object.entries(data).forEach(([file, progress]) => {
        const fileUpload = this.upload[file]
        if (fileUpload) {
          fileUpload.progress = 100 * (progress / fileUpload.size)
        }
      })

      const running = Object.values(this.upload).some(u => u.progress !== 100)
      if (!running) {
        this.upload = null
        this.fetchServerFiles()
      }
    },
    fetchLocalFiles () {
      if (this.$ws.pluginConnected) {
        this.loadingLocalFiles = true
        this.$ws.send('Files')
      }
    },
    fetchServerFiles () {
      this.loadingServerFiles = true
      this.$http.get(`/api/project/files/${this.projectPath}`)
        .then(resp => {
          this.serverFiles = resp.data
          this.loadingServerFiles = false
        })
        .catch(err => {
          console.error(err)
          this.loadingServerFiles = false
        })
    },
    startUpload () {
      const { newFiles, modifiedFiles } = this.$refs.filesBrowser
      // const files = this.localFiles.filter(f => newFiles[f.path] || modifiedFiles[f.path])
      const files = this.localFiles.filter(f => newFiles[f.path])
      const upload = {}
      files.forEach(f => {
        upload[f.path] = {
          size: f.size,
          progress: 0
        }
      })
      this.upload = upload
      this.$ws.sendJSON('SendFiles', { files, project: this.projectPath })
    }
  }
}
</script>
<style lang="scss" scoped>
.browser {
  overflow: hidden;
}
.files-page {
  .box {
    border: 1px solid #aaa;
    overflow: hidden;
  }
  .toolbar {
    grid-column: 1 / 3;
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    > * {
      justify-self: center;
      &:last-child {
        justify-self: end;
      }
    }
  }
}
</style>
