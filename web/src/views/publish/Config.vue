<template>
  <keep-alive>
    <div
      v-if="route"
      :is="route.component"
      v-bind="route.props"
      class="grow"
    />
  </keep-alive>
</template>

<script>
import DefaultRedirect from '@/mixins/DefaultRedirect'
// import LayersView from '@/views/LayersView'
import Project from '@/views/publish/Project'
import LayersView from '@/views/publish/LayersSettings'
import TopicsEditor from '@/components/TopicsEditor'

export default {
  name: 'Config',
  mixins: [ DefaultRedirect({ page: 'project' }) ],
  props: {
    page: String
  },
  data () {
    return {
    }
  },
  computed: {
    store () {
      return this.$parent
    },
    config () {
      return this.store.projectConfig
    },
    projectPath () {
      return this.store.projectPath
    },
    routes () {
      return this.config && {
        project: {
          component: Project,
          props: {
            config: this.config
          }
        },
        layers: {
          component: LayersView,
          props: {
            config: this.config
          }
        },
        topics: {
          component: TopicsEditor,
          props: {
            topics: this.config.topics,
            layers: this.store.overlays
          }
        }
      }
    },
    route () {
      return this.routes && this.routes[this.$route.params.page]
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
