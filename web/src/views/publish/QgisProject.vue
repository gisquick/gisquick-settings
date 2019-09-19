<template>
  <keep-alive>
    <div
      v-if="route"
      :is="route.component"
      v-bind="route.props"
      class="box grow"
    />
  </keep-alive>
</template>

<script>
import DefaultRedirect from '@/mixins/DefaultRedirect'
import GeneralInfo from './GeneralInfo'
import LayersInfo from './LayersInfo'

export default {
  name: "QgisProject",
  mixins: [ DefaultRedirect({ page: 'info' }) ],
  computed: {
    store () {
      return this.$parent.$data
    },
    config () {
      return this.store.projectInfo
    },
    routes () {
      return this.config && {
        layers: {
          component: LayersInfo,
          props: {
            config: this.config
          }
        },
        info: {
          component: GeneralInfo,
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
.box {
  border: 1px solid #ddd;
}
</style>
