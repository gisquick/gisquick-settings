<template>
  <v-layout column>
    <files-browser
      ref="filesBrowser"
      :src-files="src.files"
      :src-path="store.projectDirectory"
      :src-loading="fetchingLocalFiles"
      :dest-path="projectServerDir"
      :dest-loading="fetchingServerFiles"
      :dest-files="dest"
      :files-progress="filesUploadProgress"
    >
      <template v-slot:dest-toolbar>
        <v-btn icon @click="fetchServerFiles">
          <v-icon>refresh</v-icon>
        </v-btn>
      </template>
    </files-browser>
    <div class="toolbar mx-1 mt-2 mb-1">
      <div
        v-if="browser"
        class="left"
      >
        <template v-if="!fetchingLocalFiles && !fetchingServerFiles">
          <template v-if="browser.newFiles.length + browser.modifiedFiles.length > 0">
            <small
              v-show="browser.newFiles.length"
              class="mx-2 green--text"
            >
              New files: {{ browser.newFiles.length }}
            </small>
            <small
              v-show="browser.modifiedFiles.length"
              class="mx-2 orange--text"
            >
              Changed files: {{ browser.modifiedFiles.length }}
            </small>
          </template>
          <small v-else>No changes detected</small>
        </template>
      </div>
      <v-btn
        v-if="!uploadProgress"
        key="upload"
        :disabled="fetchingLocalFiles || fetchingServerFiles"
        @click="uploadFiles"
        rounded
      >
        <v-icon class="mr-2">cloud_upload</v-icon>
        <span>{{ dest.length === 0 ? 'Upload' : 'Update' }}</span>
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
      <v-spacer/>
    </div>
  </v-layout>
</template>

<script>
import _omit from 'lodash/omit'
import { dirname } from 'path'

import FilesBrowser from '@/components/FilesBrowser'
import { mapLayers, layersList } from '@/utils.js'
import { createUpload } from '@/upload.js'

export default {
  name: 'Upload',
  components: { FilesBrowser },
  refs: ['filesBrowser'],
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
    },
    projectServerDir () {
      return dirname(this.projectPath)
    },
    browser () {
      return this.$refs.filesBrowser
    },
    filesUploadProgress () {
      return this.uploadProgress && this.uploadProgress.files
    }
  },
  watch: {
    projectServerDir: {
      immediate: true,
      handler (path) {
        if (path) {
          // this.fetchServerFiles()
        } else {
          this.dest = []
        }
      }
    },
    dest (files) {
      this.store.serverFiles = files
    }
  },
  activated () {
    this.fetchFiles()
    if (this.projectServerDir) {
      this.fetchServerFiles()
    }
  },
  methods: {
    fetchFiles () {
      this.fetchingLocalFiles = true
      this.$ws.request('ProjectFiles')
        .then(resp => {
          this.fetchingLocalFiles = false
          this.src = resp.data
        })
        .catch(err => {
          this.fetchingLocalFiles = false
          this.$notification.error(err.data || 'Error')
        })
    },
    fetchServerFiles () {
      this.fetchingServerFiles = true
      this.$http.get(`/api/project/files/${this.projectServerDir}`)
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
    saveConfig () {
      this.$http.post(`/api/project/config/${this.projectPath}`, this.config)
        .then(() => {
          this.fetchServerFiles()

          const layers = JSON.parse(JSON.stringify(this.config.layers))
          layersList(layers).forEach(l => {
            delete l.wfs
            delete l.source
            l.publish = true
          })
          this.store.projectConfig = {
            ..._omit(this.config, ['layers', 'file', 'directory']),
            base_layers: [],
            // overlays: mapLayers(layers, l => _omit(l, ['wfs', 'source'])),
            overlays: layers,
            authentication: 'all',
            use_mapcache: false,
            selection_color: '#ffff00ff',
            topics: []
          }
        })
    },
    async uploadFiles () {
      const { newFiles, modifiedFiles } = this.$refs.filesBrowser
      const files = [...newFiles, ...modifiedFiles]
      if (files.length === 0) {
        this.saveConfig()
        return
      }

      this.upload = createUpload(this.$ws, files, this.projectServerDir)
      this.uploadProgress = this.upload.info
      try {
        await this.upload.start()
        this.saveConfig()
      } catch (e) {
        if (e !== 'aborted') {
          console.error(e)
          this.$notification.error(e)
        }
        this.fetchServerFiles()
      } finally {
        this.uploadProgress = null
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
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  .v-btn {
    .v-progress-linear {
      position: absolute;
      width: 100%;
      bottom: -3px;
    }
  }
}
</style>
