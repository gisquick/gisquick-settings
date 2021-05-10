<template>
  <keep-alive>
    <div
      v-if="route"
      :is="route.component"
      v-bind="route.props"
      class="fill-height"
    />
    <!-- <div v-else>404</div> -->
  </keep-alive>
</template>

<script>
import DefaultRedirect from '@/mixins/DefaultRedirect'
import ProjectSettings from '@/views/publish/ProjectSettings'
import LayersView from '@/views/publish/LayersSettings'
import LayersPermissions from '@/views/publish/LayersPermissions'
import TopicsEditor from '@/components/TopicsEditor'
import LayerAttributes from '@/components/LayerAttributes'
// import NotFound from '@/components/NotFound'
import { layersList } from '@/utils'

export default {
  name: 'Config',
  mixins: [ DefaultRedirect({ page: 'project' }) ],
  props: {
    page: String
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
      const { overlays, settings } = this.store
      return [
        {
          path: 'project',
          component: ProjectSettings,
          props: {
            layers: this.store.projectConfig.layers,
            config: this.config,
            projectPath: this.projectPath,
            settings
          }
        }, {
          path: 'layers',
          component: LayersView,
          props: {
            order: this.store.projectConfig.layers.map(l => l.name),
            config: this.config,
            settings
          }
        }, {
          path: 'topics',
          component: TopicsEditor,
          props: {
            topics: settings.topics,
            layers: overlays
          }
        }, {
          path: 'layers/permissions',
          component: LayersPermissions,
          props: {
            layers: overlays,
            settings
          }
        }, {
          path: path => {
            const match = /([^/]+)\/attributes/.exec(path)
            return match && { layername: match[1] }
          },
          component: LayerAttributes,
          props: ({ layername }) => {
            const layer = layersList(overlays).find(l => l.name === layername)
            const layerSettings = settings.layers[layer.id]
            return { layer, layerSettings }
          }
        }
      ]
    },
    route () {
      const path = this.$route.params.page
      for (let route of this.routes) {
        const match = typeof route.path === 'string' ? route.path === path : route.path(path)
        if (match) {
          return {
            component: route.component,
            props: typeof route.props === 'function' ? route.props(match) : route.props
          }
        }
      }
      return null
    }
  }
}
</script>
