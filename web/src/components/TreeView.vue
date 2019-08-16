<script>
import { VIcon, VLayout } from 'vuetify/lib'
import _xor from 'lodash/xor'

function Group (h, item, leafComp) {
  const children = item.children.map(item => item.children ? Group(h, item, leafComp) : leafComp({ item }))
  return (
    <div>
      <v-layout>
        <v-icon>folder</v-icon><span>{item.name}</span>
      </v-layout>
      <div class="children">
      {children}
      </div>
    </div>
  )
}

export default {
  components: { VIcon, VLayout },
  props: {
    items: Array
  },
  data () {
    return {
      opened: []
    }
  },
  methods: {
    toggleGroup (item) {
      this.opened = _xor(this.opened, [item.name])
    },
    renderGroup (h, item) {
      const open = this.opened.includes(item.name)
      const groupIcon = open ? 'mdi-folder-open' : 'mdi-folder'
      const children = item.children.map(item => item.children ? this.renderGroup(h, item) : this.$scopedSlots.leaf({ item }))
      return (
        <div>
          <v-layout>
            <v-icon class="mr-1" onClick={e => this.toggleGroup(item)}>{groupIcon}</v-icon><span>{item.name}</span>
          </v-layout>
          <div class="children">
          {open && children}
          </div>
        </div>
      )
    }
  },
  render (h) {
    const { items } = this.$props
    const leafs = items.filter(item => !item.children).map(item => this.$scopedSlots.leaf({ item }))
    const folders = items.filter(item => item.children)
    // const x = folders.map(item => Group(h, item, this.$scopedSlots.leaf))
    const x = folders.map(item => this.renderGroup(h, item))
    return (
      <v-layout class="column">
        {leafs}
        {x}
      </v-layout>
    )
  }
}
</script>

<style lang="scss" scoped>
.children {
  padding-left: 20px;
}
</style>
