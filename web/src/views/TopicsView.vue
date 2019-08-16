<template>
  <div class="page py-2">
    <v-toolbar dark dense class="shrink" color="lime darken-2">
      <v-toolbar-title>Topics</v-toolbar-title>
    </v-toolbar>
    <v-layout class="column box list">
      <v-list class="grow" dense :ripple="false">
        <v-list-item
          v-for="(topic, index) in topics"
          :key="index"
          :class="{ 'v-list-item--active': selectedIndex === index}"
          @click="selectedIndex = index"
          @dblclick="editable = topic"
          :ripple="false"
        >
          <v-list-item-content>
            <span
              v-if="topic !== editable"
              v-text="topic.title"
              @click="selectedIndex = index"
              @dblclick="editable = topic"
            />
            <input
              v-else
              v-model="topic.title"
              @blur="editable = null"
              outlined
              hide-details
            />
            <!-- <v-text-field
              :disabled="topic !== editable"
              v-model="topic.title"
              @blur="editable = null"
              hide-details
            /> -->
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
    <v-layout class="column box abstract">
      <h3>Description</h3>
      <v-textarea
        :disabled="!activeTopic"
        :value="activeTopic && activeTopic.abstract"
        @input="activeTopic.abstract = $event"
        class="pt-0"
        hide-details
      />
    </v-layout>
    <v-layout class="column box layers">
      <h3>Layers</h3>
      <v-treeview
        v-if="activeTopic"
        :items="overlayLayers"
        item-key="name"
        item-children="layers"
        selectable
        v-model="activeTopic.visible_overlays"
        selected-color="grey darken-3"
        dense
      >
        <template v-slot:label="{ item }">
          <span>{{ item.title || item.name }}</span>
        </template>
      </v-treeview>
    </v-layout>
  </div>
</template>

<script>
export default {
  name: 'Topics',
  props: {
    config: Object
  },
  data () {
    return {
      selectedIndex: 0,
      editable: null
    }
  },
  computed: {
    topics () {
      return this.config.topics
    },
    overlayLayers () {
      return this.config.layers
    },
    activeTopic () {
      return this.topics[this.selectedIndex]
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
      this.config.topics.splice(this.selectedIndex, 1)
    }
  }
}
</script>

<style lang="scss" scoped>
.page {
  display: grid;
  grid-template-columns: 1fr 2fr;
  grid-template-rows: auto auto 1fr;
  grid-gap: 4px;
  // max-height: 100%;
  .v-toolbar {
    grid-row: 1 / 2;
    grid-column: 1 / 3;
  }
  .list {
    grid-row: 2 / 4;
    grid-column: 1 / 2;
    ::v-deep .v-toolbar__content {
      padding: 0;
    }
  }
  .abstract {
    grid-row: 2 / 3;
    grid-column: 2 / 3;
    // min-height: 150px;
  }
  .layers {
    grid-row: 3 / 4;
    grid-column: 2 / 3;
    overflow: auto;
  }
  .box {
    border: 1px solid #ccc;
    border-radius: 3px;
  }
}
.v-list {
  min-width: 240px;
}
</style>
