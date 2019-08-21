<template>
  <v-snackbar
    bottom
    :value="open"
    v-bind="params"
    @input="close"
  >
    <template v-if="params">
      <span v-text="params.text"/>
      <v-btn
        v-if="params.action"
        @click="params.action.handler"
        color="yellow"
        text
      >
        {{ params.action.text }}
      </v-btn>
      <v-btn
        v-else
        @click="close"
        :color="contrastColor"
        text
      >
        Close
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script>
export default {
  data () {
    return {
      params: null,
      action: null
    }
  },
  computed: {
    open () {
      return Boolean(this.params)
    },
    contrastColor () {
      return (this.params && this.params.contrastColor) || 'yellow'
    }
  },
  methods: {
    close () {
      this.params = null
      this.action = null
    },
    _open (params) {
      if (this.open) {
        if (this.params.text === params.text) {
          // ignore same notifications
          return
        }
        this._close()
        setTimeout(() => {
          this._open(params)
        }, 300)
      } else {
        this.params = params
      }
    },
    show (text, opts = {}) {
      this._open({ ...opts, text })
    },
    warn (text, opts = {}) {
      this._open({ ...opts, text, color: 'orange', contrastColor: 'black' })
    },
    error (text, opts = {}) {
      this._open({ ...opts, text, color: 'red darken-4' })
    }
  }
}
</script>
