<template>
  <div class="layers-view py-2">
    <v-toolbar dark dense class="shrink" color="lime darken-2">
      <v-toolbar-title>Base layers</v-toolbar-title>
    </v-toolbar>
    <v-toolbar dark dense class="shrink" color="lime darken-2">
      <v-toolbar-title>Overlay layers</v-toolbar-title>
    </v-toolbar>
    <v-treeview
      :items="baseLayers"
      class="groups py-1"
      open-all
      dense
      rounded
      hoverable
      activatable
      item-key="name"
      :active="selection"
      @update:active="selectItem"
    >
      <template v-slot:prepend="{ item, open }">
        <v-icon v-if="item.children">
          {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
        </v-icon>
        <icon
          v-else-if="item.type === 'vector'"
          :name="geomIcons[item.geom_type]"
        />
      </template>
    </v-treeview>
    <v-treeview
      :items="overlayLayers"
      class="layers py-1"
      open-all
      dense
      rounded
      hoverable
      activatable
      item-key="name"
      item-children="layers"
      :active="selection"
      @update:active="selectItem"
    >
      <template v-slot:prepend="{ item, open }">
        <v-icon v-if="item.children">
          {{ open ? 'mdi-folder-open' : 'mdi-folder' }}
        </v-icon>
        <icon
          v-else-if="item.type === 'vector'"
          :name="geomIcons[item.geom_type]"
        />
      </template>
      <template v-slot:append="{ item, open }">
        <v-layout shrink align-center ml-4>
          <!-- <icon
            v-if="item.type === 'vector'"
            :name="geomIcons[item.geom_type]"
            class="mx-2"
          /> -->
          <v-checkbox v-if="!item.children" color="primary" class="my-0" hide-details/>
        </v-layout>
      </template>
    </v-treeview>
    <v-btn
      rounded
      class="switch my-1"
      @click="switchSelected"
    >
      <v-icon class="mr-2">swap_horiz</v-icon>
      <span>Switch</span>
    </v-btn>
  </div>
</template>

<script>
import _xor from 'lodash/xor'
import LayersStore from '@/LayersStore.vue'

export default {
  name: 'Layers',
  props: {
    config: Object
  },
  data () {
    return {
      // layers: null,
      // baseLayersNames: [],
      selectedItem: null,
      selection: []
    }
  },
  computed: {
    store () {
      return this.$storage.get('layers')
    },
    layers () {
      return this.config.base_layers.concat(this.config.layers)
    },
    baseLayers () {
      return this.config.base_layers
    },
    overlayLayers () {
      return this.config.layers
    },
    geomIcons () {
      return {
        POINT: 'point',
        LINE: 'line',
        POLYGON: 'polygon'
      }
    }
  },
  created () {
    this.$storage.add('layers', LayersStore)
  },
  methods: {
    selectItem (items) {
      const item = items[0]
      if (item) {
        const isRootNode = this.layers.some(l => l.name === item)
        this.selection = isRootNode ? items : []
      }
    },
    switchSelected () {
      if (this.selection.length) {
        this.store.baseLayersNames = _xor(this.store.baseLayersNames, this.selection)
        this.selection = [...this.selection] // clone to update UI
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.box {
  border: 1px solid #aaa;
}
.layers-view {
  max-height: 100%;
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: auto 1fr auto;
  grid-column-gap: 4px;

  .switch {
    grid-column: 1 / 3;
    max-width: 200px;
    align-self: center;
    justify-self: center;
  }
}
.v-treeview {
  min-width: 360px;
  border: 1px solid #ddd;
  overflow: auto;
  user-select: none;
  ::v-deep .v-treeview-node__root {
    margin: 0;
  }
  .v-treeview-node__root.v-treeview-node--active .v-treeview-node__content .icon {
    color: inherit;
  }
}
</style>
