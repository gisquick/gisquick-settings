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
      <template v-slot:item.content_type="{ item }">
        <attribute-content-type
          class="my-1 pt-0 mx-auto"
          :attribute="item"
          :value="attributesSettings[item.name].content_type"
          @input="setContentType(item, $event)"
        />
      </template>
    </v-data-table>

    <v-divider/>
    <components-select
      class="shrink align-start mx-2 mt-4"
      label="Info Panel"
      v-model="layer.infopanel_component"
    />
  </v-layout>
</template>

<script>
import AttributeContentType from '@/components/AttributeContentType.vue'
import ComponentsSelect from '@/components/ComponentsSelect.vue'
import { createMap, setNestedProperty, deleteNestedProperty } from '@/helpers'

function lookupTable (items) {
  // return items.reduce((lookup, v) => ({...lookup, [v]: true }), {})
  return items.reduce((lookup, v) => {
    return lookup[v] = true, lookup
  }, {})
}

export default {
  name: 'LayerAttributes',
  components: { AttributeContentType, ComponentsSelect },
  props: {
    layer: Object,
    layerSettings: Object
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
    attributesSettings () {
      return createMap(this.attributes, 'name', a => this.layerSettings?.attributes?.[a.name] || {})
    },
    columns () {
      return [
        {
          text: 'Name',
          value: 'name'
        }, {
          text: 'Alias',
          value: 'alias'
        }, {
          text: 'Type',
          value: 'type',
          sortable: false
        }, {
          text: 'Content Type',
          value: 'content_type',
          sortable: false
        }, {
          text: 'Attribute Table',
          value: 'attribute_table',
          sortable: false,
          align: 'center'
        }, {
          text: 'Info Panel',
          value: 'info_panel',
          sortable: false,
          align: 'center'
        }
      ]
    }
  },
  watch: {
    layerSettings: {
      immediate: true,
      handler (ls) {
        this.lookupAttr = lookupTable(ls.attr_table_fields || this.layer.attributes.map(a => a.name))
        this.lookupInfo = lookupTable(ls.info_panel_fields || this.layer.attributes.map(a => a.name))
      }
    },
    lookupAttr: {
      deep: true,
      handler (lookup) {
        this.$set(this.layerSettings, 'attr_table_fields', Object.keys(lookup).filter(k => lookup[k]))
      }
    },
    lookupInfo: {
      deep: true,
      handler (lookup) {
        this.$set(this.layerSettings, 'info_panel_fields', Object.keys(lookup).filter(k => lookup[k]))
      }
    }
  },
  methods: {
    setContentType (attr, value) {
      const key = `attributes.${attr.name}.content_type`
      if (value !== null) {
        setNestedProperty(this.layerSettings, key, value)
      } else {
        deleteNestedProperty(this.layerSettings, key)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.box {
  border: 1px solid #ccc;
  overflow: auto;
  background-color: #fff;
}
</style>
