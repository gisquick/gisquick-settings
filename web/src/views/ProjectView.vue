<template>
  <div class="page">
    <scripts-store/>
    <portal to="menu-breadcrumbs" v-if="pageVisible && projectConfig">
      <v-icon>keyboard_arrow_right</v-icon>
      <v-btn text color="orange" style="text-transform:none">
        {{ projectConfigOriginal.title }}
      </v-btn>
    </portal>
    <portal to="menu-actions" v-if="pageVisible">
      <!-- <project-menu :user="user" :folder="folder"/> -->
      <project-menu v-bind="$route.params"/>
    </portal>

    <v-layout class="column py-2 left-panel mt-2 mr-1 elevation-3">
      <timeline
        :files="{
          label: 'Files',
          link: {name: 'files'}
        }"
        :settings="{
          label: 'Settings',
          link: {name: 'settings'},
          activeClass: 'primary--text',
          sublinks: [
            {label: 'Project', page: 'project'},
            {label: 'Layers', page: 'layers'},
            {label: 'Topics', page: 'topics'}
          ]
        }"
        :publish="{
          label: this.originalSettings === null ? 'Publish' : 'Update',
          disabled: !metaChanged && !settingsChanged
        }"
        @publish="saveChanges"
      />
      <v-spacer/>
    </v-layout>

    <div class="content my-2 mx-1">
      <expander/>
      <keep-alive>
        <router-view v-if="projectConfig && settings"/>
      </keep-alive>
    </div>
  </div>
</template>

<script>
import isEqual from 'lodash/isEqual'
import keyBy from 'lodash/keyBy'
import mapValues from 'lodash/mapValues'
import Page from '@/mixins/Page'
import ProjectMenu from '@/components/ProjectMenu'
import Timeline from '@/components/Timeline'
import ScriptsStore from '@/components/ScriptsStore'
import { layersList, filterLayers, scalesToResolutions } from '@/utils'

function clone (data) {
  return JSON.parse(JSON.stringify(data))
}

export default {
  name: 'Project',
  mixins: [ Page ],
  components: { ProjectMenu, Timeline, ScriptsStore },
  props: {
    user: String,
    folder: String,
    projectName: String
  },
  data () {
    return {
      projectConfig: null,
      projectConfigOriginal: null,
      originalSettings: null,
      settings: null
    }
  },
  computed: {
    projectPath () {
      return [this.user, this.folder, this.projectName].join('/')
    },
    overlays () {
      const { layers, overlays } = this.settings
      const items = this.projectConfig.layers.filter(l => overlays.includes(l.id || l.name))
      return filterLayers(items, l => layers[l.id]?.publish)
    },
    metaChanged () {
      return !isEqual(this.projectConfig, this.projectConfigOriginal)
    },
    settingsChanged () {
      return !isEqual(this.settings, this.originalSettings)
    }
  },
  watch: {
    projectPath: {
      immediate: true,
      handler: 'loadProject'
    }
  },
  methods: {
    async loadProject (projectName) {
      if (projectName) {
        try {
          const meta = await this.fetchProjectMeta()
          const settings = await this.fetchGisquickSettings()
          this.projectConfig = meta
          this.projectConfigOriginal = clone(meta)

          // meta.layers.forEach(l => Object.defineProperty(l, 'settings', { configurable: true, value: data.layers[l.id] }))
          if (settings) {
            this.settings = settings
            this.originalSettings = clone(settings)
          } else {
            this.settings = this.newSettings(meta)
            this.originalSettings = null
          }
          const newLayers = layersList(this.projectConfig.layers).filter(l => !this.settings.layers[l.id])
          newLayers.forEach(l => this.$set(settings.layers, l.id, {
            publish: true,
            hidden: false,
            attributes: {}
          }))
        } catch (err) {
          console.error(err)
          this.$notification.error('Failed to load project data')
        }
      } else {
        this.projectConfig = null
        this.projectConfigOriginal = null
        this.settings = null
        this.originalSettings = null
      }
    },
    async fetchProjectMeta () {
      const { data } = await this.$http.get(`/api/project/meta/${this.projectPath}`)
      return data
    },
    newSettings (meta) {
      const layers = mapValues(keyBy(layersList(meta.layers), 'id'), () => ({
        publish: true,
        hidden: false,
        attributes: {}
      }))
      return {
        layers,
        extent: meta.extent,
        scales: meta.scales,
        tile_resolutions: scalesToResolutions(meta.scales, meta.units),
        title: meta.title,
        overlays: meta.layers.map(l => l.id || l.name),
        base_layers: [],
        authentication: 'all',
        use_mapcache: false,
        topics: []
      }
    },
    async fetchGisquickSettings () {
      const { data } = await this.$http
        .get(`/api/project/settings/${this.projectPath}`)
        .catch(err => {
          if (err.response.status === 404) {
            return { data: null }
          }
          throw err
        })
      return data
    },
    async saveMeta () {
      const project = this.projectConfig.project || this.projectPath
      await this.$http.post(`/api/project/meta/${project}`, this.projectConfig)
      this.projectConfigOriginal = JSON.parse(JSON.stringify(this.projectConfig))
    },
    async saveSettings () {
      const project = this.projectConfig.project || this.projectPath
      await this.$http.post(`/api/project/settings/${project}`, this.settings)
      this.originalSettings = JSON.parse(JSON.stringify(this.settings))
    },
    async saveChanges () {
      const tasks = []
      if (this.metaChanged) {
        tasks.push(this.saveMeta())
      }
      if (this.settingsChanged) {
        tasks.push(this.saveSettings())
      }
      try {
        await Promise.all(tasks)
        this.$notification.show('Updated!')
      } catch (err) {
        console.error(err)
        this.$notification.error('Failed to save changes!')
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.page {
  display: contents;
  .left-panel {
    grid-column: 1 / 2;
  }
}
.left-panel {
  background-color: #373737;
  min-width: 280px;
  overflow: auto;
}
.content {
  grid-column: 2 / 3;
  width: 1200px;
  overflow: auto;
  @media (max-width: 1600px) {
    width: auto;
  }
  @media (max-width: 1450px) {
    width: auto;
    grid-column: 2 / 4;
  }
}
</style>
