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
    <v-layout class="column left-panel py-2 mt-2 mr-1">
      <timeline
        :checkin="{
          label: '1. Check-in',
          link: { name: 'qgis-project' },
          disabled: !$ws.pluginConnected,
          sublinks: [
            {label: 'General', page: 'general', status: checkin.general},
            {label: 'Layers', page: 'layers', status: checkin.layers}
          ]
        }"
        :files="{
          label: '2. Upload',
          link: {name: 'publish-upload'},
          disabled: !progress.projectValid || !$ws.pluginConnected
        }"
        :settings="{
          label: '3. Settings',
          link: {name: 'publish-config'},
          disabled: !progress.filesUploaded,
          sublinks: [
            {label: 'Project', page: 'project'},
            {label: 'Layers', page: 'layers'},
            {label: 'Topics', page: 'topics'}
          ]
        }"
        :publish="{
          label: '4. Publish',
          disabled: this.projectConfig === null || !visitedLinks['publish-config'],
          published: this.progress.published
        }"
        :visited="visitedLinks"
        @publish="publish"
      />
      <v-spacer/>
    </v-layout>

    <v-layout class="column content mx-1">
      <expander/>
      <keep-alive>
        <plugin-disconnected
          v-if="$route.meta.requiresPlugin && !$ws.pluginConnected"
          class="grow disconnect-msg my-2"
        />
        <router-view v-else class="scroll-area my-2"/>
      </keep-alive>
    </v-layout>
  </div>
</template>

<script>
import { basename, extname } from 'path'
import { filterLayers, scalesToResolutions } from '@/utils'
import Page from '@/mixins/Page'
import PluginDisconnected from '@/components/PluginDisconnected'
import Timeline from '@/components/Timeline'

export default {
  name: 'PublishWizardView',
  mixins: [ Page ],
  components: { PluginDisconnected, Timeline },
  data () {
    return {
      projectInfo: null,
      projectConfig: null,
      serverFiles: null,
      publishProgress: {
        published: false
      },
      visitedLinks: {},
      checkin: {
        general: null,
        layers: null
      }
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
        const filename = basename(this.projectInfo.file)
        const projectName = filename.substring(0, filename.length - extname(filename).length)
        return `${this.user.username}/${basename(this.projectDirectory)}/${projectName}`
      }
      return ''
    },
    overlays () {
      return this.projectConfig && filterLayers(this.projectConfig.overlays, l => l.publish)
    },
    progress () {
      const { published } = this.publishProgress
      return {
        projectInfo: this.projectInfo !== null,
        projectValid: this.checkin.general && this.checkin.general !== 'error' && this.checkin.layers && this.checkin.layers !== 'error',
        filesUploaded: this.projectConfig !== null,
        published
      }
    }
  },
  activated () {
    const unbind = this.$ws.bind('ProjectChanged', this.onProjectChange)
    this.$once('hook:deactivated', unbind)
  },
  beforeRouteLeave (to, from, next) {
    this.resetProjectData()
    next()
  },
  watch: {
    '$ws.pluginConnected': {
      handler (connected) {
        if (!connected && this.$route.meta.requiresPlugin) {
          this.resetProjectData()
        }
      }
    },
    $route: {
      immediate: true,
      handler (route) {
        if (route.path.startsWith('/publish/') && route.name !== 'qgis-project' && !this.progress.projectValid) {
          // invalid state => redirect to beginning
          this.$router.push({ name: 'publish' })
          return
        }
        this.$set(this.visitedLinks, route.name, true)
      }
    }
  },
  methods: {
    resetProjectData () {
      Object.assign(this.$data, this.$options.data())
    },
    onProjectChange () {
      // TODO: check routes and don't redirect always
      this.resetProjectData()
      this.$router.push({ name: 'publish' })
    },
    publish () {
      const meta = {
        ...this.projectConfig,
        tile_resolutions: scalesToResolutions(this.projectConfig.scales, this.projectConfig.units),
        publish_date: new Date().toUTCString(),
        publish_date_unix: new Date().getTime()/1000|0,
        gislab_user: this.user.username
      }
      this.$http.post(`/api/project/meta/${this.projectPath}`, meta)
        .then(() => {
          this.$notification.show('Published!')
          this.publishProgress.published = true
        })
    }
  }
}
</script>

<style lang="scss" scoped>
.page {
  display: contents;
  .left-panel {
    grid-column: 1 / 2;
    overflow: auto;
    background-color: #353535;
    min-width: 280px;
  }
  .content {
    grid-row: 2 / 3;
    grid-column: 2 / 3;
    width: 1200px;
    overflow: auto;
    @media (max-width: 1600px) {
      width: auto;
    }
    @media (max-width: 1450px) {
      grid-column: 2 / 4;
    }
  }
  .scroll-area {
    overflow: auto;
  }
}
::v-deep .disconnect-msg {
  border: 2px dashed #ddd;
}
</style>
