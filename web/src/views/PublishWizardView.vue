<template>
  <div class="page">
    <portal to="menu-breadcrumbs" v-if="pageVisible">
      <v-icon>keyboard_arrow_right</v-icon>
      <v-btn text rounded color="orange">
        <v-icon class="mr-1">$vuetify.icons.qgis</v-icon>
        Publish
      </v-btn>
    </portal>

    <!-- Timeline -->
    <v-layout class="column left-panel mr-1">
      <v-timeline dense class="mr-2 py-0">
        <v-timeline-item
          class="flat grey--text pb-2"
          color="transparent"
          icon-color="grey"
          large
        >
          <template v-slot:icon>
            <icon name="desktop" class="desktop-img"/>
          </template>
        </v-timeline-item>

        <!-- <v-layout align-center>
          <v-divider class="ml-2"/>
          <icon name="desktop" class="mx-2 desktop-img" style="width:36px"/>
          <v-divider class="mr"/>
        </v-layout> -->

        <v-expand-transition>
          <div v-if="!$ws.pluginConnected">
            <v-timeline-item
              icon-color="deep-orange"
              color="transparent"
              icon="close"
              class="pb-1"
              small
            />
          </div>
        </v-expand-transition>

        <!-- Step 1 -->
        <timeline-primary-link
          label="1. Check project"
          :to="{name: 'qgis-project'}"
          color="lime darken-2"
          icon="$vuetify.icons.qgis"
          :disabled="!$ws.pluginConnected"
          :links="[
            {text: 'General', page: 'info'},
            {text: 'Layers', page: 'layers'}
          ]"
        />

        <!-- Step 2 -->
        <timeline-primary-link
          label="2. Upload files"
          :to="{name: 'publish-upload'}"
          color="lime darken-2"
          icon="cloud_upload"
          :disabled="!$ws.pluginConnected"
        />

        <!-- <v-divider class="ml-3 my-4"/> -->
        <!-- <v-timeline-item
          class="flat grey--text py-2 align-center"
          color="transparent"
          icon-color="grey"
        >
          <template v-slot:icon>
            <icon name="server" class="desktop-img"/>
          </template>
          <v-divider/>
        </v-timeline-item> -->

        <v-layout align-center>
          <v-divider class="ml-2"/>
          <icon name="server" class="mx-2 desktop-img" style="width:36px"/>
          <v-divider/>
        </v-layout>

        <!-- Step 3 -->
        <timeline-primary-link
          label="3. Configure"
          :to="{name: 'publish-config'}"
          color="primary"
          icon="settings"
          :disabled="!serverActionsEnabled"
          :links="[
            {text: 'Project', page: 'project'},
            {text: 'Layers', page: 'layers'},
            {text: 'Topics', page: 'topics'}
          ]"
        />
        <!-- Step 4 -->
        <timeline-primary-link
          label="4. Publish"
          :xto="{name: 'publish-publish'}"
          color="primary"
          icon="public"
          :disabled="!serverActionsEnabled"
          @click="publish"
        />

        <v-timeline-item
          class="flat grey--text pb-0 pt-2 align-center"
          color="transparent"
          icon-color="grey"
          large
        >
          <template v-slot:icon>
            <icon name="browser" class="desktop-img"/>
          </template>
        </v-timeline-item>
      </v-timeline>
      <v-spacer/>
    </v-layout>

    <div class="content main">
      <keep-alive>
        <plugin-disconnected
          v-if="$route.meta.requiresPlugin && !$ws.pluginConnected"
          class="grow"
        />
        <router-view v-else class="scroll-area"/>
        <!-- <router-view v-if="$ws.pluginConnected" class="scroll-area"/>
        <plugin-disconnected v-else class="grow"/> -->
      </keep-alive>
    </div>
  </div>
</template>

<script>
import _omit from 'lodash/omit'
import { basename, extname } from 'path'
import { layersList, filterLayers, scalesToResolutions } from '@/utils'
import Page from '@/mixins/Page'
import PluginDisconnected from '@/components/PluginDisconnected'
import TimelinePrimaryLink from '@/components/TimelinePrimaryLink'

export default {
  name: 'PublishWizardView',
  mixins: [ Page ],
  components: { PluginDisconnected, TimelinePrimaryLink },
  data () {
    return {
      projectInfo: null,
      projectConfig: null,
      serverFiles: null
    }
  },
  computed: {
    user () {
      return this.$root.user
    },
    projectDirectory () {
      return this.projectInfo && this.projectInfo.directory
    },
    projectPath () {
      if (this.projectDirectory) {
        return `${this.user.username}/${basename(this.projectDirectory)}`
      }
      return ''
    },
    serverActionsEnabled () {
      return this.projectConfig !== null
      // if (this.projectInfo && this.serverFiles) {
      //   const projectFile = basename(this.projectInfo.file)
      //   return this.serverFiles.some(f => f.path === projectFile)
      // }
      // return false
    },
    overlays () {
      return this.projectConfig && filterLayers(this.projectConfig.layers, l => l.publish)
    }
  },
  created () {
    // this.$ws.bind('PluginConnected', this.fetchProjectInfo)
    // this.fetchProjectInfo()
  },
  activated () {
    this.$ws.bind('ProjectChanged', this.onProjectChange)
    if (!this.projectInfo) {
      this.fetchProjectInfo()
    }
  },
  deactivated () {
    this.$ws.unbind('ProjectChanged', this.onProjectChange)
  },
  beforeDestroy () {
    // this.$ws.unbind('PluginConnected', this.fetchProjectInfo)
  },
  watch: {
    '$ws.pluginConnected': {
      immediate: true,
      handler (connected) {
        if (this.pageVisible && connected && !this.projectInfo) {
          this.fetchProjectInfo()
        }
      }
    }
  },
  methods: {
    onProjectChange () {
      this.projectInfo = null
      this.projectConfig = null
      this.serverFiles = null
      this.$router.push({ name: 'publish' })
      this.fetchProjectInfo()
    },
    fetchProjectInfo () {
      this.$ws.request('ProjectInfo')
        .then(resp => {
          const config = JSON.parse(resp)
          layersList(config.layers).forEach(l => {
            l.publish = true
            l.hidden = false
          })
          Object.assign(config, {
            base_layers: [],
            authentication: 'all',
            use_mapcache: false,
            selection_color: "#ffff00ff"
          })
          this.projectInfo = config
        })
    },
    publish () {
      const meta = _omit(this.projectConfig, ['file', 'directory'])
      const overlays = JSON.parse(JSON.stringify(this.overlays))
      layersList(overlays).forEach(l => {
        l.queryable = l.wfs
        delete l.wfs
        delete l.source
        delete l.publish
      })
      meta.overlays = overlays
      delete meta.layers
      
      Object.assign(meta, {
        publish_date: new Date().toUTCString(),
        publish_date_unix: new Date().getTime()/1000|0,
        gislab_user: this.user.username
      })
      meta.tile_resolutions = scalesToResolutions(meta.scales, meta.units)
      console.log(meta)
      const projFile = this.projectInfo.file
      const metafile = basename(projFile).replace(extname(projFile), ".meta")
      this.$http.post(`/api/project/meta/${this.projectPath}/${metafile}`, meta)
        .then(() => {
          this.$notification.show('Published!')
        })
    }
  }
}
</script>

<style lang="scss" scoped>
.page {
  overflow: hidden;
  display: grid;
  grid-template-columns: minmax(auto, 1fr) auto 1fr;
  grid-template-rows: auto 1fr;
  .timeline {
    grid-row: 1 / 2;
    grid-column: 2 / 3;
  }
  .left-panel {
    grid-row: 1 / 3;
    grid-column: 1 / 2;
    overflow: auto;
  }
  .main {
    grid-row: 2 / 3;
    grid-column: 2 / 3; 
  }
  .scroll-area {
    overflow: auto;
  }
}
.v-timeline {
  position: relative;
  max-width: 300px;
  &:before {
    top: 16px;
    bottom: 16px;
    height: auto;
  }
  .flat {
    // height: 3em;
    ::v-deep .v-timeline-item__dot {
      background-color: #fafafa;
      box-shadow: none;
    }
  }
  .icon.desktop-img {
    color: #bbb;
    width: 100%;
    height: 32px;
  }
  .v-divider {
    border-style: dashed;
  }
}
.timeline {
  position: relative;
  display: grid;
  grid-template-columns: auto 1fr auto 1fr auto 1fr auto;
  grid-template-rows: auto auto;
  align-items: center;
  
  .submenu {
    grid-row: 2 / 3;
    font-size: 14px;
  }
  
  .v-btn {
    &:first-child {
      transform-origin: 0 50%;
    }
    text-transform: none;
    .icon {
      color: inherit;
      transition: transform 0.3s linear;
      &.flip {
        transform: rotateY(180deg);
      }
    }
    &.step {
      // pointer-events: none;
      min-width: 180px;
      height: 36px;
      // transition: all 0.3s linear;
      // transform: scale(0.75, 0.75);
      &.v-btn--active {
        // height: 42px;
        // transform: scale(1, 1);
        transform: scale(1.15, 1.15);
      }
      &.theme--light {
        background-color: #ddd;
      }
    }
  }
  .v-divider {
    border-style: dashed;
  }
}
</style>
