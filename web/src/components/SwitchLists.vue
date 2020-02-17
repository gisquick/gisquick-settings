<template>
  <div class="switch-lists">
    <v-layout class="box column grid-l">
      <v-layout class="shrink align-center px-2 grey lighten-2">
        <label class="text-center mr-4">Asigned</label>
        <v-text-field
          v-model="valueFilter"
          class="white py-0 my-1"
          append-icon="search"
          placeholder="Filter"
          clearable
          rounded
          hide-details
        />
      </v-layout>
      <v-divider/>
      <input-container
        @keydown.delete="removeSelected"
        class="layout scrollable"
      >
        <v-list dense class="grow py-0">
          <v-list-item
            v-for="(item, index) in displayedValueItems"
            :key="index"
            @click="selected = item"
            :input-value="selected === item"
            color="primary"
            :ripple="false"
          >
            <slot name="item" :item="item">{{ item[itemText] }}</slot>
          </v-list-item>
        </v-list>
      </input-container>
    </v-layout>

    <v-btn
      icon
      class="switch mx-1"
      :disabled="!selected"
      @click="switchItem"
    >
      <v-icon>swap_horizontal_circle</v-icon>
    </v-btn>

    <v-layout class="box column grid-r">
      <v-layout class="shrink align-center px-2 grey lighten-4">
        <label class="text-center mr-4">Available</label>
        <v-text-field
          v-model="itemsFilter"
          class="white my-1 py-0"
          append-icon="search"
          placeholder="Filter"
          clearable
          rounded
          hide-details
        />
      </v-layout>
      <v-divider/>
      <input-container class="layout scrollable">
        <v-list dense class="grow py-0">
          <v-list-item
            v-for="(item, index) in displayedAvailableItems"
            :key="index"
            @click="selected = item"
            :input-value="selected === item"
            color="primary"
            :ripple="false"
          >
            <slot name="item" :item="item">{{ item[itemText] }}</slot>
          </v-list-item>
        </v-list>
      </input-container>
    </v-layout>
    <v-btn
      text small rounded
      class="text-none my-1"
      @click="removeAll"
    >
      Remove all
    </v-btn>
  </div>
</template>

<script>
import InputContainer from '@/components/InputContainer'

import { searchRegex } from '@/utils'

export default {
  components: { InputContainer },
  props: {
    value: Array,
    items: Array,
    itemText: {
      type: String,
      default: 'text'
    },
    itemValue: String
  },
  data () {
    return {
      selected: null,
      valueFilter: '',
      itemsFilter: ''
    }
  },
  computed: {
    itemsMap () {
      const map = {}
      if (this.itemValue) {
        this.items.forEach(i => { map[i[this.itemValue]] = i })
      }
      return map
    },
    valueItems () {
      if (!this.itemValue) {
        return this.value
      }
      return this.value.map(v => this.itemsMap[v])
    },
    displayedValueItems () {
      if (this.valueFilter) {
        const regex = searchRegex(this.valueFilter)
        return this.valueItems.filter(i => i[this.itemText].search(regex) !== -1)
      }
      return this.valueItems
    },
    displayedAvailableItems () {
      let items = this.items.filter(i => !this.valueItems.includes(i))
      if (this.itemsFilter) {
        const regex = searchRegex(this.itemsFilter)
        items = items.filter(i => i[this.itemText].search(regex) !== -1)
      }
      return items
    }
  },
  methods: {
    _itemValue (item) {
      return this.itemValue ? item[this.itemValue] : item
    },
    switchItem () {
      if (this.selected) {
        if (this.displayedValueItems.includes(this.selected)) {
          this.$emit('input', this.value.filter(i => i !== this._itemValue(this.selected)))
          // this.$emit('input', this.valueItems.filter(i => i !== this.selected).map(this._itemValue))
        } else if (this.displayedAvailableItems.includes(this.selected)) {
          this.$emit('input', this.value.concat([this._itemValue(this.selected)]))
        }
      }
    },
    removeSelected () {
      if (this.selected) {
        this.$emit('input', this.value.filter(i => i !== this._itemValue(this.selected)))
      }
    },
    removeAll () {
      this.$emit('input', [])
    }
  }
}
</script>

<style lang="scss" scoped>
.switch-lists {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  grid-template-rows: 1fr auto;
  overflow: hidden;
  min-height: 0;
  height: 100%;

  .grid-l {
    grid-column: 1 / 2;
  }
  .grid-r {
    grid-column: 3 / 4;
  }
  .v-btn {
    &.switch {
      align-self: center;
    }
  }
  label {
    font-size: 14px;
    font-weight: 500;
  }
}
.box {
  border: 1px solid #ccc;
  background-color: #fff;
  border-radius: 3px;
  min-height: 0;
  max-height: 100%;
  min-width: 240px;
}
.scrollable {
  max-height: 100%;
  overflow: auto;
}
</style>
