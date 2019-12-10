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
          <slot name="src-toolbar"/>
        </v-toolbar-items>
      </v-toolbar>
      <template v-if="$ws.pluginConnected">
        <v-text-field
          v-if="srcPath"
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
        <div v-else class="scroll-container grow mt-1">
          <v-treeview
            v-if="srcPath"
            :items="localFilesTree"
            item-key="path"
            class="mt-2 px-2"
          >
            <template v-slot:label="{ item, open }">
              <v-layout :style="item.children ? dirsStyles[item.path] : filesStyles[item.path]">
                <v-icon v-if="item.children">
                  {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
                </v-icon>
                <file-icon v-else :path="item.path"/>
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
          <slot v-else name="src-no-content"/>
        </div>
      </template>
      <plugin-disconnected v-else class="disconnect-msg"/>
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
      <div v-else class="scroll-container grow mt-1">
        <v-treeview
          v-if="destLoading || (destFiles && destFiles.length > 0)"
          :items="serverFilesTree"
          item-key="path"
          class="mt-2 px-2"
          dense
        >
          <template v-slot:prepend="{ item, open }">
            <v-icon v-if="item.children">
              {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
            </v-icon>
            <file-icon v-else :path="item.path"/>
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
import { dirname, basename, extname } from 'path'
// import { win32 as dirname } from 'path-dirname'

import _keyBy from 'lodash/keyBy'
import _mapValues from 'lodash/mapValues'
import PluginDisconnected from '@/components/PluginDisconnected'


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
          node = { name: name, children: [], path: dir }
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

const FileIcon = {
  functional: true,
  props: {
    path: String
  },
  render (h, ctx) {
    const ext = extname(ctx.props.path).toLowerCase()
    let icon = 'mdi-file-document-outline'
    if (ext === '.qgs' || ext === '.qgz') {
      icon = '$vuetify.icons.qgis'
    } else if (ext === '.sqlite') {
      icon = '$vuetify.icons.db'
    } else if (ext.match(/\.(jpeg|jpg|gif|png|svg|tif|tiff)$/) !== null) {
      icon = 'image'
    }
    return <v-icon>{icon}</v-icon>
  }
}
export default {
  name: 'FilesBrowser',
  components: { PluginDisconnected, FileIcon },
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
    loading () {
      return this.srcLoading || this.destLoading
    },
    newFiles () {
      if (this.loading) {
        return []
      }
      return this.localFiles.filter(lf => !this.serverFiles.find(sf => sf.path === lf.path))
    },
    newFilesMap () {
      return _keyBy(this.newFiles, 'path')
    },
    modifiedFiles () {
      if (this.loading) {
        return []
      }
      return this.localFiles.filter(lf => this.serverFiles.find(sf => sf.path === lf.path && sf.hash !== lf.hash))
    },
    modifiedFilesMap () {
      return _keyBy(this.modifiedFiles, 'path')
    },
    filesStyles () {
      const styles = {}
      this.localFiles.forEach(f => {
        if (this.newFilesMap[f.path]) {
          var color = '#689F38'
        } else if (this.modifiedFilesMap[f.path]) {
          color = '#FFA000'
        } else {
          color = '#222'
        }
        styles[f.path] = { color }
      })
      return styles
    },
    dirsInfo () {
      const dirsInfo = {}
      const groupInfo = item => {
        if (item.children) {
          // compute nested groupInfo first
          item.children.filter(i => i.children).forEach(groupInfo)

          const isNew = item.children.every(i => i.children ? dirsInfo[i.path].isNew : this.newFilesMap[i.path])
          const isMod = !isNew && item.children.some(
            i => i.children
            ? dirsInfo[i.path].isNew || dirsInfo[i.path].isMod
            : this.modifiedFilesMap[i.path] || this.newFilesMap[i.path]
          )
          dirsInfo[item.path] = { isNew, isMod }
        }
      }
      this.localFilesTree.filter(item => item.children).forEach(groupInfo)
      return dirsInfo
    },
    dirsStyles () {
      return _mapValues(this.dirsInfo, gi => ({
        color: gi.isNew ? '#689F38' : gi.isMod ? '#FFA000' : '#222'
      }))
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
      padding-right: 2px;
      margin-right: 4px;
    }
    ::v-deep {
     .v-treeview-node__label {
        margin-left: 0;
      }
      .v-treeview-node__root {
        min-height: 38px;
      }
    }
  }
}
 ::v-deep .disconnect-msg {
   padding: 50px;
  img {
    // width: 400px;
  }
}
</style>
