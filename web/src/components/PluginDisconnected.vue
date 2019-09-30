<template functional>
  <v-layout
    column align-center justify-center
    :class="[data.class, data.staticClass]"
  >
    <h4 class="my-4">QGIS plugin is not connected!</h4>
    <!-- <small class="mb-4">But don't try to find it, it's not public yet.</small> -->
    <img src="@/assets/qgis_disconnected.svg" class="my-2">
    <small class="my-4">* A special version of plugin is required, which is not available yet.</small>
  </v-layout>
</template>

<script>
export default {
  functional: true
}
</script>

<style lang="scss" scoped>
h4 {
  color: #F44336;
  font-weight: 500;
}
small {
  color: #777;
  font-weight: 500;
}
img {
  width: 600px;
  max-width: 80%;
}
</style>
