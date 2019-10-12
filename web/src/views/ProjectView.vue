<template>
  <div class="page">
    <portal to="menu-breadcrumbs" v-if="pageVisible && projectConfig">
      <v-icon>keyboard_arrow_right</v-icon>
      <v-btn text color="orange" style="text-transform:none">
        {{ projectConfig.title }}
      </v-btn>
    </portal>
    <portal to="menu-actions" v-if="pageVisible">
      <!-- <project-menu :user="user" :folder="folder"/> -->
      <project-menu v-bind="$route.params"/>
    </portal>

    <v-layout class="column mr-2 py-2 left-panel mt-2 elevation-3">
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
          disabled: true
        }"
      />
      <v-spacer/>
    </v-layout>

    <div class="content my-2">
      <keep-alive>
        <router-view v-if="projectConfig"/>
      </keep-alive>
    </div>
  </div>
</template>

<script>
import Page from '@/mixins/Page'
import ProjectMenu from '@/components/ProjectMenu'
import Timeline from '@/components/Timeline'
import { layersList, filterLayers } from '@/utils'

export default {
  name: 'Project',
  mixins: [ Page ],
  components: { ProjectMenu, Timeline },
  props: {
    user: String,
    folder: String,
    projectName: String
  },
  data () {
    return {
      layers: [],
      projectConfig: null,
      projectInfo: null
    }
  },
  computed: {
    projectPath () {
      return [this.user, this.folder, this.projectName].join('/')
    },
    overlays () {
      return this.projectConfig && filterLayers(this.projectConfig.overlays, l => l.publish)
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
      // const params = { PROJECT: this.projectPath }
      // this.$http.get('/project.json', { params })
      this.$http.get(`/api/project/meta/${this.projectPath}`)
        .then(resp => {
          // TODO: load projectInfo file when available
          const layers = [...resp.data.base_layers, ...resp.data.overlays]
          this.projectInfo = {
            layers
          }
          // resp.data.layers = resp.data.overlays
          // layersList(resp.data.layers).forEach(l => {
          //   l.publish = true
          // })
          this.projectConfig = resp.data
        })
        .catch(err => {
          console.error(err)
          this.projectConfig = null
          this.projectInfo = null
        })
    },
    saveChanges () {

    }
  }
}
</script>

<style lang="scss" scoped>
.page {
  overflow: hidden;
  display: grid;
  grid-template-columns: minmax(auto, 1fr) auto 1fr;
}
.content > * {
  width: 100%;
}
.left-panel {
  background-color: #3f3f3f;
  min-width: 250px;
  overflow: auto;
}
</style>
