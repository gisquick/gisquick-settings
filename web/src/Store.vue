<script>
import Vue from 'vue'

export default {
  data () {
    return {
      items: [],
      models: {}
    }
  },
  created () {
    Vue.prototype.$storage = {
      add: (id, component) => {
        if (!this.models[id]) {
          const ComponentClass = Vue.extend(component)
          const instance = new ComponentClass({ parent: this })
          Vue.set(this.models, id, instance)
        }
        // const exist = this.items.some(i => i.id === id)
        // if (!exist) {
        //   this.items.push({ id, component })
        // }
      },
      get: (id) => {
        // const item = this.items.find(i => i.id === id)
        // console.log(id, this.$refs[id])
        // return item && this.$refs[id]
        // return this.$refs[id]
        // console.log(this.$children.length)
        // return this.$children[0]
        return this.models[id]
      }
    }
  },
  render () {
    return null
  }
}
</script>
