export default function DefaultRedirect (params) {
  return {
    beforeRouteEnter (to, from, next) {
      to.params.page ? next() : next({ name: to.name, params: { ...to.params, ...params } })
    },
    beforeRouteUpdate (to, from, next) {
      // to.params.page ? next() : next({ name: to.name, params })
      to.params.page ? next() : next({ name: to.name, params: { ...to.params, ...params } })
    }
  }
}
