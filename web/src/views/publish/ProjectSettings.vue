<template>
  <v-layout class="column box">
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
      <label><small>Extent:</small></label>
      <v-layout ml-2>
        <!-- <v-text-field
          label="Extent"
          :value="config.extent"
          @change="setExtent"
          :rules="[validators.extent]"
          class="mr-2"
        /> -->
        <v-text-field
          label="X-Min"
          v-model.number="config.extent[0]"
          type="number"
          class="mr-2"
        />
        <v-text-field
          label="Y-Min"
          v-model.number="config.extent[1]"
          type="number"
          class="mr-2"
        />
        <v-text-field
          label="X-Max"
          v-model.number="config.extent[2]"
          type="number"
          class="mr-2"
        />
        <v-text-field
          label="Y-Max"
          v-model.number="config.extent[3]"
          type="number"
          class="mr-2"
        />
        <v-select
          placeholder="Set layer extent"
          :items="layersExtents"
          @input="setExtent"
        />
      </v-layout>
      <v-checkbox
        label="Map cache"
        color="primary"
        v-model="config.use_mapcache"
      />
      <label>
        <small>Map Preview</small>
      </label>
      <map-preview
        :project="projectPath"
        :config="config"
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

function validateExtent (value) {
  if (!Array.isArray(value) || value.length !== 4) {
    return "Not valid extent"
  }
  return true
}

export default {
  name: 'ProjectSettings',
  components: {
    'map-preview': () => import('@/components/MapPreview')
  },
  props: {
    layers: Array,
    config: Object,
    projectPath: String
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
      return layersList(this.layers)
        .filter(l => l.extent)
        .map(l => ({
          text: l.name,
          value: roundExtent(l.extent)
        }))
    },
    validators () {
      return {
        extent: validateExtent
      }
    }
  },
  methods: {
    setExtent (extent) {
      if (typeof extent === 'string') {
        this.config.extent = extent.split(',').filter(v => !Number.isNaN(v)).map(parseFloat)
      } else {
        this.config.extent = extent
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.box {
  border: 1px solid #ccc;
  overflow: auto;
}
label {
  opacity: 0.65;
}
</style>
