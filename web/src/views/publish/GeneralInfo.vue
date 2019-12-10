<template>
  <v-layout column>
    <h4>Project info</h4>
    <v-text-field
      label="Project file"
      :value="config.file"
      class="shrink mx-4"
      disabled
    />
    <h4>Map settings</h4>
    <v-text-field
      label="Projection"
      :value="config.projection.code"
      class="shrink mx-4"
      disabled
    />
    <v-text-field
      label="Units"
      :value="config.units"
      class="shrink mx-4"
      disabled
    />
    <v-text-field
      v-if="scalesText"
      label="Predefined scales"
      :value="scalesText"
      class="shrink mx-4"
      disabled
    />
    <v-layout v-else align-center pl-4 my-1 shrink>
      <v-icon color="orange">warning</v-icon>
      <span class="orange--text mx-2">Map scales are not defined</span>
      <small class="grey--text">(Project > Properties > General > Project Predefined Scales)</small>
    </v-layout>
    <h4 class="mt-1">Print templates</h4>
    <v-list
      v-if="config.composer_templates.length"
      class="mx-4 my-2"
    >
      <v-list-item
        v-for="template in config.composer_templates"
        :key="template.name"
      >
        {{ template.name }}
      </v-list-item>
    </v-list>
    <span v-else class="orange--text pl-4 my-1">None</span>
    <h4>Server options</h4>
    <v-checkbox
      label="Add geometry in feature response"
      class="mx-4"
      :input-value="config.server.wms_add_geometry"
      disabled
    />
  </v-layout>
</template>

<script>
export default {
  name: 'GeneralInfo',
  props: {
    config: Object
  },
  computed: {
    scalesText () {
      return this.config.scales.map(s => `1: ${s}`).join(', ')
    },
    error () {
      return false
      // return this.config.scales.length === 0
    },
    status () {
      if (this.error) {
        return 'error'
      }
      return 'ok'
    }
  },
  watch: {
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
h4 {
  color: #fff;

  padding: 0.25em 0.5em;
  margin: 0;
  background-color: #eee;
  border: 1px solid #ddd;
  border-width: 1px 0;
  color: #555;
  font-weight: 500;
}
.v-list {
  border: 1px solid #eee;
}
</style>
