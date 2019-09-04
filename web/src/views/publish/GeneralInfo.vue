<template>
  <v-layout column>
    <h4>Project CRS</h4>
    <v-text-field
      label="Code"
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
    <h4>Print templates</h4>
    <v-list class="mx-4 my-2">
      <v-list-item
        v-for="template in config.composer_templates"
        :key="template.name"
      >
        {{ template.name }}
      </v-list-item>
    </v-list>
  </v-layout>
</template>

<script>
export default {
  name: 'GeneralInfo',
  props: {
    config: Object
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
