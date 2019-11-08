<template>
  <div
    v-if="route"
    :is="route.component"
    v-bind="route.props"
    class="box grow"
    @status="updateStatus"
  />
  <v-layout
    v-else-if="error"
    class="column box grow px-2 py-2 align-center justify-center"
  >
    <h2 class="mb-4">
      <v-icon size="30" color="orange" class="mr-1">warning</v-icon>
      <span>QGIS project data is not available</span>
    </h2>
    <p v-if="error.code === 404">Please open the project you would like to publish in QGIS desktop application</p>
    <p v-else class="red--text">
      <v-icon color="red" class="mr-1">error</v-icon>
      {{ error.msg }}
    </p>
    <v-btn
      v-if="error.details"
      @click="showErrorDetails = !showErrorDetails"
      class="my-2"
      text small
    >
      <span v-if="showErrorDetails">Hide details</span>
      <span v-else>Show details</span>
    </v-btn>
    <pre
      v-if="showErrorDetails"
      v-text="error.details"
      class="mt-4 box px-2 py-2"
    />
    <v-btn
      v-if="error.code === 405"
      @click="fetchProjectInfo(true)"
      class="my-2"
      dark rounded
    >
      <v-icon class="mr-2">layers_clear</v-icon>
      <!-- <v-icon class="mr-2">sync_problem</v-icon> -->
      <span>Skip layers with error</span>
    </v-btn>
  </v-layout>
</template>

<script>
import DefaultRedirect from '@/mixins/DefaultRedirect'
import GeneralInfo from './GeneralInfo'
import LayersInfo from './LayersInfo'

export default {
  name: "QgisProject",
  mixins: [ DefaultRedirect({ page: 'general' }) ],
  data () {
    return {
      error: null,
      showErrorDetails: false
    }
  },
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
  activated () {
    const unbind = this.$ws.bind('ProjectChanged', () => this.fetchProjectInfo())
    this.$once('hook:deactivated', unbind)
    if (!this.store.projectInfo) {
      this.fetchProjectInfo()
    }
  },
  methods: {
    fetchProjectInfo (skipLayersWithError=false) {
      if (!this.$ws.pluginConnected) {
        return
      }
      this.error = null
      const params = {
        skip_layers_with_error: skipLayersWithError
      }
      this.$ws.request('ProjectInfo', params)
        .then(resp => {
          this.store.projectInfo = resp.data
        })
        .catch(err => {
          this.store.projectInfo = null
          this.error = {
            code: err.status,
            msg: err.data || 'Error',
            details: err.traceback
          }
        })
    },
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
p {
  margin: 8px 0;
}
pre {
  font-size: 85%;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
