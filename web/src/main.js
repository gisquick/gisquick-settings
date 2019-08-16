import Vue from 'vue'
import VueRouter from 'vue-router'
import PortalVue from 'portal-vue'
import http from './http.js'
import App from './App.vue'
import Icon from './components/Icon.vue'
import vuetify from './plugins/vuetify'
import router from './router.js'
import * as filters from './filters.js'

Vue.config.productionTip = false
Vue.use(VueRouter)
Vue.use(PortalVue)
Vue.component('icon', Icon)

Vue.filter('filesize', filters.formatFileSize)
Vue.prototype.$http = http

if (process.env.NODE_ENV === 'development') {
  var initialize = new Promise((resolve, reject) => {
    http.get('/dev/')
      .then(resp => resolve(resp.data))
      .catch(reject)
  })
} else {
  initialize = new Promise(resolve => {
    resolve(JSON.parse(document.getElementById('app-data').textContent))
  })
}

initialize.then(data => {
  new Vue({
    router,
    vuetify,
    data () {
      return {
        user: data.user,
        projects: []
      }
    },
    render: h => h(App)
  }).$mount('#app')
})
