<script>
const svgIcons = require.context('../../icons', false, /.*\.svg$/)

function requireAll (requireContext) {
  return requireContext.keys().map(requireContext)
}
requireAll(svgIcons)

export default {
  functional: true,
  props: {
    name: String,
    color: String,
    size: String
  },
  render (h, ctx) {
    const { name, color, size } = ctx.props
    ctx.data.staticClass = ctx.data.staticClass ? ctx.data.staticClass + ' icon' : 'icon'
    const attrs = {
      'xlink:href': `#${name}`
    }
    const style = {}
    if (color) {
      if (color.startsWith('#')) {
        style.color = color
      } else {
        if (color.includes(' ')) {
          const [base, shade] = color.split(' ')
          ctx.data.staticClass += ` ${base}--text text--${shade}`
        } else {
          ctx.data.staticClass += ` ${color}--text`
        }
      }
    }
    if (size) {
      style.height = size
      // style.width = 'auto'
      style.maxWidth = '100%'
    }
    ctx.data.style = style
    const use = h('use', { attrs })
    return h('svg', ctx.data, [use])
  }
}
</script>

<style>
.icon {
  fill: currentColor;
  width: 1.5em;
  height: 1.5em;
}
</style>
