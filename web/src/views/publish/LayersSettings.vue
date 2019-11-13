<template>
  <v-layout>
    <layers-table
      label="Base layers"
      :items="config.base_layers"
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
          color="grey darken-1"
          class="my-0 py-1 justify-center"
          :ripple="false"
          hide-details
        />
      </template>
    </layers-table>

    <layers-table
      label="Overlay layers"
      :items="config.overlays"
      :headers="overlaysHeaders"
      :opened.sync="opened"
      :selected="selected"
      selected-class="amber lighten-4"
      class="ml-2 mr-1 my-1 elevation-2"
      @click:row="onClick"
    >
      <template v-slot:leaf.publish="{ item }">
        <v-checkbox
          v-model="item.publish"
          color="grey darken-1"
          class="my-0 py-1 justify-center"
          :ripple="false"
          hide-details
        />
      </template>
      <template v-slot:leaf.hidden="{ item }">
        <v-checkbox
          v-model="item.hidden"
          color="grey darken-1"
          class="my-0 py-1 justify-center"
          :ripple="false"
          hide-details
        />
      </template>
      <template v-slot:leaf.attributes="{ item }">
        <v-btn text>
          <icon name="attribute-table" size="18"/>
        </v-btn>
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
    order: Array
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
        // }, {
        //   text: 'Attributes',
        //   value: 'attributes'
        }
      ]
    },
    selectedLayer () {
      if (!this.selected) {
        return null
      }
      const { base_layers, overlays } = this.config
      return base_layers.find(l => l.name === this.selected) || overlays.find(l => l.name === this.selected)
    },
    canSwapSelectedLayer () {
      return this.selectedLayer && !specialBaseLayersTypes.includes(this.selectedLayer.type)
    }
  },
  methods: {
    onClick (item) {
      if (this.config.base_layers.includes(item) || this.config.overlays.includes(item)) {
        this.selected = item.name
      }
    },
    swapSelectedLayer () {
      if (this.config.base_layers.find(i => i.name === this.selected)) {
        var src = this.config.base_layers
        var dest = this.config.overlays
      } else {
        src = this.config.overlays
        dest = this.config.base_layers
      }
      const layer = src.find(i => i.name === this.selected)
      src.splice(src.indexOf(layer), 1)
      const destIndex = this.order
        .filter(name => name === this.selected || dest.find(i => i.name === name))
        .indexOf(this.selected)
      dest.splice(destIndex, 0, layer)
    }
  }
}
</script>

<style lang="scss" scoped>
.swap.v-btn {
  position: absolute;
  top: 50%;
  left: 50%;
  min-width: 36px;
  padding: 0;
  transform: translate(-50%, -50%);
}
</style>
