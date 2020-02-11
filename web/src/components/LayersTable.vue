<script>
import { VIcon, VLayout, VBtn, VTextField } from 'vuetify/lib'
import _xor from 'lodash/xor'

import { layersGroups, filterLayers, searchRegex } from '@/utils'

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

/* Filter that test/match also groups */
export function filterLayers2(items, test) {
  const list = []
  items.forEach(item => {
    if (item.layers) {
      if (test(item)) {
        list.push(item)
      } else {
        const children = filterLayers2(item.layers, test)
        if (children.length) {
          list.push({
            ...item,
            layers: children
          })
        }
      }
    } else if (test(item)) {
      list.push(item)
    }
  })
  return list
}

export default {
  name: 'LayersTable',
  components: { VIcon, VLayout, VBtn, VTextField },
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
  data () {
    return {
      filter: ''
    }
  },
  computed: {
    displayedItems () {
      if (this.filter) {
        const regex = searchRegex(this.filter)
        return filterLayers(this.items, l => l.name.search(regex) !== -1)
      }
      return this.items
    }
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
            style={indentStyle}
            onClick={() => this.$emit('click:row', item)}
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
        <tr class={{[this.selectedClass]: item.name === this.selected}}>
          <td
            colspan={this.headers.length + 1}
            class="group-header py-1"
            onClick={() => this.$emit('click:row', item)}
            key={item.name}
          >
            <v-layout>
              <v-icon
                class="mr-1"
                style={paddingStyle}
                vOn:click_stop={() => this.toggleGroup(item)}
              >
                {groupIcon}
              </v-icon>
              <span>{item.name}</span>
            </v-layout>
          </td>
        </tr>
      )
      return [groupNode, ...children]
    },
    expandAll () {
      const groups = layersGroups(this.items)
      this.$emit('update:opened', groups.map(g => g.name))
    },
    renderLayersHeader (h) {
      const slot = this.$scopedSlots['header.layer']
      return (
        <th class="header">
          <v-layout class="align-center">
            <v-icon onClick={this.expandAll}>mdi-folder-open</v-icon>
            <span class="mx-2">{this.label}</span>
            <span class="spacer"/>
            <v-text-field
              placeholder="Filter"
              class="filter my-0 py-0 mx-2"
              append-icon="search"
              rounded={true}
              clearable={true}
              hide-details={true}
              value={this.filter}
              onInput={v => this.filter = v}
            />
            {slot ? slot() : null}
          </v-layout>
        </th>
      )
    }
  },
  render (h) {
    const items = this.displayedItems
    const children = items.map(item => item.layers ? this.renderGroup(h, item, 0) : this.renderLeaf(h, item, 0))
    const headers = [
      this.renderLayersHeader(h),
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
    white-space: nowrap;
  }
  th {
    .filter {
      max-width: 320px;
      background-color: #fff;
    }
  }
  td {
    height: 32px;
  }
  .header {
    background-color: #eee;
    padding-top: 4px;
    padding-bottom: 4px;
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
