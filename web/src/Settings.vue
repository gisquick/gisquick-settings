<template>
  <v-layout column>
    <v-store/>

    <v-app-bar app dark dense>
      <img src="./assets/text_logo_dark.svg" class="logo"/>
      <router-view name="menu"/>

      <img :src="pluginStatusImg" class="qgis-status mr-2"/>
      <v-menu>
        <template v-slot:activator="{ on: menu }">
          <v-btn v-on="{ ...menu }" icon>
            <v-icon>menu</v-icon>
          </v-btn>
        </template>
        <v-list dense class="py-1">
          <v-list-item @click="logout">
            <v-list-item-title>Logout</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>

    <keep-alive>
      <router-view v-if="user && $ws.connected" class="content"/>
    </keep-alive>
  </v-layout>
</template>

<script>
import Vue from 'vue'
import VStore from '@/Store.vue'
import WebsocketMessenger from '@/ws.js'

export default {
  components: { VStore },
  computed: {
    user () {
      return this.$root.user
    },
    pluginStatusImg () {
      return this.$ws.pluginConnected ? require('./assets/qgis-icon32.svg') : require('./assets/qgis-icon-black32.svg')
    }
  },
  created () {
    console.log('Settings:Created')
    const protocol = location.protocol.endsWith('s:') ? 'wss' : 'ws'
    const socket = WebsocketMessenger(`${protocol}://${location.host}/ws/app`)
    Vue.util.defineReactive(socket, 'connected')
    Vue.util.defineReactive(socket, 'pluginConnected')
    Vue.prototype.$ws = socket
  },
  beforeDestroy () {
    this.$ws.close()
    this.$ws = null
  },
  methods: {
    logout () {
      this.$http.get('/logout/').then(() => location.reload())
    }
  }
}
</script>

<style lang="scss" scoped>
.logo {
  height: calc(100% - 20px);
  box-sizing: content-box;
}
.qgis-status {
  margin-left: auto;
}
.content {
  margin: 48px auto 8px auto;
  display: flex;
  height: calc(100vh - 48px);
  overflow: auto;

  @media (max-width: 960px) {
    width: 100%;
  }
  @media (min-width: 1440px) {
    width: 1200px;
  }
}
</style>
