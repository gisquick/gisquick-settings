<template>
  <v-layout class="components-select">
    <v-select
      class="components-list"
      clearable
      outlined
      dense
      :label="label"
      :items="items"
      :value="value"
      :error-messages="errors"
      @input="$emit('input', $event)"
    />

    <v-dialog
      max-width="800"
      v-model="dialog"
      scrollable
    >
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          icon
          class="ml-2 my-2"
          v-bind="attrs"
          v-on="on"
        >
          <v-icon>settings</v-icon>
        </v-btn>
      </template>
      <v-card>
        <v-toolbar dense dark color="grey darken-4">
          <v-toolbar-title>Modules</v-toolbar-title>
          <v-spacer/>
          <v-toolbar-items>
          <v-btn icon @click="dialog = false">
            <v-icon>close</v-icon>
          </v-btn>
          </v-toolbar-items>
        </v-toolbar>
        <v-card-text class="px-0 py-0">
          <scripts-settings
            :namespace="namespace"
            :scripts="scripts"
            @update:scripts="updateScripts"
          />
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-layout>
</template>

<script>
import ScriptsSettings from '@/components/ScriptsSettings.vue'

export default {
  components: { ScriptsSettings },
  props: {
    label: String,
    value: String,
  },
  data () {
    return {
      dialog: false,
      active: false
    }
  },
  computed: {
    namespace () {
      return this.$root.scriptsStore.namespace
    },
    scripts () {
      // lazy load of scripts
      if (this.active) {
        return this.$root.scriptsStore.getScripts()?.data
      }
      return {}
    },
    components () {
      return this.$root.scriptsStore.components
    },
    missingModule () {
      return this.value && !this.$root.scriptsStore.loading && !this.components.includes(this.value)
    },
    errors () {
      if (this.missingModule) {
        return 'Component not found in available modules'
      }
      return []
    },
    items () {
      if (this.missingModule) {
        // add current value to the items
        return [this.value].concat(this.components)
      }
      return this.components
    }
  },
  activated () {
    this.active = true
  },
  deactivated () {
    this.active = false
  },
  methods: {
    updateScripts (data) {
      this.$root.scriptsStore.update(data)
    }
  }
}
</script>

<style lang="scss" scoped>
.v-select.components-list {
  max-width: 500px;
}
::v-deep .v-subheader {
  height: 18px;
  justify-content: center;
  &::before, &::after {
    content: "";
    border-bottom: 1px solid #eee;
    flex-grow: 1;
  }
  &::before {
    margin-right: 8px;
  }
  &::after {
    margin-left: 8px;
  }
}
</style>
