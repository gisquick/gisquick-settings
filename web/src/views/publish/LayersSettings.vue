<template>
  <v-layout class="layers-settings">
    <layers-table
      label="Base layers"
      :items="baseLayers"
      :headers="baseLayersHeaders"
      :opened.sync="opened"
      :selected="selected"
      selected-class="amber lighten-4"
      class="ml-1 mr-2 my-1 elevation-2"
      @click:row="onClick"
    >
      <template v-slot:leaf.publish="{ item }">
        <v-checkbox
          v-model="item.publish"
          color="secondary"
          class="my-0 py-1 justify-center"
          :ripple="false"
          hide-details
        />
      </template>
    </layers-table>

    <layers-table
      label="Layers"
      :items="overlayLayers"
      :headers="overlaysHeaders"
      :opened.sync="opened"
      :selected="selected"
      selected-class="amber lighten-4"
      class="ml-2 mr-1 my-1 elevation-2"
      @click:row="onClick"
    >
      <template v-slot:header.layer="{}">
        <router-link to="layers/permissions">
          <v-icon class="px-1 py-1">vpn_lock</v-icon>
        </router-link>
      </template>
      <template v-slot:leaf.publish="{ item }">
        <v-checkbox
          v-model="settings.layers[item.id].publish"
          color="secondary"
          class="my-0 py-1 justify-center"
          :ripple="false"
          hide-details
        />
      </template>
      <template v-slot:leaf.hidden="{ item }">
        <v-checkbox
          v-model="settings.layers[item.id].hidden"
          color="secondary"
          class="my-0 py-1 justify-center"
          :ripple="false"
          hide-details
        />
      </template>
      <template v-slot:leaf.attributes="{ item }">
        <router-link
          class="layout row justify-center"
          :to="`layers/${item.name}/attributes`"
        >
          <icon name="attribute-table" size="18"/>
        </router-link>
      </template>
    </layers-table>

    <v-btn
      rounded
      class="swap grey lighten-2"
      :disabled="!canSwapSelectedLayer"
      :outlined="!canSwapSelectedLayer"
      @click="swapSelectedLayer"
    >
      <v-icon>swap_horiz</v-icon>
    </v-btn>
  </v-layout>
</template>

<script>
import LayersTable from '@/components/LayersTable'

const specialBaseLayersTypes = ['blank', 'osm', 'mapbox', 'bing']

export default {
  name: 'LayersSettings',
  components: { LayersTable },
  props: {
    config: Object,
    order: Array,
    settings: Object
  },
  data () {
    return {
      opened: [],
      selected: null
    }
  },
  computed: {
    baseLayersHeaders () {
      return [
        {
          text: 'Publish',
          value: 'publish'
        }
      ]
    },
    overlaysHeaders () {
      return [
        {
          text: 'Publish',
          value: 'publish'
        }, {
          text: 'Hidden',
          value: 'hidden'
        }, {
          text: 'Attributes',
          value: 'attributes'
        }
      ]
    },
    baseLayers () {
      return this.config.layers.filter(l => this.settings.base_layers.includes(l.id || l.name))
    },
    overlayLayers () {
      return this.config.layers.filter(l => this.settings.overlays.includes(l.id || l.name))
    },
    selectedLayer () {
      // return this.selected && this.config.layers.find(l => l.id ? l.id === this.selected : l.name === this.selected)
      return this.selected
    },
    canSwapSelectedLayer () {
      return this.selectedLayer && !specialBaseLayersTypes.includes(this.selectedLayer.type)
    }
  },
  methods: {
    onClick (item) {
      // this.selected = item.id || item.name
      this.selected = item
      // if (this.config.base_layers.find(l => l.name === item.name) || this.config.overlays.find(l => l.name === item.name)) {
        // this.selected = item.name
      // }
    },
    swapSelectedLayer () {
      const value = this.selected.id || this.selected.name
      const isBaseLayer = this.settings.base_layers.includes(value)
      const src = isBaseLayer ? 'base_layers' : 'overlays'
      const dest = isBaseLayer ? 'overlays' : 'base_layers'
      this.settings[src] = this.settings[src].filter(i => i !== value)
      this.settings[dest].push(value)
    }
  }
}
</script>

<style lang="scss" scoped>
.layers-settings {
  position: relative;
  overflow: hidden;
  max-height: 100%;
  flex: 1 1;
  .table-container {
    overflow: auto;
  }
}
.swap.v-btn {
  position: absolute;
  top: 50%;
  left: 50%;
  min-width: 36px;
  padding: 0;
  transform: translate(-50%, -50%);
}
</style>
