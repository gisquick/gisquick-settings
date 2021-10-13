<template>
  <div class="topics">
    <v-toolbar
      dark
      color="grey darken-2"
      height="42"
    >
      <v-toolbar-title>Topics</v-toolbar-title>
    </v-toolbar>
    <v-layout class="column box list">
      <v-list class="grow" dense :ripple="false">
        <v-list-item
          v-for="(topic, index) in topics"
          :key="index"
          :class="{ 'v-list-item--active': selectedIndex === index}"
          @click="selectedIndex = index"
          :ripple="false"
        >
          <v-list-item-content>
            <span
              v-text="topic.title"
              @click="selectedIndex = index"
            />
          </v-list-item-content>
        </v-list-item>
      </v-list>
      <v-divider/>
      <v-toolbar dense flat class="shrink">
        <v-btn class="mx-1 px-1 grow" @click="addTopic">
          <v-icon>add_circle</v-icon>
          <span>Add</span>
        </v-btn>
        <v-btn class="mx-1 px-1 grow" @click="removeSelectedTopic">
          <v-icon>delete</v-icon>
          <span>Remove</span>
        </v-btn>
      </v-toolbar>
    </v-layout>
    <v-layout class="form box px-2 py-2">
      <v-layout column v-if="activeTopic">
        <v-text-field
          label="Title"
          v-model="activeTopic.title"
        />
        <v-textarea
          label="Description"
          v-model="activeTopic.abstract"
          placeholder=" "
          rows="2"
          auto-grow
        />
        <label>
          <small>Layers</small>
        </label>
        <layers-table
          v-if="activeTopic"
          label="Layer"
          :items="availableLayers"
          :headers="[]"
          :opened.sync="opened"
        >
          <template v-slot:leaf="{ item }">
            <!-- :disabled="settings.layers[item.id].hidden" -->
            <v-checkbox
              color="secondary"
              class="my-0 py-1"
              :label="item.name"
              :ripple="false"
              :input-value="topicLayersLookup[item.id]"
              @change="toggleLayer(item)"
              hide-details
            />
          </template>
        </layers-table>
      </v-layout>
      <label v-else>No topic is selected</label>
    </v-layout>
  </div>
</template>

<script>
import LayersTable from '@/components/LayersTable.vue'
import { filterLayers, lookupTable } from '@/utils'

export default {
  name: 'TopicsEditor',
  components: { LayersTable },
  props: {
    topics: Array,
    layers: Array,
    settings: Object
  },
  data () {
    return {
      selectedIndex: 0,
      opened: []
    }
  },
  computed: {
    availableLayers () {
      return filterLayers(this.layers, l => !this.settings.layers[l.id].hidden)
    },
    activeTopic () {
      return this.topics[this.selectedIndex]
    },
    topicLayersLookup () {
      return lookupTable(this.activeTopic.visible_overlays)
    }
  },
  methods: {
    addTopic () {
      this.topics.push({
        title: 'New',
        abstract: '',
        visible_overlays: []
      })
    },
    removeSelectedTopic () {
      // if (this.selectedIndex >= 0) {
      this.topics.splice(this.selectedIndex, 1)
    },
    toggleLayer (layer) {
      if (this.topicLayersLookup[layer.id]) {
        this.activeTopic.visible_overlays = this.activeTopic.visible_overlays.filter(id => id != layer.id)
      } else {
        this.activeTopic.visible_overlays.push(layer.id)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.topics {
  display: grid;
  grid-template-columns: 1fr 2fr;
  grid-template-rows: auto 1fr;
  // max-height: 100%;
  .v-toolbar {
    grid-row: 1 / 2;
    grid-column: 1 / 3;
  }
  .list {
    grid-row: 2 / 4;
    grid-column: 1 / 2;
    margin-right: 4px;
    ::v-deep .v-toolbar__content {
      padding: 0;
    }
  }
  .form {
    grid-row: 2 / 3;
    grid-column: 2 / 3;
    overflow: auto;
    .v-input {
      flex-grow: 0;
    }
  }
  .box {
    border: 1px solid #ccc;
    border-radius: 3px;
  }
  label {
    opacity: 0.75;
  }
}
.v-list {
  min-width: 240px;
}
.layers-table {
  .v-input--checkbox {
    ::v-deep {
      label {
        margin-left: 4px;
        color: #111;
      }
    }
  }
}
</style>
