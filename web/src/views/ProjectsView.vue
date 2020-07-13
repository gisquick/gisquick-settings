<template>
  <v-layout class="content column elevation-2 my-2 mx-1">
    <portal to="menu-actions" v-if="pageVisible">
      <projects-menu @upload-finished="fetchProjects"/>
    </portal>
    <expander/>
    <v-data-table
      :headers="headers"
      :items="projects"
      :loading="loading"
      disable-pagination
      hide-default-footer
    >
      <template v-slot:item.title="{ item, value }">
        <router-link :to="{ path: '/' + item.project }">{{ value || item.project }}</router-link>
      </template>
      <template v-slot:item.url="{ value }">
        <a :href="value" target="_blank">
          <v-icon>map</v-icon>
        </a>
      </template>
      <template v-slot:item.cache="{ value }">
        <v-icon v-text="checkIcon[value]"/>
      </template>
      <template v-slot:item.ows_url="{ value }">
        <a :href="`${value}&service=WMS&request=GetCapabilities`">
          <v-icon>link</v-icon>
        </a>
      </template>
    </v-data-table>
  </v-layout>
</template>

<script>
import Page from '@/mixins/Page'
import ProjectsMenu from '@/components/ProjectsMenu'

const Headers = [
  {
    text: 'Project',
    align: 'left',
    sortable: true,
    value: 'title'
  }, {
    text: 'Map',
    align: 'left',
    sortable: false,
    value: 'url',
    width: 75
  }, {
    text: 'Authentication',
    align: 'left',
    sortable: false,
    value: 'authentication'
  }, {
    text: 'Publication time',
    align: 'left',
    sortable: false,
    value: 'publication_time'
  }, {
    text: 'Cache',
    align: 'left',
    sortable: false,
    value: 'cache'
  }, {
    text: 'WMS',
    align: 'left',
    sortable: false,
    value: 'ows_url'
  }
]

export default {
  name: 'Projects',
  mixins: [ Page ],
  components: { ProjectsMenu },
  props: {
    user: String
  },
  data () {
    return {
      loading: false
    }
  },
  computed: {
    projects () {
      return this.$root.projects
    },
    headers () {
      return Headers
    },
    checkIcon () {
      return {
        true: 'check',
        false: 'remove'
      }
    },
    // links () {
    //   return Object.assign({}, ...this.projects.map(p => ({ [p.project]: p.project.split('/').slice(1).join('/') })))
    // }
  },
  watch: {
    user: 'fetchProjects'
  },
  activated () {
    this.fetchProjects()
  },
  methods: {
    fetchProjects () {
      this.loading = true
      const url = this.user ? `/api/projects/${this.user}/` : '/api/projects/'
      this.$http.get(url)
        .then(resp => {
          this.$root.projects = resp.data.projects
          this.loading = false
        })
        .catch(err => {
          this.loading = false
          this.$notification.error('Failed to load projects!')
        })
    }
  }
}
</script>

<style lang="scss" scoped>
.v-data-table {
  width: 100%;
}
a {
  .v-icon {
    color: inherit;
    opacity: 0.8;
  }
}
.content {
  grid-row: 2 / 3;
  grid-column: 2 / 3;
  width: 1200px;
  overflow: auto;
  background-color: #fff;
  @media (max-width: 1600px) {
    width: auto;
  }
  @media (max-width: 1450px) {
    grid-column: 1 / 4;
  }
}
</style>
