<template>
  <v-layout class="column box">
    <v-form class="px-3" ref="form">
      <v-text-field
        label="Title"
        placeholder=" "
        v-model="config.title"
        :rules="validators.title"
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
        <v-menu bottom max-height="400" min-width="220">
          <template v-slot:activator="{ on }">
            <v-btn
              v-on="on"
              class="mt-2"
              icon
            >
              <v-icon>menu</v-icon>
            </v-btn>
          </template>
          <v-list dense>
            <v-list-item @click="drawToolActive = true">
              <v-list-item-content>Draw on map</v-list-item-content>
              <v-list-item-icon>
                <v-icon>picture_in_picture</v-icon>
              </v-list-item-icon>
            </v-list-item>
            <v-list-item
              v-if="allLayersExtent"
              @click="setExtent(allLayersExtent)"
            >
              <v-list-item-content>Extent of all layers</v-list-item-content>
              <v-list-item-icon>
                <v-icon>zoom_out_map</v-icon>
              </v-list-item-icon>
            </v-list-item>
            <v-list-item-group>
              <v-layout row align-center @click.stop="">
                <v-divider/>
                <small class="mx-2 my-1 grey--text">Use layer extent</small>
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
      <v-layout row ml-0 mr-2 justify-space-between>
        <scales-list
          label="Scales"
          :value="config.scales"
          @input="updateScales"
          @click:scale="zoomToScale"
          class="scales-list mt-1 shrink mr-2"
          :rules="validators.scales"
        />
        <v-layout class="map-preview column">
          <label>
            <small>Map Preview</small>
          </label>
          <map-preview
            :project="projectPath"
            :config="config"
            class="box grow"
            ref="mapPreview"
          >
            <draw-extent v-if="drawToolActive" @draw="onExtentDraw"/>
          </map-preview>
        </v-layout>
      </v-layout>
      <v-layout justify-space-between>
        <v-checkbox
          label="Map cache"
          color="primary"
          class="mt-2"
          v-model="config.use_mapcache"
        />
        <v-btn
          text
          outlined
          color="red darken-2"
          @click="deleteCache"
        >
          Delete cache
        </v-btn>
      </v-layout>
    </v-form>
    <v-spacer/>
  </v-layout>
</template>

<script>
import { extend, getCenter } from 'ol/extent'
import { required } from '@/validators'
import { layersList, scalesToResolutions } from '@/utils'
import ScalesList from '@/components/ScalesList'
const MapPreview = () => import(/* webpackChunkName: "olmap" */ '@/components/MapPreview');
const DrawExtent = () => import(/* webpackChunkName: "olmap" */ '@/components/DrawExtent');

function roundExtent (extent) {
  const maxVal = Math.max(...extent.map(v => Math.abs(v)))
  if (maxVal > 1000) {
    return extent.map(v => Math.round(v))
  }
  return extent
}

export default {
  name: 'ProjectSettings',
  components: { ScalesList, MapPreview, DrawExtent },
  props: {
    layers: Array,
    config: Object,
    projectPath: String
  },
  data () {
    return {
      drawToolActive: false
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
    },
    validators () {
      return {
        title: [required],
        scales: [v => v.length === 0 ? 'Map scales must be defined' : true]
      }
    }
  },
  mounted () {
    // this.$refs.form.validate()
  },
  methods: {
    setExtent (extent) {
      this.config.extent = roundExtent(extent)
    },
    zoomToScale (scale) {
      const res = scalesToResolutions([scale], this.config.units)[0]
      const olMap = this.$refs.mapPreview.map
      olMap.getView().setCenter(getCenter(this.config.extent))
      olMap.getView().setResolution(res)
    },
    updateScales (scales) {
      this.config.scales = scales
      this.config.tile_resolutions = scalesToResolutions(scales, this.config.units)
    },
    onExtentDraw (extent) {
      this.setExtent(extent)
      this.drawToolActive = false
    },
    deleteCache () {
      this.$http.delete(`/api/project/cache/${this.projectPath}`)
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
.scales-list {
  width: 200px;
}
.map-preview {
  max-width: 900px;
  margin-bottom: 22px;
}
</style>
