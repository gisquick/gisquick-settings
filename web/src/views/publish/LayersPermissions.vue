<template>
  <div class="grid">
    <v-layout class="align-center header px-2 py-2">
      <span>Roles</span>
      <v-spacer/>
      <v-switch
        class="my-0 py-0"
        color="primary"
        :input-value="accessControl && accessControl.enabled"
        @change="setProjectAccess"
        :ripple="false"
        hide-details
      />
    </v-layout>

    <v-layout class="roles column grow">
      <input-container
        @keydown.delete="removeRole"
        class="grow layout"
      >
      <v-list dense class="grow py-0">
        <v-list-item
          v-for="(role, index) in roles"
          :key="index"
          :input-value="selectedRole === role"
          @click="selectedRole = role"
          color="primary"
          :ripple="false"
          v-text="role.name"
        />
      </v-list>
      </input-container>
      <v-divider/>
      <div class="toolbar shrink">
        <v-btn
          text
          :disabled="!accessControl"
          @click="addRole"
        >
          <v-icon class="mr-2">add_circle</v-icon>
          <span>Add</span>
        </v-btn>
        <v-divider vertical/>
        <v-btn
          text
          :disabled="!selectedRole"
          @click="removeRole"
        >
          <v-icon class="mr-2">remove_circle</v-icon>
          <span>Remove</span>
        </v-btn>
      </div>
    </v-layout>

    <v-layout
      v-if="selectedRole"
      class="form column my-1"
    >
      <v-row class="mx-1 my-2 shrink no-gutters">
        <v-col>
          <v-text-field
            label="Name"
            class="mx-2"
            v-model="selectedRole.name"
            hide-details
          />
        </v-col>
        <v-col>
          <v-select
            label="Authentication"
            v-model="selectedRole.auth"
            :items="authTypes"
            class="mx-2"
            prepend-icon="security"
            hide-details
          />
        </v-col>
      </v-row>

      <!-- Users -->
      <v-layout
        v-if="selectedRole.auth === 'users'"
        class="column shrink"
      >
        <v-layout class="align-center header px-3 py-0">
          <span>Users</span>
        </v-layout>
        <v-text-field
          class="mx-3 my-2"
          placeholder="No users asigned"
          :value="roleUsersText"
          :append-outer-icon="editUsers ? 'edit_off' : 'edit'"
          @click:append-outer="editUsers = !editUsers"
          readonly
          hide-details
        />
        <switch-lists
          v-if="editUsers"
          class="mx-3 my-2"
          item-text="label"
          item-value="username"
          :items="users"
          v-model="selectedRole.users"
        />

      </v-layout>

      <!-- Layers -->
      <layers-table
        label="Layer"
        :items="config.overlays"
        :headers="overlaysHeaders"
        :opened.sync="opened"
        class="layers"
      >
        <template v-slot:leaf.view="{ item }">
          <v-checkbox
            v-model="layersPermissions[item.serverName || item.name].view"
            color="secondary"
            class="my-0 py-1 justify-center"
            :ripple="false"
            hide-details
          />
        </template>
        <template v-slot:leaf.insert="{ item }">
          <v-checkbox
            v-model="layersPermissions[item.serverName || item.name].insert"
            color="secondary"
            class="my-0 py-1 justify-center"
            :ripple="false"
            hide-details
          />
        </template>
        <template v-slot:leaf.update="{ item }">
          <v-checkbox
            v-model="layersPermissions[item.serverName || item.name].update"
            color="secondary"
            class="my-0 py-1 justify-center"
            :ripple="false"
            hide-details
          />
        </template>
        <template v-slot:leaf.delete="{ item }">
          <v-checkbox
            v-model="layersPermissions[item.serverName || item.name].delete"
            color="secondary"
            class="my-0 py-1 justify-center"
            :ripple="false"
            hide-details
          />
        </template>
      </layers-table>
    </v-layout>
  </div>
</template>

<script>
import LayersTable from '@/components/LayersTable'
import InputContainer from '@/components/InputContainer'
import SwitchLists from '@/components/SwitchLists'
import { layersList } from '@/utils'

export default {
  name: 'LayersPermissions',
  components: { LayersTable, InputContainer, SwitchLists },
  props: {
    config: Object
  },
  data () {
    return {
      opened: [],
      users: [],
      selectedRole: null,
      selectedUser: null,
      editUsers: false,
    }
  },
  computed: {
    overlaysHeaders () {
      return [
        {
          text: 'View',
          value: 'view'
        }, {
          text: 'Insert',
          value: 'insert'
        }, {
          text: 'Update',
          value: 'update'
        }, {
          text: 'Delete',
          value: 'delete'
        }
      ]
    },
    authTypes () {
      return [
        {
          text: 'All',
          value: 'all'
        }, {
          text: 'All authenticated',
          value: 'authenticated'
        }, {
          text: 'Anonymous',
          value: 'anonymous'
        }, {
          text: 'Selected users',
          value: 'users'
        }
      ]
    },
    roleUsersText () {
      return this.selectedRole && (this.selectedRole.users || []).join(', ')
    },
    layersPermissions () {
      return this.selectedRole && this.selectedRole.permissions.layers
    },
    accessControl () {
      return this.config.access_control
    },
    roles () {
      return (this.accessControl && this.accessControl.roles) || []
    }
  },
  created () {
    this.fetchUsers()
  },
  watch: {
    accessControl: {
      immediate: true,
      handler (accessControl) {
        this.selectedRole = null
        if (accessControl && accessControl.roles) {
          accessControl.roles = accessControl.roles.map(role => {
            const layerKeys = layersList(this.config.overlays).map(l => l.serverName || l.name)
            const layerPerms = {}
            layerKeys.forEach(key => {
              layerPerms[key] = role.permissions.layers[key] || {
                view: true,
                insert: false,
                update: false,
                delete: false
              }
            })
            return {
              ...role,
              permissions: { layers: layerPerms }
            }
          })
        }
      }
    }
  },
  methods: {
    createRole (params) {
      const layersPerms = {}
      layersList(this.config.overlays).forEach(l => {
        layersPerms[l.serverName || l.name] = {
          view: true,
          insert: false,
          update: false,
          delete: false
        }
      })
      return {
        name: 'New',
        auth: 'users',
        ...params,
        users: [],
        permissions: {
          layers: layersPerms
        }
      }
    },
    setProjectAccess (enabled) {
      enabled = Boolean(enabled)
      if (this.accessControl) {
        this.accessControl.enabled = enabled
      } else {
        const accessControl = {
          enabled,
          roles: [
            this.createRole({
              name: 'Public',
              auth: this.config.authentication === 'all' ? 'all' : 'authenticated'
            }
          )]
        }
        this.$set(this.config, 'access_control', accessControl)
      }
    },
    setAuthentication (type) {
      this.selectedRole.auth = type
      if (type !== 'users' && this.selectedRole.users.length) {
        this.selectedRole.users = []
      }
    },
    fetchUsers () {
      this.$http.get('/api/users/')
        .then(resp => {
          this.users = resp.data.map(u => ({
            ...u,
            label: u.username === u.full_name ? u.username : `${u.username} (${u.full_name})`
          }))
          // const users = new Array(20).fill(0).map((x, i) => ({
          //   id: 100 + i,
          //   username: `user_${i}`,
          //   full_name: `user_${i}`
          // }))
          // this.users = resp.data.concat(users)
        })
    },
    addRole () {
      this.accessControl.roles.push(this.createRole())
    },
    removeRole () {
      this.accessControl.roles = this.accessControl.roles.filter(r => r !== this.selectedRole)
      this.selectedRole = null
    }
  }
}
</script>

<style lang="scss" scoped>
.header {
  font-weight: 500;
  font-size: 15px;
  grid-column: 1 / 3;
  border: solid #ccc;
  border-width: 1px 0 1px 0;
  background-color: #eee;
  min-height: 40px;
}
.grid {
  display: grid;
  grid-template-columns: minmax(20%, auto) 1fr;
  grid-template-rows: auto 1fr;
  overflow: hidden;
  max-height: 100%;
  flex: 1 1;
  border: solid #ccc;
  border-width: 0 1px 1px 1px;
}
.roles {
  overflow: auto;
  border-right: 1px solid #ccc;
  .toolbar {
    display: grid;
    grid-template-columns: 1fr auto 1fr;
  }
}
.layers {
  overflow: visible;
}
.form {
  overflow: auto;
}
.switch-lists {
  min-height: 180px;
  max-height: 300px;
}
</style>
