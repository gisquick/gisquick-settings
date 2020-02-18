<template>
  <v-app class="fill-height">
    <settings v-if="user && initialized"/>
    <login-dialog v-if="!user" @login="onLogin"/>
    <portal-target name="dialogs"/>
    <notifications ref="notification"/>
  </v-app>
</template>

<script>
import Vue from 'vue'
import WebsocketMessenger from '@/ws.js'
import Settings from '@/Settings.vue'
import LoginDialog from '@/components/LoginDialog.vue'
import Notifications from '@/components/Notifications.vue'

export default {
  name: 'app',
  components: { LoginDialog, Settings, Notifications },
  data () {
    return {
      initialized: false
    }
  },
  computed: {
    user () {
      return this.$root.user
    }
  },
  watch: {
    user: {
      immediate: true,
      handler (user) {
        if (user) {
          this.createWebsocketConnection()
        } else {
          if (this.$ws) {
            this.$ws.close()
            this.$ws = null
          }
        }
      }
    }
  },
  mounted () {
    Vue.prototype.$notification = this.$refs.notification
  },
  methods: {
    onLogin (user) {
      this.$root.user = user
    },
    createWebsocketConnection () {
      const protocol = location.protocol.endsWith('s:') ? 'wss' : 'ws'
      const ws = WebsocketMessenger(`${protocol}://${location.host}/ws/app`)
      Vue.util.defineReactive(ws, 'connected')
      Vue.util.defineReactive(ws, 'pluginConnected')
        ws.onopen().then(() => {
        this.initialized = true
      })
      Vue.prototype.$ws = ws
    }
  }
}
</script>

<style lang="scss">
html, body {
  /* font-size: 16px; */
  height: 100%;
  overflow: hidden!important;
}
.theme--light .icon {
  color: rgba(0, 0, 0, 0.54);
}
// .theme--dark .icon {
//   color: #fff;
// }
a {
  text-decoration: none;
}

.theme--light.v-text-field > .v-input__control > .v-input__slot:before {
  border-color: rgba(0, 0, 0, 0.15);
}
.v-divider--vertical {
  min-height: 0;
}
</style>
