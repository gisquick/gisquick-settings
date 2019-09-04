<script>
export default {
  functional: true,
  props: {
    name: String
  },
  render (h, ctx) {
    const attrs = {
      'xlink:href': `#${ctx.props.name}`
    }
    const use = h('use', { attrs })
    return h('svg', ctx.data, [use])
  }
}
</script>

<style lang="scss" scoped>
.v-icon {
  fill: currentColor;
}
</style>
