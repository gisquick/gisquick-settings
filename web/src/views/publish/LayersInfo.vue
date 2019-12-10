<template>
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
        <v-spacer/>
        <v-tooltip v-if="sourceInfo[item.name].error" bottom color="red darken-1">
          <template v-slot:activator="{ on }">
            <v-icon v-on="on" color="red darken-1">report</v-icon>
          </template>
          <span v-text="sourceInfo[item.name].error"/>
        </v-tooltip>
      </v-layout>
    </template>
  </layers-table>
</template>

<script>
import LayersTable from '@/components/LayersTable'
import { layersList, layersGroups } from '@/utils'
import path from 'path'

export default {
  name: 'Layers',
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
      const clientWinPlatform = /windows/i.test(this.$ws.clientInfo)
      layersList(this.layers).forEach(l => {
        let text = l.source
        const schema = l.source.split('://')[0]
        const info = { icon: this.sourceIcons[schema] }
        if (schema === 'file') {
          text = text.substring(7)
          if (clientWinPlatform) {
            text = text.replace(/^[/]+/, '') // strip '/' from the beginning
          }
          if (text.startsWith(projectDir)) {
            text = text.substring(projectDir.length)
          } else {
            info.error = 'Data file is located outside of project directory'
          }
        }
        info.text = text
        sources[l.name] = info
      })
      return sources
    },
    status () {
      const error = Object.values(this.sourceInfo).some(s => s.error)
      return error ? 'error' : 'ok'
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
