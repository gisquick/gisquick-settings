<template>
  <div>
    <v-data-table
      class="grow mx-1 my-2 elevation-2"
      :headers="headers"
      :items="projects"
      disable-pagination
      hide-default-footer
    >
      <template v-slot:item.title="{ item, value }">
        <router-link :to="{path: item.project}">{{ value }}</router-link>
      </template>
      <template v-slot:item.url="{ value }">
        <a :href="value">
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
    this.$http.get('/projects.json').then(resp => {
      this.$root.projects = resp.data.projects
    })
  }
}
</script>

<style lang="scss" scoped>
a {
  .v-icon {
    color: inherit;
    opacity: 0.8;
  }
}
</style>
