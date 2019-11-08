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
      :files-progress="filesUploadProgress"
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
      <template v-slot:src-no-content>
        <v-layout fill-height align-center justify-center>
          <v-icon color="red" class="mx-2">error</v-icon>
          <div class="error--text">{{ localFilesError }}</div>
        </v-layout>
      </template>
    </files-browser>
    <div class="toolbar mx-1 my-1">
      <v-spacer/>
      <v-btn
        v-if="!uploadProgress"
        key="upload"
        :disabled="!canUpload"
        @click="uploadFiles"
        rounded
      >
        <v-icon class="mr-2">cloud_upload</v-icon>
        <span>Upload</span>
      </v-btn>
      <v-btn
        v-else
        key="cancel"
        @click="upload.abort()"
        rounded
      >
        <v-icon class="mr-2">cancel</v-icon>
        <span>Cancel</span>
        <v-progress-linear
          height="3"
          :value="uploadProgress.totalProgress"
        />
      </v-btn>

      <v-btn rounded :href="`/api/project/download/${projectPath}`" class="self-align-end">
        <v-icon class="mr-2">archive</v-icon>
        <span>Download</span>
      </v-btn>
    </div>
  </v-layout>
</template>

<script>
import { createUpload } from '@/upload.js'
import FilesBrowser from '@/components/FilesBrowser'

export default {
  name: 'Files',
  components: { FilesBrowser },
  refs: ['filesBrowser'],
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
      localFilesError: '',
      loadingServerFiles: false,
      loadingLocalFiles: false,
      uploadProgress: null
    }
  },
  computed: {
    pluginConnected () {
      return this.$ws.pluginConnected
    },
    projectPath () {
      return this.user && this.folder && `${this.user}/${this.folder}`
    },
    browser () {
      return this.$refs.filesBrowser
    },
    canUpload () {
      if (this.browser && !this.loadingLocalFiles && !this.loadingServerFiles) {
        const { newFiles, modifiedFiles } = this.browser
        return newFiles && modifiedFiles && newFiles.length + modifiedFiles.length > 0
      }
      return false
    },
    filesUploadProgress () {
      return this.uploadProgress && this.uploadProgress.files
    }
  },
  activated () {
    const unbind = this.$ws.bind('ProjectChanged', this.fetchLocalFiles)
    const unwatch = this.$watch('pluginConnected', connected => {
      if (connected && !this.loadingLocalFiles) {
        this.fetchLocalFiles()
      }
    })
    this.$once('hook:deactivated', unbind)
    this.$once('hook:deactivated', unwatch)
    this.fetchLocalFiles()
  },
  watch: {
    projectPath: {
      immediate: true,
      handler (path) {
        if (path) {
          this.fetchServerFiles()
        } else {
          this.serverFiles = []
        }
      }
    }
  },
  methods: {
    fetchLocalFiles () {
      if (this.$ws.pluginConnected) {
        this.loadingLocalFiles = true
        this.localFilesError = ''
        this.$ws.request('ProjectFiles')
          .then(resp => {
            this.loadingLocalFiles = false
            const data = resp.data
            this.localFiles = data.files
            this.localDirectory = data.directory
          })
          .catch(err => {
            this.loadingLocalFiles = false
            this.localFiles = []
            this.localDirectory = ''
            this.localFilesError = err.data || 'Faild to list project files'
          })
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
    async uploadFiles () {
      const { newFiles, modifiedFiles } = this.$refs.filesBrowser
      const files = [...newFiles, ...modifiedFiles]
      // const files = newFiles

      this.upload = createUpload(this.$ws, files, this.projectPath)
      this.uploadProgress = this.upload.info
      try {
        await this.upload.start()
      } catch (e) {
        if (e !== 'aborted') {
          console.error(e)
        }
      } finally {
        this.fetchServerFiles()
        this.uploadProgress = null
      }
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
    .v-btn {
      .v-progress-linear {
        position: absolute;
        width: 100%;
        bottom: -3px;
      }
    }
  }
}
</style>
