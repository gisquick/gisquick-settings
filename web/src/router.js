import VueRouter from 'vue-router'
import FilesView from '@/views/FilesView'
import LayersView from '@/views/LayersView'
import TopicsView from '@/views/TopicsView'
import ProjectView from '@/views/ProjectView'
import ProjectsView from '@/views/ProjectsView'
import PublishWizardView from '@/views/PublishWizardView'
import QgisProject from '@/views/publish/QgisProject'
import Upload from '@/views/publish/Upload'
import Config from '@/views/publish/Config'

// import LayersStore from '@/LayersStore'

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
          path: 'config/:page?',
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
      path: '/:user/:folder/:projectName/',
      name: 'project',
      redirect: { name: 'layers' },
      props: true,
      component: ProjectView,
      children: [
        {
          path: '/',
          redirect: 'layers'
        },
        {
          path: 'layers',
          name: 'layers',
          component: LayersView
          // components: {
          //   default: LayersView,
          //   store: LayersStore
          // }
        },
        {
          path: 'topics',
          name: 'topics',
          component: TopicsView
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
