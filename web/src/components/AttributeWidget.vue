<template>
  <v-select
    :items="widgetOptions"
    :disabled="disabled"
    :value="value"
    @input="updateValue"
    item-value="value"
    hide-details
  />
</template>

<script>
export default {
  props: {
    attribute: Object
  },
  computed: {
    dataType () {
      return this.attribute.type.split('(')[0]
    },
    disabled () {
      return this.dataType !== 'TEXT'
    },
    value () {
      return this.attribute.widget || null
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
          value: {
            type: 'media',
            accept: 'image/*'
          }
        }
      ]
    }
  },
  methods: {
    updateValue (val) {
      const value = val !== null ? val : undefined
      // this.attribute.widget = value
      this.$set(this.attribute, 'widget', value)
    }
  }
}
</script>
