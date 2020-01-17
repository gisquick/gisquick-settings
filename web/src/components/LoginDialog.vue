<template>
  <v-dialog
    fullscreen scrollable
    v-model="open"
    content-class="login-dialog"
  >
    <div class="bg-logo-container">
      <img src="../assets/image_logo.svg">
    </div>
    <v-container>
      <v-card dark color="transparent">
        <v-layout column class="header">
          <img class="logo" src="../assets/text_logo_dark.svg">
          <h3>Sign In to Continue</h3>
        </v-layout>

        <v-card-text>
          <p
            v-if="authenticationError"
            class="text-xs-center red--text">{{ authenticationError }}
          </p>
          <v-form ref="form" v-model="valid" :lazy-validation="true">
            <v-text-field
              label="Username"
              placeholder=" "
              v-model="username"
              @keyup.enter="login"
              color="primary"
            />
            <v-text-field
              label="Password"
              placeholder=" "
              v-model="password"
              :append-icon="passwordVisible ? 'visibility' : 'visibility_off'"
              :type="passwordVisible ? 'text' : 'password'"
              @click:append="passwordVisible = !passwordVisible"
              @keyup.enter="login"
              color="primary"
            />
            <v-btn light color="grey lighten-3" @click="login">
              Login
            </v-btn>
          </v-form>
        </v-card-text>

      </v-card>
    </v-container>
  </v-dialog>
</template>

<script>

export default {
  name: 'LoginDialog',
  data () {
    return {
      open: true,
      authenticationError: '',
      valid: true,
      username: '',
      password: '',
      passwordVisible: false
    }
  },
  methods: {
    login () {
      const form = new FormData()
      form.append('username', this.username)
      form.append('password', this.password)
      this.$http.post('/api/auth/login/', form)
        .then((resp) => {
          this.authenticationError = null
          this.$emit('login', resp.data)
          this.open = false
        })
        .catch(resp => {
          this.authenticationError = 'Authentication Failed'
        })
    }
  }
}
</script>

<style lang="scss">
.login-dialog {
  background-color: black;
  align-items: center;
  .v-card {
    max-width: 400px;

    .header {
      justify-content: center;
      align-items: center;
      .logo {
        width: 80%;
        margin-bottom: 1em;
      }
      h3, h4 {
        margin: 0.5em 0;
        font-weight: 400;
      }
      h4 {
        opacity: 0.5;
      }
    }

    form {
      border: 2px dashed #444;
      padding: 0.5em 1em;
      display: flex;
      flex-direction: column;
      .input-group label:after {
        /* Remove asterics from required fields */
        display: none;
      }
      .btn {
        margin: 0.25em 0;
      }
      input:-webkit-autofill {
        box-shadow: 0 0 0px 1000px #000 inset;
        -webkit-text-fill-color: #fff
      }
    }
  }

  .bg-logo-container {
    position: absolute;
    overflow: hidden;
    top: 0;
    bottom: 0;
    left: calc(400px + (100% - 400px)/ 2);
    width: calc((100% - 400px)/ 2);
    img {
      position: absolute;
      height: 100%;
      left: 40%;
      top: 0;
    }
  }
}
</style>
