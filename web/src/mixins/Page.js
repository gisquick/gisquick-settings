export default {
  data () {
    return {
      pageVisible: false
    }
  },
  created () {
    // console.log('created')
  },
  mounted () {
    // console.log('mounted')
  },
  activated () {
    // console.log('activated')
    this.pageVisible = true
  },
  deactivated () {
    // console.log('deactivated')
    this.pageVisible = false
  }
}
