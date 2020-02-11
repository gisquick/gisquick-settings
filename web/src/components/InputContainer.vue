<template>
  <div
    :class="{'v-input--is-focused': focus}"
    @click="focus = true"
  >
    <slot :focus="focus"/>
  </div>
</template>

<script>
export default {
  data () {
    return {
      focus: false
    }
  },
  methods: {
    _clickOutsideTest (e) {
      if (!this.$el.contains(e.target)) {
        this.focus = false
      }
    },
    onKeyDown (e) {
      this.$emit('keydown', e)
    }
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
