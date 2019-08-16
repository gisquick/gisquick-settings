<script>
import Vue from 'vue'

export default {
  name: 'LayersStore',
  data () {
    return {
      layers: [],
      baseLayersNames: []
    }
  },
  computed: {
    overlayLayers () {
      return this.layers.filter(item => !this.baseLayersNames.includes(item.name))
    },
    baseLayers () {
      return this.layers.filter(item => this.baseLayersNames.includes(item.name))
    }
  },
  created () {
    this.$ws.bind('Layers', this.onMessage)
    this.$ws.send('Layers')
  },
  beforeDestroy () {
    this.$ws.removeEventListener('message', this.onMessage)
  },
  methods: {
    onMessage (e, msg) {
      const data = JSON.parse(msg)
      this.layers = data.layers
    }
  },
  render (h) {
    return null
  }
}
</script>
