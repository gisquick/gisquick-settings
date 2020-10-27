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
          label: 'Update',
          disabled: !configChanged
        }"
        @publish="saveChanges"
      />
      <v-spacer/>
    </v-layout>

    <div class="content my-2 mx-1">
      <expander/>
      <keep-alive>
        <router-view v-if="projectConfig"/>
      </keep-alive>
    </div>
  </div>
</template>

<script>
import isEqual from 'lodash/isEqual'
import Page from '@/mixins/Page'
import ProjectMenu from '@/components/ProjectMenu'
import Timeline from '@/components/Timeline'
import ScriptsStore from '@/components/ScriptsStore'
import { layersList, filterLayers } from '@/utils'

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
      layers: [],
      projectConfig: null,
      projectConfigOriginal: null,
      projectInfo: null
    }
  },
  computed: {
    projectPath () {
      return [this.user, this.folder, this.projectName].join('/')
    },
    overlays () {
      return this.projectConfig && filterLayers(this.projectConfig.overlays, l => l.publish)
    },
    configChanged () {
      return !isEqual(this.projectConfig, this.projectConfigOriginal)
    }
  },
  watch: {
    projectPath: {
      immediate: true,
      handler (projectName) {
        if (projectName) {
          this.fetchProjectConfig()
        } else {
          this.projectConfig = null
        }
      }
    }
  },
  methods: {
    fetchProjectConfig () {
      this.$http.get(`/api/project/meta/${this.projectPath}`)
        .then(resp => {
          // TODO: load projectInfo file when available
          const layers = [...resp.data.base_layers, ...resp.data.overlays]
          // resp.data.layers = resp.data.overlays
          layersList(layers)
            .filter(l => l.publish === undefined)
            .forEach(l => {
              l.publish = true
            })
          this.projectInfo = {
            layers
          }
          this.projectConfig = resp.data
          this.projectConfigOriginal = JSON.parse(JSON.stringify(this.projectConfig))
        })
        .catch(err => {
          console.error(err)
          this.projectConfig = null
          this.projectInfo = null
        })
    },
    saveChanges () {
      const project = this.projectConfig.project || this.this.projectPath
      this.$http.post(`/api/project/meta/${project}`, this.projectConfig)
        .then(() => {
          this.$notification.show('Updated!')
          this.projectConfigOriginal = JSON.parse(JSON.stringify(this.projectConfig))
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
