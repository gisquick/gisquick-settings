<template>
  <v-layout class="scripts-settings column">
    <!-- <div class="list py-2">
      <div
        v-for="(item, mod) in scripts"
        :key="item.path"
        class="script-info mx-4"
      >
        <div>
          <span class="title" v-text="mod"/>
          (<a target="_blank" :href="`/api/project/static/${namespace}/${item.path}`">{{ item.path }}</a>)
        </div>
        <div class="components">Components: {{ item.components.join(', ') }}</div>
        <v-btn icon @click="deleteScript(mod)">
          <v-icon>delete</v-icon>
        </v-btn>
      </div>
    </div> -->

    <div class="list py-2">
      <div
        v-for="(item, mod) in scripts"
        :key="item.path"
        class="script-info2 mx-4 mb-2"
      >
        <span class="title" v-text="mod"/>
        <div class="link">
          <span>File: </span>
          <a target="_blank" :href="`/api/project/static/${namespace}/${item.path}`">{{ item.path }}</a>
        </div>
        <div class="components">Components: {{ item.components.join(', ') }}</div>
        <v-btn icon @click="deleteScript(mod)">
          <v-icon>delete</v-icon>
        </v-btn>
      </div>
    </div>
    <p v-if="noScriptsDefined" class="px-4 text-center">
      There are no JavaScript modules defined in the project
    </p>

    <v-divider/>
    <v-btn
      v-if="!uploadMode"
      color="primary"
      class="align-self-center my-4"
      @click="uploadMode = true">
      <v-icon class="mr-2">add</v-icon>
      <span>Upload module</span>
    </v-btn>
    <template v-else>
      <v-layout class="mx-4 mt-2">
        <v-icon color="orange" size="18" class="mx-1">warning</v-icon>
        <small class="orange--text">Select only trusted JavaScript files!</small>
      </v-layout>
      <div class="upload mx-4 mb-2">
        <v-file-input
          ref="fileInput"
          show-size
          accept=".js"
          prepend-icon="note_add"
          placeholder="Select file"
          :messages="messages"
          :error-messages="errorMessages"
          @change="loadScript"
        />
        <v-btn
          rounded
          class="mt-2"
          :disabled="!scriptValid || uploading"
          @click="upload"
        >
          <v-icon class="mr-2">cloud_upload</v-icon>
          <span>Upload</span>
        </v-btn>
      </div>
    </template>
  </v-layout>
</template>

<script>

export async function externalComponent (url, filename) {
  return new Promise((resolve, reject) => {
    const name = filename.match(/^(.*?)\.umd\.(.+\.)?js$/)?.[1]
    if (!name) {
      reject('Invalid JavaScript module, expected UMD module')
    }
    const script = document.createElement('script')
    script.async = true
    script.addEventListener('load', () => {
      resolve({
        element: script,
        module: window[name]
      })
    })
    script.addEventListener('error', () => {
      const err = new Error(`Error loading ${url}`)
      reject(err)
    })
    script.src = url
    document.head.appendChild(script)
  })
}

export default {
  props: {
    namespace: String,
    scripts: Object
  },
  data () {
    return {
      error: null,
      uploadError: null,
      uploadMode: false,
      uploadScript: null,
      uploading: false
    }
  },
  computed: {
    noScriptsDefined () {
      return !this.scripts || Object.keys(this.scripts).length === 0
    },
    validationError () {
      if (this.error) {
        return this.error
      }
      if (this.uploadScript) {
        const components = this.uploadScript.components
        if (!Array.isArray(components) || components.length < 1) {
          return 'No components found'
        }
        if (components.some(n => !n)) {
          return "Some components doesn't have defined name"
        }
        // TODO: check duplicate components names in other modules
      }
      return null
    },
    errorMessages () {
      const err = this.validationError || this.uploadError
      return err ? [err] : []
    },
    messages () {
      if (this.uploadScript) {
        return [`Components: ${this.uploadScript.components?.join(', ')}`]
      }
      return []
    },
    scriptValid () {
      return this.uploadScript && !this.validationError
    }
  },
  methods: {
    async loadScript (f) {
      if (!f) {
        if (this.script) {
          URL.revokeObjectURL(this.script.element.src)
          this.script.element.remove()
          this.script = null
        }
        this.uploadScript = null
        this.error = null
        return
      }
      const objectURL = URL.createObjectURL(f)
      try {
        const script = await externalComponent(objectURL, f.name)
        let components
        if (script.module.__esModule) {
          components = script.module.default.map(c => c.name)
        } else {
          components = [script.module.name]
        }
        this.script = script
        this.uploadScript = Object.freeze({
          file: f,
          components
        })
        this.error = null
      } catch (err) {
        this.error = err
      }
    },
    async upload () {
      const { file, components } = this.uploadScript
      const form = new FormData()
      const path = 'components/' + file.name
      const info = { path, components: components.filter(n => !!n) }
      form.append('info', JSON.stringify(info))
      form.append('script', file, path)
      const { user, folder } = this.$route.params
      this.uploading = true
      this.uploadError = false
      try {
        const { data } = await this.$http.post(`/api/project/script/${user}/${folder}`, form)
        this.$emit('update:scripts', data)
        this.uploadScript = null
        this.uploadMode = false
        this.$refs.fileInput.reset()
      } catch (err) {
        this.uploadError = 'Failed to upload file'
      } finally {
        this.uploading = false
      }
    },
    async deleteScript (name) {
      const { user, folder } = this.$route.params
      const { data } = await this.$http.delete(`/api/project/script/${user}/${folder}/${name}`)
      this.$emit('update:scripts', data)
    }
  }
}
</script>

<style lang="scss" scoped>
.scripts-settings {
  min-height: 0;
  max-height: 100%;
  overflow: hidden;
  .upload {
    display: grid;
    grid-template-columns: 1fr auto;
    align-items: start;
    gap: 8px;
  }
  .script-info {
    display: grid;
    grid-template-columns: 1fr auto;
    &:not(:first-child) {
      border-top: 1px solid #eee;
    }
    .components {
      grid-area: 2 / 1 / 3 / 2;
    }
    .v-btn {
      grid-area: 1 / 2 / 2 / 3;
    }
  }
  .script-info2 {
    display: grid;
    grid-template-columns: 1fr auto;
    &:not(:first-child) {
      border-top: 1px solid #eee;
    }
    .link {
      grid-area: 2 / 1 / 3 / 2;
    }
    .components {
      grid-area: 3 / 1 / 4 / 2;
    }
    .v-btn {
      grid-area: 1 / 2 / 3 / 3;
      align-self: center;
    }
  }
  .list {
    display: flex;
    flex-direction: column;
    overflow: auto;
    min-height: 0;
    flex-shrink: 1;
  }
}
</style>
