export default function DefaultRedirect (params) {
  return {
    beforeRouteEnter (to, from, next) {
      to.params.page ? next() : next({ name: to.name, params })
    },
    beforeRouteUpdate (to, from, next) {
      to.params.page ? next() : next({ name: to.name, params })
    }
  }
}
