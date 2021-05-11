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
    this.pageVisible = true
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
