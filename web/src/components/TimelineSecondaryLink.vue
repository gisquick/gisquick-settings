<template>
  <v-timeline-item
    class="pb-0 my-3"
    :color="dotColor"
    small
  >
    <router-link
      :to="to"
      :active-class="activeClass"
      class="link"
      tag="span"
    >
      <slot>{{ label }}</slot>
    </router-link>
  </v-timeline-item>
</template>

<script>
export default {
  props: {
    label: String,
    to: Object,
    color: String,
    activeClass: String
  },
  computed: {
    linkPath () {
      return this.$router.resolve(this.to).href
    },
    isActive () {
      return this.$route.fullPath === this.linkPath
    },
    dotColor () {
      return this.isActive ? this.color : 'grey'
    }
  }
}
</script>

<style lang="scss" scoped>
.link {
  margin-left: 20px;
  color: #555;
  cursor: pointer;
}
</style>
