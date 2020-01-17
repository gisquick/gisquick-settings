import Vue from 'vue'
import VueRouter from 'vue-router'
import PortalVue from 'portal-vue'
import { ReactiveRefs } from 'vue-reactive-refs'

import http from './http.js'
import App from './App.vue'
import ServerError from './ServerError'
import Icon from './components/Icon.vue'
import Expander from './components/Expander.vue'
import vuetify from './plugins/vuetify'
import router from './router.js'
import * as filters from './filters.js'

Vue.config.productionTip = false
Vue.use(VueRouter)
Vue.use(PortalVue)
Vue.use(ReactiveRefs)
Vue.component('icon', Icon)
Vue.component('expander', Expander)

Vue.filter('filesize', filters.formatFileSize)
Vue.prototype.$http = http

function createApp (data) {
  new Vue({
    router,
    vuetify,
    data () {
      return {
        ...data,
        projects: []
      }
    },
    render: h => h(App)
  }).$mount('#app')
}

http.get('/api/app/')
  .then(resp => {
    const data = resp.data
    if (data.user && data.user.is_guest) {
      data.user = null
    }
    createApp(data)
  })
  .catch(() => {
    new Vue({
      vuetify,
      render: h => h(ServerError)
    }).$mount('#app')
  })
