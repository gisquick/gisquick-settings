<template>
  <div class="elevation-2">
    <v-data-table
      :headers="headers"
      :items="projects"
      :loading="loading"
      disable-pagination
      hide-default-footer
    >
      <template v-slot:item.title="{ item, value }">
        <router-link :to="{path: item.project}">{{ value }}</router-link>
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
  </div>
</template>

<script>
const Headers = [
  {
    text: 'Project',
    align: 'left',
    sortable: true,
    value: 'title'
  }, {
    text: 'Map',
    align: 'left',
    sortable: true,
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
  activated () {
    this.fetchProjects()
  },
  methods: {
    fetchProjects () {
      this.loading = true
      this.$http.get('/projects.json')
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
</style>
