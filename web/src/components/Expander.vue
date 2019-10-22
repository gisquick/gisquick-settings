
<script>

const text = Array(200).fill('.').join(' ')

export default {
  functional: true,
  render (h) {
    return <span class="expander">{text}</span>
  }
}
</script>

<style scoped>
.expander {
  line-height: 0;
  visibility: hidden;
}
</style>
