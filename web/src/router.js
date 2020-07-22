import VueRouter from 'vue-router'
import FilesView from '@/views/FilesView'
import ProjectView from '@/views/ProjectView'
import ProjectsView from '@/views/ProjectsView'
import PublishWizardView from '@/views/PublishWizardView'
import QgisProject from '@/views/publish/QgisProject'
import Upload from '@/views/publish/Upload'
import Config from '@/views/publish/Config'

export default new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'projects',
      component: ProjectsView
    },
    {
      path: '/publish',
      name: 'publish',
      component: PublishWizardView,
      redirect: { name: 'qgis-project' },
      children: [
        {
          path: 'qgis/:page?',
          name: 'qgis-project',
          component: QgisProject,
          meta: { requiresPlugin: true }
        },
        {
          path: 'upload',
          name: 'publish-upload',
          component: Upload,
          meta: { requiresPlugin: true }
        },
        {
          path: 'config/:page(.*)?',
          name: 'publish-config',
          component: Config
        },
        {
          path: 'final',
          name: 'publish-final'
        }
      ]
    },
    {
      path: '/:user',
      name: 'user-projects',
      component: ProjectsView,
      props: true
    },
    {
      path: '/:user/:folder/:projectName/',
      name: 'project',
      props: true,
      component: ProjectView,
      redirect: { name: 'settings' },
      children: [
        {
          path: 'settings/:page(.*)?',
          name: 'settings',
          component: Config
        },
        {
          path: 'files',
          name: 'files',
          component: FilesView,
          props: true
        }
      ]
    }
  ]
})
