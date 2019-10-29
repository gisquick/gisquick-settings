<template>
  <v-input
    :value="value"
    :rules="rules"
    class="layout column"
    :class="{
      'primary--text': focus,
      'v-input--is-focused': focus
    }"
  >
    <small slot="label">Scales</small>
    <v-layout class="column grow box">
      <v-list
        dense
        class="grow"
        @click.native="focus = true"
      >
        <v-list-item
          v-for="(scale, index) in value"
          :key="index"
          :input-value="selectedIndex === index"
          @click="selectedIndex = index"
          @dblclick="setEditable(index)"
          color="primary"
          :ripple="false"
        >
          <span class="mr-1">1:</span>
          <span
            v-if="index !== editingIndex"
            class="pl-1 mr-1"
            v-text="scale"
          />
          <input
            v-else
            ref="editField"
            type="number"
            class="pl-1 mr-1 edit-field"
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
          :disabled="selectedIndex === null"
          @click="removeAt(selectedIndex)"
        >
          <v-icon>remove</v-icon>
        </v-btn>
      </v-layout>
    </v-layout>
  </v-input>
</template>

<script>
export default {
  props: {
    value: Array,
    rules: Array,
    label: String
  },
  data () {
    return {
      focus: false,
      selectedIndex: null,
      editingIndex: null,
      editValue: 1 
    }
  },
  methods: {
    add () {
      this.$emit('input', this.value.concat(1))
    },
    removeAt (index) {
      const scales = this.value.filter((s, i) => i !== index)
      this.$emit('input', scales)
    },
    setEditable (index) {
      this.editValue = this.value[index]
      this.editingIndex = index
      this.$nextTick(() => {
        this.$refs.editField[0].focus()
      })
    },
    finishEdit () {
      const scales = this.value.slice()
      scales[this.editingIndex] = this.editValue
      scales.sort((a, b) => b - a)
      this.$emit('input', scales)
      this.editingIndex = null
    },
    onKeyDown (evt) {
      if (evt.target.tagName === 'INPUT') {
        return
      }
      if (evt.key === 'Delete') {
        this.removeAt(this.selectedIndex)
        evt.preventDefault()
      } else if (evt.key === 'ArrowUp') {
        if (this.selectedIndex > 0) {
          this.selectedIndex--
        }
        evt.preventDefault()
      } else if (evt.key === 'ArrowDown') {
        if (this.selectedIndex < this.value.length - 1) {
          this.selectedIndex++
        }
        evt.preventDefault()
      }
    },
    _clickOutsideTest (e) {
      if (!this.$el.contains(e.target)) {
        this.focus = false
      }
    }
  },
  beforeDestroy () {
    document.removeEventListener('keydown', this.onKeyDown)
    document.removeEventListener('keydown', this._clickOutsideTest)
  },
  watch: {
    focus (focus) {
      if (focus) {
        document.addEventListener('keydown', this.onKeyDown, false)
        document.addEventListener('click', this._clickOutsideTest)
      } else {
        document.removeEventListener('keydown', this.onKeyDown)
        document.removeEventListener('click', this._clickOutsideTest)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.v-input {
  color: #ddd;
  &.error--text, &.v-input--is-focused {
    .box {
      border-color: currentColor;
    }
    ::v-deep {
      .v-label {
        color: inherit;
      }
    }
  }
}
::v-deep {
  .v-label {
    transition: none;
    width: 100%;
  }
  .v-input__slot {
    flex-direction: column;
    flex-grow: 1;
  }
  .v-messages {
    flex-grow: 0;
  }
}
.box {
  border: 1px solid;
  width: 100%;
}
.v-list-item {
  input.edit-field {
    opacity: 0.8;
    color: inherit;
    outline: auto;
  }
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
