<script>
import { VIcon, VLayout } from 'vuetify/lib'
import _xor from 'lodash/xor'

const DefaultValue = {
  functional: true,
  props: ['value'],
  render (h, ctx) {
    return <span>{ctx.props.value}</span>
  }
}

function layerIcon (layer) {
  if (layer.type === 'vector') {
    return {
      POINT: 'point',
      LINE: 'line',
      POLYGON: 'polygon'
    }[layer.geom_type]
  } else {
    return 'raster'
  }
}

export default {
  name: 'LayersTable',
  components: { VIcon, VLayout },
  props: {
    label: String,
    items: Array,
    headers: Array,
    opened: {
      type: Array,
      default: () => []
    },
    selected: String,
    selectedClass: String
  },
  methods: {
    toggleGroup (item) {
      // this.opened = _xor(this.opened, [item.name])
      this.$emit('update:opened', _xor(this.opened, [item.name]))
    },
    renderLeaf (h, item, depth) {
      let cmp
      const indentStyle = {
        paddingLeft: `${28 * depth}px`
      }
      if (this.$scopedSlots.leaf) {
        cmp = this.$scopedSlots.leaf({ item, depth })
        cmp[0].data.style = indentStyle
      } else {
        cmp = (
          <v-layout
            shrink align-center
            style={indentStyle} onClick={() => this.$emit('click:row', item)}
          >
            <icon name={layerIcon(item)} class="grey--text mr-2"/>
            <span>{item.name}</span>
          </v-layout>
        )
      }

      const values = this.headers.map(header => {
        const slot = this.$scopedSlots[`leaf.${header.value}`]
        const comp = slot ? slot({ item, depth }) : h(DefaultValue, { props: { value: item[header.value] } })
        return <td>{comp}</td>
      })
      return (
        <tr class={{[this.selectedClass]: item.name === this.selected}}>
          <td>{cmp}</td>
          {values}
        </tr>
      )
    },
    renderGroup (h, item, depth) {
      const open = this.opened.includes(item.name)
      const groupIcon = open ? 'mdi-folder-open' : 'mdi-folder'
      let children = []
      if (open) {
        children = item.layers.map(item => item.layers ? this.renderGroup(h, item, depth + 1) : this.renderLeaf(h, item, depth + 1))
      }
      const paddingStyle = {
        paddingLeft: `${28 * depth}px`,
      }
      const groupNode = (
        <tr>
          <td
            colspan={this.headers.length + 1}
            class="group-header py-1"
            onClick={e => this.toggleGroup(item)}
            key={item.name}
          >
            <v-layout>
              <v-icon class="mr-1" style={paddingStyle}>{groupIcon}</v-icon>
              <span>{item.name}</span>
            </v-layout>
          </td>
        </tr>
      )
      return [groupNode, ...children]
    }
  },
  render (h) {
    const { items } = this.$props
    const children = items.map(item => item.layers ? this.renderGroup(h, item, 0) : this.renderLeaf(h, item, 0))
    const headers = [
      <th class="header text-left">{this.label}</th>,
      ...this.headers.map(header => <th class={`header text-${header.align || 'center'}`} width="1">{header.text}</th>)
    ]
    return (
      <div class="table-container">
        <table class="layers-table">
          <thead>
            <tr>{headers}</tr>
          </thead>
          <tbody>
            {children}
          </tbody>
        </table>
      </div>
    )
  }
}
</script>

<style lang="scss" scoped>
.table-container {
  overflow: auto;
  flex: 1 1;
}
.layers-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  th, td {
    padding: 0 12px;
  }
  .header {
    background-color: #eee;
    padding-top: 8px;
    padding-bottom: 8px;
    border-top: 1px solid #ddd;
    border-bottom: 1px solid #ddd;
    font-weight: 500;
    // opacity: 0.75;
    font-size: 14px;
    position: sticky;
    top: 0;
    z-index: 1;
    // box-sizing: border-box;
  }
  .group-header {
    cursor: pointer;
    user-select: none;
    height: 36px;
    span {
      font-weight: 500;
      // line-height: 28px;
      // display: inline-block;
      // margin-top: 4px;
    }
    i {
      // vertical-align: center;
    }
  }
  tbody {
    tr:hover {
      background-color: #eee;
    }
  }
  .v-input--checkbox {
    ::v-deep {
      .v-input--selection-controls__input {
        margin-right: 0;
      }
    }
  }
}
</style>