<template>
  <div class="map">
    <v-fade-transition>
      <div v-if="loading" class="progressbar px-2 py-1">
        <label class="mr-2"><small>Loading</small></label>
        <v-progress-circular size="20" width="2" indeterminate/>
      </div>
    </v-fade-transition>
    <div v-if="error" class="error-msg red--text">
      <small>Failed to render map</small>
    </div>
    <slot/>
  </div>
</template>

<script>
import Map from 'ol/Map'
import View from 'ol/View'
import ImageLayer from 'ol/layer/Image'
import ImageWMS from 'ol/source/ImageWMS'
import VectorSource from 'ol/source/Vector'
import VectorLayer from 'ol/layer/Vector'
import Feature from 'ol/Feature'
import Style from 'ol/style/Style'
import Stroke from 'ol/style/Stroke'
import { fromExtent } from 'ol/geom/Polygon'
import { getCenter } from 'ol/extent'
import { get as getProj } from 'ol/proj'
import { register } from 'ol/proj/proj4'
import { defaults } from 'ol/interaction'
import MouseWheelZoom from 'ol/interaction/MouseWheelZoom'
import proj4 from 'proj4'

import 'ol/ol.css'

import { layersList } from '@/utils'


export default {
  props: {
    project: String,
    config: Object
  },
  data () {
    return {
      error: false,
      loading: false
    }
  },
  computed: {
    extent () {
      return this.config.extent
    }
  },
  mounted () {
    this.createMap()
  },
  beforeDestroy () {
    if (this.map) {
      this.map.dispose()
    }
  },
  watch: {
    project () {
      if (this.map) {
        this.map.dispose()
      }
    },
    config () {
      if (this.map) {
        this.map.dispose()
      }
      this.createMap()
    },
    extent (extent) {
      const f = this.extentSource.getFeatures()[0]
      f.setGeometry(fromExtent(extent))
    }
  },
  methods: {
    createMap () {
      const config = this.config

      // const timestamp = config.publish_date_unix
      // const project = timestamp ? `${this.project}_${timestamp}` : this.project
      const project = config.project || this.project
      const url = `/api/project/map?MAP=${project}.qgs`

      const layers = layersList(config.overlays).filter(l => l.visible)
      layers.sort((a, b) => (b.drawing_order || 0) - (a.drawing_order || 0))

      let projection = getProj(config.projection.code)
      if (!projection) {
        proj4.defs(config.projection.code, config.projection.proj4)
        register(proj4)
        projection = getProj(config.projection.code)
      }
      const source = new ImageWMS({
        url: url,
        params: {
          LAYERS: layers.map(l => l.serverName || l.name).join(','),
          FORMAT: 'image/png',
          TRANSPARENT: 'false'
        },
        serverType: 'qgis',
        ratio: 1.0
      })
      const mapLayer = new ImageLayer({ source })
      source.on('imageloadstart', () => {
        this.loading = true
      })
      source.on('imageloadend', () => {
        this.error = false
        this.loading = false
      })
      source.on('imageloaderror', () => {
        this.error = true
        this.loading = false
      })
      this.extentSource = new VectorSource({
        features: [new Feature({ geometry: fromExtent(config.extent) })]
      })
      const extentLayer = new VectorLayer({
        source: this.extentSource,
        style: new Style({
          stroke: new Stroke({
            color: [3, 169, 244, 0.8],
            width: 2
          }),
          fill: null
        })
      })
      const zoomInteraction = new MouseWheelZoom({
        condition: evt => evt.type === 'wheel' && evt.originalEvent.ctrlKey
      })
      this.map = new Map({
        layers: [mapLayer, extentLayer],
        view: new View({
          projection: projection,
          center: getCenter(config.extent),
          zoom: 0,
          // resolutions: config.tile_resolutions,
        }),
        interactions: defaults({mouseWheelZoom: false}).extend([zoomInteraction])
      })
      this.map.setTarget(this.$el)
      this.map.getView().fit(this.config.extent, { padding: [50, 50, 50, 50] })
    }
  }
}
</script>

<style lang="scss" scoped>
.map {
  height: 500px;
  // width: 100%;
  position: relative;
  background-color: #fff;

  .progressbar {
    position: absolute;
    top: 0;
    right: 0;
    opacity: 0.75;
    z-index: 100;
    background-color: rgba(255, 255, 255, 0.6);
  }
  .error-msg {
    position: absolute;
    top: 45%;
    text-align: center;
    width: 100%;
  }
}
</style>
