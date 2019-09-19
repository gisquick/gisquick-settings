import Vue from 'vue'
import 'material-icons/iconfont/material-icons.scss'
import Vuetify from 'vuetify/lib'

import VIcon from '@/components/VIcon.vue'

Vue.use(Vuetify)

function customIcon (name) {
  return {
    component: VIcon,
    props: { name }
  }
}
const customIcons = ['db', 'qgis']

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
    values: Object.assign({}, ...customIcons.map(name => ({ [name]: customIcon(name) })))
  }
})
