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
      <template v-slot:item.cache="{ value }">
        <v-icon v-text="checkIcon[value]"/>
      </template>
    </v-data-table>
  </div>
</template>

<script>
const Headers = [
  {
    text: 'Title',
    align: 'left',
    sortable: true,
    value: 'title'
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
