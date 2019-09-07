<template>
  <v-layout column>
    <v-form class="px-3">
      <v-text-field
        label="Title"
        v-model="config.title"
      />
      <v-select
        label="Authentication"
        :items="authOpts"
        v-model="config.authentication"
      />
      <v-checkbox
        label="Map cache"
        color="primary"
        v-model="config.use_mapcache"
      />
    </v-form>
    <v-spacer/>
  </v-layout>
</template>

<script>
export default {
  name: 'ProjectSettings',
  props: {
    config: Object
  },
  data () {
    return {
    }
  },
  computed: {
    authOpts () {
      return [
        {
          text: 'All',
          value: 'all'
        }, {
          text: 'All authenticated',
          value: 'authenticated'
        }, {
          text: 'Only project owner',
          value: 'owner'
        }
      ]
    }
  }
}
</script>

<style lang="scss" scoped>

</style>
