<template>
  <v-select
    :items="widgetOptions"
    :disabled="disabled"
    :value="value || null"
    @input="updateValue"
    item-value="value"
    hide-details
  />
</template>

<script>
export default {
  props: {
    attribute: Object,
    value: String
  },
  computed: {
    dataType () {
      return this.attribute.type.split('(')[0]
    },
    disabled () {
      return this.dataType !== 'TEXT'
    },
    widgetOptions () {
      return [
        {
          text: 'Plain value',
          value: null
        }, {
          text: 'Hyperlink',
          value: 'url'
        }, {
          text: 'Media image file',
          value: 'media;image/*'
        }
      ]
    }
  },
  methods: {
    updateValue (value) {
      this.$emit('input', value)
    }
  }
}
</script>
