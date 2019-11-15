<script>
import Draw, { createBox } from 'ol/interaction/Draw'

export default {
  mounted () {
    const olMap = this.$parent.map
    const draw = new Draw({
      type: 'Circle',
      geometryFunction: createBox()
    })
    olMap.addInteraction(draw)
    draw.on('drawend', evt => this.$emit('draw', evt.feature.getGeometry().getExtent()))
    this.$once('hook:beforeDestroy', () => olMap.removeInteraction(draw))
  },
  render: h => null
}
</script>
