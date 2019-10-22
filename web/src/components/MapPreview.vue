<template>
  <div class="map my-2"/>
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
import Fill from 'ol/style/Fill'
import Stroke from 'ol/style/Stroke'
import { fromExtent } from 'ol/geom/Polygon'
import { getCenter } from 'ol/extent'
import { get as getProj } from 'ol/proj'
import { register } from 'ol/proj/proj4'
import proj4 from 'proj4'

import 'ol/ol.css'

import { layersList } from '@/utils'


export default {
  props: {
    project: String,
    config: Object
  },
  computed: {
    extent () {
      return this.config.extent
    }
  },
  mounted () {
    this.createMap()
  },
  watch: {
    project () {
      if (this.map) {
        this.map.dispose()
      }
    },
    config (config) {
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

      const layers = layersList(config.overlays)
        .filter(l => l.visible)
        .map(l => l.serverName || l.name)

      let projection = getProj(config.projection.code)
      if (!projection) {
        proj4.defs(config.projection.code, config.projection.proj4)
        register(proj4)
        projection = getProj(config.projection.code)
      }
      const layer = new ImageLayer({
        visible: true,
        // extent: config.extent,
        source: new ImageWMS({
          url: url,
          // resolutions: this.resolutions,
          params: {
            LAYERS: layers.join(','),
            FORMAT: 'image/png',
            TRANSPARENT: 'false'
          },
          serverType: 'qgis',
          ratio: 1.0
        })
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
      this.map = new Map({
        layers: [layer, extentLayer],
        view: new View({
          projection: projection,
          center: getCenter(config.extent),
          zoom: 0,
          // resolutions: config.tile_resolutions,
        }),
        // controls: []
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
  max-width: 900px;
  position: relative;
  background-color: #fff;
  border: 1px solid #ddd;
}
</style>
