<template>
  <v-layout column align-center justify-center>
    <h4 class="my-4">QGIS plugin is not connected!</h4>
    <img src="@/assets/qgis_disconnected.svg" class="my-2">
    <small class="mt-5 mb-2">You can download required plugin here (64-bit OS):</small>
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
    <small class="orange--text mt-3">
      <v-icon small color="orange" class="mr-1 mb-1">warning</v-icon>
      <span>After changes in API, previous version will not work (2.0.0). Update to the latest version.</span>
    </small>
  </v-layout>
</template>

<script>
export default {

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
  text-align: center;
}
img {
  width: 600px;
  max-width: 80%;
}
</style>
