<template>
  <v-timeline-item
    large
    class="py-0 my-2"
    :color="color"
    :icon="icon"
  >
    <v-btn
      :to="to"
      :color="isActive ? color : ''"
      :dark="isActive && !disabled"
      active-class=""
      rounded class="step"
      :disabled="disabled"
      @click="$emit('click', $event)"
    >
      <slot>{{ label }}</slot>
    </v-btn>
  </v-timeline-item>
</template>

<script>
export default {
  props: {
    label: String,
    to: Object,
    icon: String,
    color: String,
    disabled: Boolean
  },
  computed: {
    linkPath () {
      return this.to && this.$router.resolve(this.to).href
    },
    isActive () {
      return this.$route.fullPath === this.linkPath || this.$route.fullPath.startsWith(this.linkPath + '/')
    }
  }
}
</script>

<style lang="scss" scoped>
.v-timeline-item {

  .v-btn {
    width: 100%;
    height: 2.5em;
    margin-top: 0.28em;
    text-transform: none;
    font-size: 17px;
    justify-content: flex-start;
    // border-top-right-radius: 6px;
    // border-bottom-right-radius: 6px;
    // border-bottom-left-radius: 6px;
    &:after {
      content: "";
      position: absolute;
      width: 2em;
      border-bottom: 2px dotted #ccc;
      left: -2em;
      box-shadow: none;
    }
  }
}
</style>
