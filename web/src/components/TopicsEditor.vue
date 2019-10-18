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
        <v-treeview
          v-if="activeTopic"
          :items="layers"
          item-key="name"
          item-children="layers"
          item-disabled="hidden"
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
      <label v-else>No topic is selected</label>
    </v-layout>
  </div>
</template>

<script>
export default {
  name: 'TopicsEditor',
  props: {
    topics: Array,
    layers: Array
  },
  data () {
    return {
      selectedIndex: 0
    }
  },
  computed: {
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
      this.topics.splice(this.selectedIndex, 1)
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
</style>
