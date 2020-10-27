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
      const val = this.attribute.content_type
      return val === undefined ? null : val
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
      // value = value !== null ? value : undefined
      if (value === null) {
        this.$delete(this.attribute, 'content_type')
      } else {
        this.$set(this.attribute, 'content_type', value)
      }
    }
  }
}
</script>
