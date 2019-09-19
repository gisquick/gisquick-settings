<template>
  <div class="content">
    <portal to="menu-breadcrumbs" v-if="pageVisible">
      <v-icon>keyboard_arrow_right</v-icon>
      <v-btn text color="orange">
        <v-icon class="mr-1">map</v-icon>
        {{ config && config.root_title }}
      </v-btn>
    </portal>
    <portal to="menu-actions" v-if="pageVisible">
      <!-- <project-menu :user="user" :folder="folder"/> -->
      <project-menu v-bind="$route.params"/>
    </portal>
    <keep-alive>
      <router-view v-if="config" :config="config"/>
    </keep-alive>
  </div>
</template>

<script>
import Page from '@/mixins/Page'
import ProjectMenu from '@/components/ProjectMenu'

export default {
  name: 'Project',
  mixins: [ Page ],
  components: { ProjectMenu },
  props: {
    user: String,
    folder: String,
    projectName: String
  },
  data () {
    return {
      config: null
    }
  },
  computed: {
    projectPath () {
      return [this.user, this.folder, this.projectName].join('/')
    },
    // overlayLayers () {
    //   return this.layers.filter(item => !this.baseLayersNames.includes(item.name))
    // },
    // baseLayers () {
    //   return this.layers.filter(item => this.baseLayersNames.includes(item.name))
    // }
  },
  watch: {
    projectName: {
      immediate: true,
      handler (projectName) {
        if (projectName) {
          this.fetchProjectConfig()
        } else {
          this.config = null
        }
      }
    }
  },
  methods: {
    fetchProjectConfig () {
      const params = { PROJECT: this.projectPath }
      this.$http.get('/project.json', { params })
        .then(resp => {
          this.config = resp.data
        })
    }
  }
}
</script>

<style lang="scss" scoped>
.content > * {
  width: 100%;
}
</style>
