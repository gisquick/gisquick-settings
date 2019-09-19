<template>
  <v-layout column>
    <v-form class="px-3">
      <v-text-field
        label="Title"
        v-model="config.title"
      />
      <v-layout>
        <v-text-field
          label="Extent"
          v-model="config.extent"
          class="mr-2"
        />
        <v-select
          placeholder="Set layer extent"
          :items="layersExtents"
          @input="setExtent"
        />
      </v-layout>
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
import { layersList } from '@/utils'

function roundExtent (extent) {
  const maxVal = Math.max(...extent.map(v => Math.abs(v)))
  if (maxVal > 1000) {
    return extent.map(v => Math.round(v))
  }
  return extent
}

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
    },
    layersExtents () {
      return layersList(this.config.layers)
        .filter(l => l.extent)
        .map(l => ({
          text: l.name,
          value: roundExtent(l.extent)
        }))
    }
  },
  methods: {
    setExtent (extent) {
      this.config.extent = extent
    }
  }
}
</script>

<style lang="scss" scoped>

</style>
