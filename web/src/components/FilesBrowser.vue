<template>
  <v-layout class="browser">
    <!-- Local -->
    <v-layout class="column box">
      <v-toolbar
        dark
        elevation="2"
        class="shrink"
        color="grey darken-2"
        height="42"
      >
        <v-toolbar-title>Local files</v-toolbar-title>
        <v-spacer/>
        <v-toolbar-items>

        </v-toolbar-items>
      </v-toolbar>
      <template v-if="$ws.pluginConnected">
        <v-text-field
          label="Directory"
          :value="srcPath"
          class="shrink mt-2 px-4"
          readonly
          disabled
          hide-details
        />
        <v-layout
          v-if="srcLoading"
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
                  :class="{'progress-active': filesProgress && filesProgress[item.path]}"
                >
                  <span>{{ item.size | filesize }}</span>
                  <v-progress-linear
                    v-if="filesProgress && filesProgress[item.path]"
                    :value="filesProgress[item.path].progress"
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
    <v-layout class="column box ml-1">
      <v-toolbar
        dark
        elevation="2"
        class="shrink"
        color="grey darken-2"
        height="42"
      >
        <v-toolbar-title>Server files</v-toolbar-title>
        <v-spacer/>
        <v-toolbar-items>
          <slot name="dest-toolbar"/>
        </v-toolbar-items>
      </v-toolbar>
      <v-text-field
        label="Project"
        :value="destPath"
        class="shrink mt-2 px-4"
        readonly
        disabled
        hide-details
      />
      <v-layout
        v-if="destLoading"
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
          v-if="destLoading || destFiles.length > 0"
          :items="serverFilesTree"
          class="mt-2 px-2"
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
        <v-layout v-else fill-height align-center justify-center>
          <span class="grey--text">Empty</span>
        </v-layout>
      </div>
    </v-layout>
  </v-layout>
</template>

<script>
import Path from 'path'
import { dirname, basename } from 'path'
import _keyBy from 'lodash/keyBy'

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
  name: 'FilesBrowser',
  props: {
    srcPath: String,
    destPath: String,
    srcLoading: Boolean,
    destLoading: Boolean,
    srcFiles: {
      type: Array,
      default: () => []
    },
    destFiles: {
      type: Array,
      default: () => []
    },
    filesProgress: Object
  },
  computed: {
    localFiles () {
      return [...this.srcFiles].sort(compareFilenames)
    },
    serverFiles () {
      return [...this.destFiles].sort(compareFilenames)
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
      const files = this.localFiles.filter(lf => this.serverFiles.find(sf => sf.path === lf.path && sf.hash !== lf.hash))
      return _keyBy(files, 'path')
    },
    filesStyles () {
      const styles = {}
      this.localFiles.forEach(f => {
        if (this.newFiles[f.path]) {
          var color = '#689F38'
        } else if (this.modifiedFiles[f.path]) {
          color = '#FFA000'
        } else {
          color = '#444'
        }
        styles[f.path] = { color }
      })
      return styles
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
.browser {
  // display: grid;
  // grid-gap: 4px;
  // grid-template-columns: 1fr 1fr;
  // grid-template-rows: 1fr auto;

  .box {
    border: 1px solid #aaa;
    overflow: hidden;
    flex: 1 1;
  }
  .v-treeview {
    .v-icon {
      color: inherit;
      opacity: 0.75;
    }
  }
}
</style>
