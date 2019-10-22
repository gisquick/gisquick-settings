<template>
  <!-- <tree-table
    v-if="config"
    :items="config.layers"
    :headers="headers"
    class="tree"
  >
    <template v-slot:leaf="{ item, depth }">
      <v-layout shrink align-center>
        <icon
          v-if="item.type === 'vector'"
          :name="geomIcons[item.geom_type]"
          class="grey--text mr-2"
        />
        <span>{{ item.name }}</span>
      </v-layout>
    </template>
    <template v-slot:leaf.visible="{ item }">
      <v-checkbox
        v-model="item.visible"
        color="grey darken-1"
        class="my-0 py-1 justify-center"
        :ripple="false"
        hide-details
      />
    </template>
    <template v-slot:leaf.hidden="{ item }">
      <v-checkbox
        v-model="item.hidden"
        :disabled="!item.visible"
        color="grey darken-1"
        class="my-0 py-1 justify-center"
        :ripple="false"
        hide-details
      />
    </template>
    <template v-slot:leaf.queryable="{ item }">
      <span>?</span>
    </template>
  </tree-table> -->
  <layers-table
    label="Layer"
    :items="layers"
    :headers="headers"
    :opened.sync="opened"
  >
    <template v-slot:leaf="{ item, depth }">
      <v-layout shrink align-center>
        <icon
          v-if="item.type === 'vector'"
          :name="geomIcons[item.geom_type]"
          class="grey--text mr-2"
        />
        <icon
          v-else-if="item.type === 'raster'"
          name="raster"
          class="grey--text mr-2"
        />
        <span>{{ item.name }}</span>
      </v-layout>
    </template>
    <!-- <template v-slot:leaf.publish="{ item }">
      <v-checkbox
        v-model="item.publish"
        color="grey darken-1"
        class="my-0 py-1 justify-center"
        :ripple="false"
        hide-details
      />
    </template> -->
    <template v-slot:leaf.wfs="{ item }">
      <v-icon v-if="item.wfs">check</v-icon>
      <v-icon v-else>remove</v-icon>
    </template>
    <template v-slot:leaf.source="{ item }">
      <v-layout
        v-if="item.source"
        class="align-center source"
      >
        <v-icon size="20" class="mr-1" v-text="sourceInfo[item.name].icon"/>
        <span>{{ sourceInfo[item.name].text }}</span>
      </v-layout>
    </template>
  </layers-table>
</template>

<script>
// import TreeView from '@/components/TreeView'
// import TreeTable from '@/components/TreeTable'
import LayersTable from '@/components/LayersTable'
import { layersList, layersGroups } from '@/utils'
import path from 'path'

export default {
  name: 'Layers',
  // components: { TreeView, TreeTable, LayersTable },
  components: { LayersTable },
  props: {
    config: Object
  },
  data () {
    return {
      opened: []
    }
  },
  computed: {
    layers () {
      return this.config.layers
    },
    headers () {
      return [
        // {
        //   text: 'Publish',
        //   value: 'publish',
        // },
        {
          text: 'WFS',
          value: 'wfs'
        },
        {
          text: 'CRS',
          value: 'projection',
          align: 'left'
        },
        {
          text: 'Provider',
          value: 'provider_type',
          align: 'left'
        },
        {
          text: 'Source',
          value: 'source',
          align: 'left'
        }
      ]
    },
    geomIcons () {
      return {
        POINT: 'point',
        LINE: 'line',
        POLYGON: 'polygon'
      }
    },
    sourceIcons () {
      return {
        file: 'mdi-file-document-outline',
        postgresql: '$vuetify.icons.db',
        http: 'public',
        https: 'public'
      }
    },
    sourceInfo () {
      const sources = {}
      const projectDir = this.config.directory + path.sep
      layersList(this.layers).forEach(l => {
        let text = l.source
        const schema = l.source.split('://')[0]
        if (schema === 'file') {
          text = text.substring(7)
          if (text.startsWith(projectDir)) {
            text = text.substring(projectDir.length)
          }
        }
        sources[l.name] = {
          icon: this.sourceIcons[schema],
          text
        }
      })
      return sources
    },
    status () {
      return 'ok'
    }
  },
  watch: {
    layers: {
      immediate: true,
      handler (layers) {
        this.opened = layersGroups(layers).map(l => l.name)
      }
    },
    status: {
      immediate: true,
      handler (status) {
        this.$emit('status', status)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.tree {
  font-size: 15px;
  // max-width: 550px;
  border: 1px solid #ddd;
  ::v-deep {
    .header {
      position: sticky;
      top: 0;
      z-index: 1;
    }
  }
  ::v-deep > * {
    align-self: center;
    padding: 0 8px;
  }
}
.layers-table {
  .source span {
    max-width: 400px;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }
}
</style>
