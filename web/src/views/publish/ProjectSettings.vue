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
      <v-layout class="extent ml-2">
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
        <v-menu bottom left max-height="360">
          <template v-slot:activator="{ on }">
            <v-btn
              v-on="on"
              class="mt-2"
              icon
            >
              <v-icon>menu</v-icon>
            </v-btn>
          </template>
          <v-list>
            <v-list-item @click="enableDrawTool">Draw on map</v-list-item>
            <v-list-item
              v-if="allLayersExtent"
              @click="setExtent(allLayersExtent)"
            >
              All layers
            </v-list-item>
            <v-list-item-group>
              <v-layout row align-center @click.stop="">
                <v-divider/>
                <small class="mx-2 grey--text">Use layer extent</small>
                <v-divider/>
              </v-layout>
            </v-list-item-group>
            <v-list-item
              v-for="(layer, i) in extentLayers"
              :key="i"
              @click="setExtent(layer.extent)"
              v-text="layer.name"
            />
          </v-list>
        </v-menu>
      </v-layout>
      <label>
        <small>Map Preview</small>
      </label>
      <map-preview
        :project="projectPath"
        :config="config"
        ref="mapPreview"
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
import { extend } from 'ol/extent'
import Draw, { createBox } from 'ol/interaction/Draw'

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
    extentLayers () {
      return layersList(this.layers).filter(l => l.extent)
    },
    allLayersExtent () {
      if (this.extentLayers.length) {
        let extent = this.extentLayers[0].extent.slice()
        this.extentLayers.slice(1).forEach(layer => {
          extend(extent, layer.extent)
        })
        return extent
      }
      return null
    }
  },
  methods: {
    setExtent (extent) {
      this.config.extent = roundExtent(extent)
    },
    enableDrawTool () {
      const olMap = this.$refs.mapPreview.map
      const draw = new Draw({
        type: 'Circle',
        geometryFunction: createBox()
      })
      olMap.addInteraction(draw)
      draw.once('drawend', evt => {
        this.setExtent(evt.feature.getGeometry().getExtent())
        olMap.removeInteraction(draw)
      })
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
