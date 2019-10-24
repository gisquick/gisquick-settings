<template>
  <v-layout column>
    <v-list dense class="grow">
      <v-list-item
        v-for="(scale, index) in value"
        :key="index"
        :input-value="selected === scale"
        @click="selected = scale"
        @dblclick="setEditable(index)"
        color="primary"
        :ripple="false"
      >
        <span class="mr-2">1:</span>
        <span
          v-if="index !== editingIndex"
          class="mr-2"
          v-text="scale"
        />
        <input
          v-else
          ref="editField"
          type="number"
          class="mr-2"
          v-model.number="editValue"
          @blur="finishEdit"
          @keydown.enter="finishEdit"
        >
        <v-spacer/>
        <v-list-item-action class="my-2" @click="$emit('click:scale', scale)">
          <v-icon>search</v-icon>
        </v-list-item-action>
      </v-list-item>
    </v-list>
    <v-divider/>
    <v-layout shrink>
      <v-btn text class="flex" @click="add">
        <v-icon>add</v-icon>
      </v-btn>
      <v-btn
        text
        class="flex"
        :disabled="selected === null"
        @click="remove(selected)"
      >
        <v-icon>remove</v-icon>
      </v-btn>
    </v-layout>
  </v-layout>
</template>

<script>
import { nextTick } from 'q'
export default {
  props: {
    value: Array
  },
  data () {
    return {
      selected: null,
      editingIndex: null,
      editValue: 1 
    }
  },
  methods: {
    add () {
      this.$emit('input', this.value.concat(1))
    },
    remove (scale) {
      const scales = this.value.filter(s => s !== scale)
      this.$emit('input', scales)
    },
    setEditable (index) {
      this.editValue = this.value[index]
      this.editingIndex = index
      nextTick(() => {
        this.$refs.editField[0].focus()
      })
    },
    finishEdit () {
      const scales = this.value.slice()
      scales[this.editingIndex] = this.editValue
      scales.sort((a, b) => b - a)
      this.$emit('input', scales)
      this.editingIndex = null
    }
  }
}
</script>

<style lang="scss" scoped>
.v-list-item {
  .v-list-item__action {
    opacity: 0;
    transition: 0.25s opacity ease;
    will-change: opacity;
  }
  &:hover {
    .v-list-item__action {
      opacity: 1;
    }
  }
  input {
    min-width: 0;
  }
}
</style>
