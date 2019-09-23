<template>
  <v-layout column fill-height>
    <v-store/>

    <v-app-bar dark dense class="shrink">
      <div class="logo">
        <img src="./assets/text_logo_dark.svg"/>
      </div>
      <!-- <router-link to="/" class="logo">
        <img src="./assets/text_logo_dark.svg"/>
      </router-link> -->

      <portal-target name="menu" class="menu layout row">
        <!-- <v-divider vertical class="my-2 mr-1"/> -->
        <v-btn rounded text :to="{name: 'projects'}">
          <v-icon class="mr-1">view_list</v-icon>
          Projects
        </v-btn>
        <portal-target name="menu-breadcrumbs" multiple/>
        <v-spacer/>
        <portal-target name="menu-actions"/>
      </portal-target>

      <v-layout justify-end>
        <img :src="pluginStatusImg" class="mr-2"/>
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
      </v-layout>
    </v-app-bar>

    <keep-alive>
      <router-view v-if="user && $ws.connected" class="fill-height"/>
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
  height: inherit;
  img {
    height: inherit;
    padding: 10px 0;
  }
}
::v-deep .content {
  margin: 0 auto;
  display: flex;
  position: relative;
  overflow: auto;
  min-width: 960px;

  @media (max-width: 960px) {
    min-width: 0;
    width: 100%;
  }
  @media (min-width: 1440px) {
    width: 1200px;
  }
}
.v-app-bar {
  ::v-deep .v-toolbar__content {
    display: grid;
    grid-template-columns: 1fr auto 1fr;
  }
  .menu {
    flex: 1 0;
    ::v-deep .v-btn--text {
      padding: 0 0.5em
    }
    @media (max-width: 960px) {
      min-width: 0;
      width: auto;
    }
    @media (min-width: 1440px) {
      width: 1200px;
    }
    ::v-deep {
      .v-btn.v-btn--text {
        &:before {
          display: none;
        }
        // text-transform: none;
      }
    }
  }
}
</style>
