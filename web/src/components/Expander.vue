
<script>

const text = Array(200).fill('.').join(' ')

export default {
  functional: true,
  render (h) {
    return <div class="expander">{text}</div>
  }
}
</script>

<style scoped>
.expander {
  line-height: 0;
  height: 0;
  visibility: hidden;
}
</style>
