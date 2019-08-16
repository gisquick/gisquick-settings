<template>
  <keep-alive>
    <router-view v-if="config" :config="config"/>
  </keep-alive>
</template>

<script>
export default {
  name: 'Project',
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
