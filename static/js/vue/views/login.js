import API from '../apicaller.js'

export default {
  data() {
    return {
      loading: false,
      password: ''
    }
  },
  computed: {
    Username: {
      get() {
        return this.$store.state.user.username
      },
      set(newVal) {
        this.$store.commit('setUsername', newVal)
      }
    },
    ...Vuex.mapState({
      LastMsg: state => {
        return state.gen.lastmsgText
      },
      IsLastMessage: state => {
        return state.gen.lastmsgText !== ''
      },
    })
  },
  methods: {
    signIn() {
      console.log('signIn ')
      let req = { username: this.Username, password: this.password }
      API.SignIn(this, req, (res) => {
        console.log('Sign In result is ', res)
      })
    },

  },
  template: `
  <v-row justify="center">
    <v-col class="mb-12" cols="12" md="10">
      <v-card>
        <v-card-title>Login</v-card-title>
        <v-container>
          <v-row justify="space-around">
            <v-col cols="8">
              <v-row justify="space-around">
                <v-text-field
                  label="Username"
                  v-model="Username"
                ></v-text-field>
              </v-row>
              <v-row>
                <v-text-field
                  label="Password"
                  type="password"
                  v-model="password"
                ></v-text-field>
              </v-row>
            </v-col>
          </v-row>
        </v-container>
        <v-card-actions>
          <v-row justify="space-around">
            <v-btn color="primary" text @click="signIn"> Sign In </v-btn>
          </v-row>
        </v-card-actions>
      </v-card>

      <v-row v-if="IsLastMessage">
        <v-container>
          {{ LastMsg }}
        </v-container>
      </v-row>
    </v-col>
  </v-row>
`
}