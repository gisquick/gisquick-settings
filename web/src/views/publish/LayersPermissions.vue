<template>
  <div class="grid column box">
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
      <v-layout class="shrink my-1">
        <v-btn text @click="addRole">
          <v-icon class="mr-2">add_circle</v-icon>
          <span>Add</span>
        </v-btn>
        <v-btn
          text
          :disabled="!selectedRole"
          @click="removeRole"
        >
          <v-icon class="mr-2">remove_circle</v-icon>
          <span>Remove</span>
        </v-btn>
      </v-layout>
    </v-layout>

    <v-layout
      v-if="selectedRole"
      class="form column my-1"
    >
      <v-text-field
        label="Name"
        class="mx-2 my-2 shrink"
        v-model="selectedRole.name"
        hide-details
      />
      <v-divider/>

      <!-- Users -->
      <v-layout class="column shrink">
        <v-layout class="align-center header px-2 py-0 grey lighten-3">
          <span>Users</span>
          <v-spacer/>
          <v-btn
            icon class="mx-1"
            @click="editUsers = !editUsers"
            :color="editUsers ? 'primary' : ''"
            :ripple="false"
          >
            <v-icon>edit</v-icon>
          </v-btn>
        </v-layout>
        <v-divider/>
        <v-text-field
          v-if="!editUsers"
          class="mx-2 my-2"
          :value="roleUsersText"
          placeholder="No users asigned"
          xappend-icon="edit"
          @click:append="editUsers = !editUsers"
          xdisabled
          readonly
          hide-details
        />
        <switch-lists
          v-if="editUsers"
          class="mx-2 my-2"
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
        class="layers my-1"
      >
        <template v-slot:leaf.view="{ item }">
          <v-checkbox
            v-model="layersPermissions[item.serverName].view"
            color="grey darken-1"
            class="my-0 py-1 justify-center"
            :ripple="false"
            hide-details
          />
        </template>
        <template v-slot:leaf.insert="{ item }">
          <v-checkbox
            v-model="layersPermissions[item.serverName].insert"
            color="grey darken-1"
            class="my-0 py-1 justify-center"
            :ripple="false"
            hide-details
          />
        </template>
        <template v-slot:leaf.update="{ item }">
          <v-checkbox
            v-model="layersPermissions[item.serverName].update"
            color="grey darken-1"
            class="my-0 py-1 justify-center"
            :ripple="false"
            hide-details
          />
        </template>
        <template v-slot:leaf.delete="{ item }">
          <v-checkbox
            v-model="layersPermissions[item.serverName].delete"
            color="grey darken-1"
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
    roleUsersText () {
      return this.selectedRole && this.selectedRole.users.join(', ')
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
    selectedRole (role) {
      if (role) {
        this.initPermissions(role)
      }
    }
  },
  methods: {
    setProjectAccess (enabled) {
      enabled = Boolean(enabled)
      if (this.accessControl) {
        this.accessControl.enabled = enabled
      } else {
        const accessControl = {
          enabled,
          roles: []
        }
        this.$set(this.config, 'access_control', accessControl)
      }
    },
    initPermissions (role) {
      const permissions = {}
      layersList(this.config.overlays).forEach(l => {
        permissions[l.serverName] = {
          view: true,
          insert: false,
          update: false,
          delete: false
        }
      })
      role.permissions.layers = {
        ...permissions,
        ...(role.permissions.layers || {})
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
      this.accessControl.roles.push({
        name: 'New',
        users: [],
        permissions: {
          layers: {}
        }
      })
    },
    removeRole () {
      this.accessControl.roles = this.accessControl.roles.filter(r => r !== this.selectedRole)
      this.selectedRole = null
    }
  }
}
</script>

<style lang="scss" scoped>
.box {
  border: 1px solid #ccc;
  background-color: #fff;
}
.header {
  font-weight: 500;
  font-size: 15px;
  grid-column: 1 / 3;
  border-bottom: 1px solid #ccc;
  background-color: #eee;
}
.grid {
  display: grid;
  grid-template-columns: minmax(20%, auto) 1fr;
  grid-template-rows: auto 1fr;
  overflow: hidden;
  max-height: 100%;
  flex: 1 1;
}
.roles {
  overflow: auto;
  border-right: 1px solid #ccc;
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
