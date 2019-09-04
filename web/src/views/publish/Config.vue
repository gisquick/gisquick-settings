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
import LayersView from '@/views/LayersView'
import TopicsView from '@/views/TopicsView'

export default {
  name: 'Config',
  mixins: [ DefaultRedirect({ page: 'layers' }) ],
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
        layers: {
          component: LayersView,
          props: {
            config: this.config
          }
        },
        topics: {
          component: TopicsView,
          props: {
            config: this.config
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
