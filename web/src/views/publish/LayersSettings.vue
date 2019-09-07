<template>
  <v-layout>
    <layers-table
      label="Base layers"
      :items="config.base_layers"
      :headers="[]"
      :opened.sync="opened"
      class="mx-1 my-1 elevation-2"
      @click:row="onClick"
    >

    </layers-table>

    <layers-table
      label="Overlay layers"
      :items="config.layers"
      :headers="overlaysHeaders"
      :opened.sync="opened"
      :selected="selected"
      selected-class="amber lighten-4"
      class="mx-1 my-1 elevation-2"
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
    </layers-table>
  </v-layout>
</template>

<script>
import LayersTable from '@/components/LayersTable'

export default {
  name: 'LayersSettings',
  components: { LayersTable },
  props: {
    config: Object
  },
  data () {
    return {
      opened: [],
      selected: null
    }
  },
  computed: {
    overlaysHeaders () {
      return [
        {
          text: 'Publish',
          value: 'publish'
        }, {
          text: 'Hidden',
          value: 'hidden'
        }
      ]
    }
  },
  methods: {
    onClick (item) {
      console.log('Click', item)
      if (this.config.layers.includes(item)) {
        this.selected = item.name
      }
    }
  }
}
</script>

<style lang="scss" scoped>
</style>
