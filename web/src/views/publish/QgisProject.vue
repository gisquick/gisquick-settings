<template>
    <div
      v-if="route"
      :is="route.component"
      v-bind="route.props"
      class="box grow"
      @status="updateStatus"
    />
</template>

<script>
import DefaultRedirect from '@/mixins/DefaultRedirect'
import GeneralInfo from './GeneralInfo'
import LayersInfo from './LayersInfo'

export default {
  name: "QgisProject",
  mixins: [ DefaultRedirect({ page: 'general' }) ],
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
          },
          validate: 'layersInfo'
        },
        general: {
          component: GeneralInfo,
          props: {
            config: this.config
          },
          validate: 'generalInfo'
        }
      }
    },
    route () {
      return this.routes && this.routes[this.$route.params.page]
    }
  },
  watch: {
    route: {
      immediate: true,
      handler (route) {
        if (route) {
          this.store.publishProgress[route.validate] = true
        }
      }
    }
  },
  methods: {
    updateStatus (status) {
      // this.$emit('status', { route: this.$route.params.page, status })
      this.store.checkin[this.$route.params.page] = status
      this.$emit('status', status)
    }
  }
}
</script>

<style lang="scss" scoped>
.box {
  border: 1px solid #ddd;
}
</style>
