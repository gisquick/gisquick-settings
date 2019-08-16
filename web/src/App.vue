<template>
  <v-app>
    <settings v-if="user"/>
    <login-dialog v-else @login="onLogin"/>
    <portal-target name="dialogs"/>
  </v-app>
</template>

<script>
import Settings from '@/Settings.vue'
import LoginDialog from '@/components/LoginDialog.vue'

export default {
  name: 'app',
  components: { LoginDialog, Settings },
  computed: {
    user () {
      return this.$root.user
    }
  },
  methods: {
    onLogin (user) {
      this.$root.user = user
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
</style>
