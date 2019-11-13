<template functional>
  <v-layout
    column align-center justify-center
    :class="[data.class, data.staticClass]"
  >
    <h4 class="my-4">QGIS plugin is not connected!</h4>
    <img src="@/assets/qgis_disconnected.svg" class="my-2">
    <!-- <small class="my-4">* A special version of plugin is required, which is not available yet.</small> -->
    <small class="mt-5 mb-2">You can download required plugin here (64-bit architecture):</small>
    <v-layout shrink>
      <v-btn href="/media/files/gisquick2_lin64.zip" text>
        <icon name="linux"/>
      </v-btn>
      <v-btn href="/media/files/gisquick2_win64.zip" text>
        <icon name="windows"/>
      </v-btn>
      <v-btn href="/media/files/gisquick2_darwin64.zip" text>
        <icon name="apple"/>
      </v-btn>
    </v-layout>
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
