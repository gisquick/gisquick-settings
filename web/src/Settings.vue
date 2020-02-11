<template>
  <div class="app-layout">

    <div class="header-bg elevation-3"/>
    <div class="logo pl-2">
      <img src="./assets/text_logo_dark.svg"/>
    </div>

    <portal-target name="menu" class="menu">
      <themeable-toolbar dark>
        <v-btn rounded text :to="{name: 'projects'}">
          <v-icon class="mr-1">view_list</v-icon>
          Projects
        </v-btn>
        <portal-target name="menu-breadcrumbs" multiple/>
        <v-spacer/>
        <portal-target name="menu-actions"/>
      </themeable-toolbar>
    </portal-target>

    <v-layout class="app-menu align-center justify-end pr-2">
      <img :src="pluginStatusImg" class="mr-2"/>
      <v-menu>
        <template v-slot:activator="{ on: menu }">
          <v-btn dark icon height="40" width="40" v-on="{ ...menu }">
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

    <keep-alive>
      <router-view/>
    </keep-alive>
  </div>
</template>

<script>
import Themeable from 'vuetify/lib/mixins/themeable'

const ThemeableToolbar = {
  mixins: [Themeable],
  render (h) {
    return <div class="layout row align-center mx-0">{this.$slots.default}</div>
  }
}

export default {
  components: { ThemeableToolbar },
  computed: {
    user () {
      return this.$root.user
    },
    pluginStatusImg () {
      return this.$ws.pluginConnected ? require('./assets/qgis-icon32.svg') : require('./assets/qgis-icon-black32.svg')
    }
  },
  methods: {
    logout () {
      this.$http.get('/api/auth/logout/').then(() => location.reload())
    }
  }
}
</script>

<style lang="scss" scoped>
.app-layout {
  width: 100%;
  height: 100%;
  display: grid;
  grid-template-columns: minmax(auto, 1fr) auto 1fr;
  grid-template-rows: 48px 1fr;
  @media (max-width: 1450px) {
    grid-template-columns: auto 1fr auto;
  }
  .header-bg {
    background-color: #3d3d3d;
    grid-column: 1 / 4;
    grid-row: 1 / 2;
  }
  .logo {
    grid-row: 1 / 2;
    grid-column: 1 / 2;
    max-width: 250px;
    width: 100%;
    height: inherit;
    img {
      height: inherit;
      padding: 10px 0;
    }
  }
  .menu {
    grid-row: 1 / 2;
    grid-column: 2 / 3;
    display: flex;
    ::v-deep {
      .v-btn.v-btn--text {
        padding: 0 8px;
        &:before {
          display: none;
        }
        // text-transform: none;
      }
    }
  }
  .app-menu {
    grid-row: 1 / 2;
    grid-column: 3 / 4;
    min-width: 92px;
    img {
      height: 24px;
      width: 24px;
    }
  }
}
</style>
