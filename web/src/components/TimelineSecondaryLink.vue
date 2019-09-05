<template>
  <v-timeline-item
    class="pb-0"
    :color="currentColor"
    small
  >
    <v-btn
      :to="to"
      :color="currentColor"
      rounded
      text
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
    color: String
  },
  computed: {
    linkPath () {
      return this.$router.resolve(this.to).href
    },
    isActive () {
      return this.$route.fullPath === this.linkPath
    },
    currentColor () {
      return this.isActive ? this.color : 'grey'
    }
  }
}
</script>

<style lang="scss" scoped>
.v-timeline-item {
  height: 2.5em;
  ::v-deep .v-timeline-item__body {
    display: flex;
    align-items: center;
  }
}
.v-btn {
  &:before {
    display: none;
  }
  text-transform: none;
}
</style>
