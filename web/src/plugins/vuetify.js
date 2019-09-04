import Vue from 'vue'
import 'material-icons/iconfont/material-icons.scss'
import Vuetify from 'vuetify/lib'

import VIcon from '@/components/VIcon.vue'

Vue.use(Vuetify)

export default new Vuetify({
  theme: {
    themes: {
      light: {
        // primary: '#00b8e5'
        primary: '#00a6e5'
      }
    }
  },
  icons: {
    values: {
      qgis: {
        component: VIcon,
        props: {
          name: 'qgis'
        }
      }
    }
  },
})
