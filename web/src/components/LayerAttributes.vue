<template>
  <v-layout class="box column">
    <v-data-table
      :headers="columns"
      :items="attributes"
      disable-pagination
      hide-default-footer
      class="grow"
    >
      <template v-slot:item.attribute_table="{ item }">
        <v-simple-checkbox
          color="secondary"
          v-model="lookupAttr[item.name]"
          :ripple="false"
        />
      </template>
      <template v-slot:item.info_panel="{ item }">
        <v-simple-checkbox
          color="secondary"
          v-model="lookupInfo[item.name]"
          :ripple="false"
        />
      </template>
      <template v-slot:item.widget="{ item }">
        <attribute-widget
          :attribute="item"
          class="my-1 pt-0"
        />
      </template>
    </v-data-table>
  </v-layout>
</template>

<script>
import AttributeWidget from '@/components/AttributeWidget.vue'

function lookupTable (items) {
  // return items.reduce((lookup, v) => ({...lookup, [v]: true }), {})
  return items.reduce((lookup, v) => {
    return lookup[v] = true, lookup
  }, {})
}

export default {
  components: { AttributeWidget },
  props: {
    layer: Object
  },
  data () {
    return {
      lookupAttr: {},
      lookupInfo: {}
    }
  },
  computed: {
    attributes () {
      return this.layer.attributes
    },
    columns () {
      return [
        {
          text: 'Name',
          sortable: true,
          value: 'name'
        }, {
          text: 'Alias',
          value: 'alias'
        }, {
          text: 'Type',
          value: 'type'
        }, {
          text: 'Display',
          value: 'widget'
        }, {
          text: 'Attribute Table',
          value: 'attribute_table'
        }, {
          text: 'Info Panel',
          value: 'info_panel'
        }
      ]
    }
  },
  watch: {
    layer: {
      immediate: true,
      handler (layer) {
        this.lookupAttr = lookupTable(layer.attr_table_fields || layer.attributes.map(a => a.name))
        this.lookupInfo = lookupTable(layer.info_panel_fields || layer.attributes.map(a => a.name))
      }
    },
    lookupAttr: {
      deep: true,
      handler (lookup) {
        this.$set(this.layer, 'attr_table_fields', Object.keys(lookup).filter(k => lookup[k]))
      }
    },
    lookupInfo: {
      deep: true,
      handler (lookup) {
        this.$set(this.layer, 'info_panel_fields', Object.keys(lookup).filter(k => lookup[k]))
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.box {
  border: 1px solid #ccc;
  overflow: auto;
}
</style>
