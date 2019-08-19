<template>
  <div class="page py-2">
    <v-layout class="column box">
      <v-toolbar
        dark dense
        elevation="2"
        class="shrink"
        color="lime darken-2"
      >
        <v-toolbar-title>Local files</v-toolbar-title>
        <v-spacer/>
        <v-toolbar-items>
          <v-btn icon @click="fetchLocalFiles">
            <v-icon>refresh</v-icon>
          </v-btn>
        </v-toolbar-items>
      </v-toolbar>
      <template v-if="$ws.pluginConnected">
        <v-text-field
          label="Directory"
          :value="localDirectory"
          class="shrink mt-2 px-4"
          readonly
          disabled
          hide-details
        />
        <v-layout
          v-if="loadingLocalFiles"
          class="column grow align-center justify-center mx-4 my-4"
        >
          <div class="subtitle-1">Loading</div>
          <v-progress-linear
            indeterminate
            rounded
            height="6"
          />
        </v-layout>
        <div v-else class="scroll-container mt-1">
          <v-treeview
            :items="localFilesTree"
            class="mt-2 px-2"
            dense
          >
            <template v-slot:label="{ item, open }">
              <v-layout :style="filesStyles[item.path]">
                <v-icon v-if="item.children">
                  {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
                </v-icon>
                <v-icon v-else>mdi-file-document-outline</v-icon>
                <span>{{ item.name }}</span>
                <v-spacer/>
                <v-layout
                  class="column filesize upload shrink ml-2"
                  :class="{'progress-active': upload && upload[item.path]}"
                >
                  <span>{{ item.size | filesize }}</span>
                  <v-progress-linear
                    v-if="upload && upload[item.path]"
                    :value="upload[item.path].progress"
                  />
                </v-layout>
              </v-layout>
            </template>
          </v-treeview>
        </div>
      </template>
      <p v-else class="mx-3 my-3">QGIS plugin not connected</p>
    </v-layout>

    <!-- Server -->
    <v-layout class="column box">
      <v-toolbar
        dark dense
        elevation="2"
        class="shrink"
        color="lime darken-2"
      >
        <v-toolbar-title>Server files</v-toolbar-title>
        <v-spacer/>
        <v-toolbar-items>
          <v-btn
            icon
            :disabled="loadingServerFiles"
            @click="fetchServerFiles"
          >
            <v-icon>refresh</v-icon>
          </v-btn>
        </v-toolbar-items>
      </v-toolbar>
      <v-text-field
        label="Project"
        :value="projectPath"
        class="shrink mt-2 px-4"
        readonly
        disabled
        hide-details
      />
      <v-layout
        v-if="loadingServerFiles"
        class="column grow align-center justify-center mx-4 my-4"
      >
        <div class="subtitle-1">Loading</div>
        <v-progress-linear
          indeterminate
          rounded
          height="6"
        />
      </v-layout>
      <div v-else class="scroll-container mt-1">
        <v-treeview
          :items="serverFilesTree"
          class="mt-4 px-2"
          dense
        >
          <template v-slot:prepend="{ item, open }">
            <v-icon v-if="item.children">
              {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
            </v-icon>
            <v-icon v-else>mdi-file-document-outline</v-icon>
          </template>
          <template v-slot:append="{ item, open }">
            <span class="ml-4">{{ item.size | filesize }}</span>
          </template>
        </v-treeview>
      </div>
    </v-layout>
    <div class="toolbar">
      <v-spacer/>
      <v-btn rounded @click="startUpload">
        <v-icon class="mr-2">cloud_upload</v-icon>
        <span>Upload</span>
      </v-btn>
      <v-btn rounded :href="`/api/project/download/${projectPath}`" class="self-align-end">
        <v-icon class="mr-2">archive</v-icon>
        <span>Download</span>
      </v-btn>
    </div>
  </div>
</template>

<script>
import Path from 'path'
import { dirname, basename } from 'path'
import _keyBy from 'lodash/keyBy'
import TreeView from '@/components/TreeView.vue'

function compareFilenames (a, b) {
  const depthA = a.path.split(Path.sep).length
  const depthB = b.path.split(Path.sep).length
  if (depthA === depthB) {
    return basename(a.path).localeCompare(basename(b.path))
  }
  return depthA - depthB
}

function filesTree (files) {
  const root = { children: [] }
  files.forEach(f => {
    let parent = root
    const dir = dirname(f.path)
    if (dir !== '.') {
      const parts = dir.split(Path.sep)
      parts.forEach(name => {
        let node = parent.children.find(i => i.name === name)
        if (!node) {
          node = {name: name, children: []}
          parent.children.push(node)
        }
        parent = node
      })
    }
    parent.children.push({
      ...f,
      name: basename(f.path)
    })
  })
  return root.children
}

export default {
  name: 'Files',
  components: { TreeView },
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
    // user () {
    //   return this.$root.user
    // },
    projectPath () {
      return this.user && this.folder && `${this.user}/${this.folder}`
    },
    localFilesTree () {
      return filesTree(this.localFiles)
    },
    serverFilesTree () {
      return filesTree(this.serverFiles)
    },
    newFiles () {
      const files = this.localFiles.filter(lf => !this.serverFiles.find(sf => sf.path === lf.path))
      return _keyBy(files, 'path')
    },
    modifiedFiles () {
      const files = this.localFiles.filter(lf => this.serverFiles.find(sf => sf.path === lf.path && sf.checksum !== lf.checksum))
      return _keyBy(files, 'path')
    },
    filesStyles () {
      const styles = {}
      this.localFiles.forEach(f => {
        if (this.newFiles[f.path]) {
          var color = 'green'
        } else if (this.modifiedFiles[f.path]) {
          color = 'orange'
        } else {
          color = '#444'
        }
        styles[f.path] = { color }
      })
      return styles
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
      this.localFiles = data.files.sort(compareFilenames)
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
      this.loadingLocalFiles = true
      this.$ws.send('Files')
    },
    fetchServerFiles () {
      this.loadingServerFiles = true
      this.$http.get(`/api/project/files/${this.projectPath}`)
        .then(resp => {
          this.serverFiles = resp.data.sort(compareFilenames)
          this.loadingServerFiles = false
        })
        .catch(err => {
          this.loadingServerFiles = false
        })
    },
    startUpload () {
      // const files = this.localFiles.filter(f => this.newFiles[f.path] || this.modifiedFiles[f.path])
      const files = this.localFiles.filter(f => this.newFiles[f.path])
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
.scroll-container {
  overflow: auto;
  flex: 1 1;
}
.upload-bar {
  width: 70px;
  margin: auto 0;
}
.filesize {
  width: 80px;
  text-align: right;

  &.upload {
    position: relative;
    &.progress-active {
      span {
        transform: scale(0.85, 0.85);
      }
    }
    span {
      transform-origin: top right;
      transition: 0.25s transform ease;
    }
    .v-progress-linear {
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
    }
  }
}
.page {
  display: grid;
  grid-gap: 4px;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr auto;

  .box {
    border: 1px solid #aaa;
    overflow: hidden;
  }
  .v-treeview {
    min-width: 300px;
    .v-icon {
      color: inherit;
      opacity: 0.75;
    }
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
