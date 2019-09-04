<template>
  <v-layout column align-center justify-center>
    <h4 class="my-4">QGIS plugin is not connected!</h4>
    <img src="@/assets/qgis_disconnected.svg" class="my-2">
  </v-layout>
</template>

<script>
export default {

}
</script>

<style lang="scss" scoped>
.layout {
  border: 2px dashed #ddd;
}
h4 {
  color: #F44336;
  font-weight: 500;
}
img {
  max-width: 600px;
}
</style>
