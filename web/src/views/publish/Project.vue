<template>
  <div>
    <div v-if="$ws.pluginConnected">
      Plugin: {{ info }}
    </div>
    <div v-else>
      Plugin not connected!
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      info: null
    }
  },
  activated () {
    this.$ws.request('ProjectInfo')
      .then(resp => {
        this.info = resp
      })
  }
}
</script>

<style lang="scss" scoped>

</style>
