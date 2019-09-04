import VueRouter from 'vue-router'
import FilesView from '@/views/FilesView'
import LayersView from '@/views/LayersView'
import TopicsView from '@/views/TopicsView'
import ProjectView from '@/views/ProjectView'
import ProjectsView from '@/views/ProjectsView'
import PublishView from '@/views/PublishView'
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
      component: PublishView,
      redirect: { name: 'qgis-project' },
      children: [
        {
          path: 'qgis/:page?',
          name: 'qgis-project',
          component: QgisProject
        },
        {
          path: 'upload',
          name: 'publish-upload',
          component: Upload
        },
        {
          path: 'config/:page?',
          name: 'publish-config',
          component: Config
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
