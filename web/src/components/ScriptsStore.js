export default {

  data () {
    return {
      resource: null
    }
  },
  computed: {
    namespace () {
      const { user, folder } = this.$route.params
      return user && `${user}/${folder}`
    },
    scripts () {
      return this.resource?.data || {}
    },
    loading () {
      return this.resource?.fetching
    },
    components () {
      const list = []
      Object.entries(this.scripts).forEach(([m, i]) => {
        list.push({
          header: m,
        })
        list.push(...i.components)
      })
      return list
    }
  },
  watch: {
    namespace () {
      this.resource = null
    }
  },
  created () {
    this.$root.scriptsStore = this
  },
  methods: {
    getScripts () {
      const namespace = this.namespace
      if (!this.resource) {
        const resource = {
          fetching: true,
          data: {},
          error: null,
        }
        this.resource = resource
        this.$http.get(`/api/project/script/${namespace}`)
          .then(resp => {
            this.resource.fetching = false
            resource.data = resp.data
          })
          .catch(err => {
            resource.error = err
          })
      }
      return this.resource
    },
    update (scripts) {
      this.resource.data = scripts
    }
  },
  render () {
    return null
  }
}