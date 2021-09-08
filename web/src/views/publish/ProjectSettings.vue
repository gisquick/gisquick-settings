<template>
  <v-layout class="column box">
    <v-form class="px-3" ref="form">
      <v-text-field
        label="Title"
        placeholder=" "
        v-model="settings.title"
        :rules="validators.title"
      />
      <v-select
        label="Authentication"
        :items="authOpts"
        v-model="settings.authentication"
      />
      <div
        v-for="{ settingField, title } in extentFields"
        :key="settingField"
        class="extent-field"
      >
        <label><small v-text="title"/></label>
        <v-layout class="extent">
          <div class="marker" :class="settingField"/>
          <v-text-field
            label="X-Min"
            v-model.number="settings[settingField][0]"
            type="number"
            class="mr-2"
          />
          <v-text-field
            label="Y-Min"
            v-model.number="settings[settingField][1]"
            type="number"
            class="mr-2"
          />
          <v-text-field
            label="X-Max"
            v-model.number="settings[settingField][2]"
            type="number"
            class="mr-2"
          />
          <v-text-field
            label="Y-Max"
            v-model.number="settings[settingField][3]"
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
              <v-list-item @click="activateDraw(settingField)">
                <v-list-item-content>Draw on map</v-list-item-content>
                <v-list-item-icon>
                  <v-icon>picture_in_picture</v-icon>
                </v-list-item-icon>
              </v-list-item>
              <v-list-item
                v-if="allLayersExtent"
                @click="setExtent(settingField, allLayersExtent)"
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
                @click="setExtent(settingField, layer.extent)"
                v-text="layer.name"
              />
            </v-list>
          </v-menu>
        </v-layout>
      </div>
      <v-checkbox
        label="Custom default view extent"
        color="primary"
        :input-value="!!settings.zoom_extent"
        @change="toggleZoomExtent"
      />
      <div class="map-layout">
        <scales-list
          label="Scales"
          :value="settings.scales"
          @input="updateScales"
          @click:scale="zoomToScale"
          class="scales-list mt-1 shrink"
          :rules="validators.scales"
        />
        <v-layout class="map-preview column">
          <label>
            <small>Map Preview</small>
          </label>
          <map-preview
            ref="mapPreview"
            class="box grow"
            :project="projectPath"
            :config="config"
            :settings="settings"
          >
            <draw-extent v-if="drawExtentField" @draw="onExtentDraw"/>
          </map-preview>
        </v-layout>
      </div>
      <v-layout justify-space-between>
        <v-checkbox
          label="Map cache"
          color="primary"
          class="mt-2"
          v-model="settings.use_mapcache"
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
    projectPath: String,
    settings: Object
  },
  data () {
    return {
      drawExtentField: null,
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
    extentFields () {
      const fields = [{
        title: 'Map extent',
        settingField: 'extent'
      }]
      if (this.settings.zoom_extent) {
        fields.push({
          title: 'Default view extent',
          settingField: 'zoom_extent'
        })
      }
      return fields
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
    toggleZoomExtent (val) {
      this.$set(this.settings, 'zoom_extent', val ? [...this.settings.extent] : null)
    },
    activateDraw (field) {
      this.drawExtentField = field
    },
    setExtent (field, extent) {
      this.settings[field] = roundExtent(extent)
    },
    zoomToScale (scale) {
      const res = scalesToResolutions([scale], this.config.units)[0]
      const olMap = this.$refs.mapPreview.map
      olMap.getView().setCenter(getCenter(this.settings.extent))
      olMap.getView().setResolution(res)
    },
    updateScales (scales) {
      this.settings.scales = scales
      this.settings.tile_resolutions = scalesToResolutions(scales, this.config.units)
    },
    onExtentDraw (extent) {
      this.setExtent(this.drawExtentField, extent)
      this.drawExtentField = null
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
.map-layout {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 12px;
  justify-content: end;
  max-height: 700px;
  overflow: hidden;

  .scales-list {
    width: 220px;
  }
  .map-preview {
    max-height: 100%;
    min-height: 640px;
    margin-bottom: 22px;
  }
}

.extent-field {
  .marker {
    width: 4px;
    margin: 4px 8px 4px 0;
    height: 44px;
    &.extent {
      background-color: rgb(206, 55, 17);
    }
    &.zoom_extent {
      background-color: rgb(48, 136, 70);
    }
  }
}
</style>
