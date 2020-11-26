import API from '../apicaller.js'

export default {
  data() {
    return {
      loading: false,
    }
  },
  computed: {
    ...Vuex.mapState({
      LastMsg: state => {
        return state.gen.lastmsgText
      },
    })
  },
  methods: {
    fetchData() {
      console.log('fetchData ')
      let req = { from: 'lastday' }
      API.FetchData(this, req)
    }
  },
  template: `
  <v-card color="grey lighten-4" flat tile>
    <v-toolbar flat dense>
      <v-toolbar-title class="subheading grey--text"
        >Dashboard Sensor</v-toolbar-title
      >
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="fetchData" :loading="loading" v-on="on">
            <v-icon>mdi-sync</v-icon>
          </v-btn>
        </template>
        <span>Update data</span>
      </v-tooltip>
      <v-spacer></v-spacer>
    </v-toolbar>
    <v-container>
      <v-divider></v-divider>
      <v-row justify="space-around">
        <v-card>
          <v-card-text>Last message</v-card-text>
          <div class="mx-4">{{ LastMsg }}</div>
        </v-card>
      </v-row>
    </v-container>
  </v-card>
`
}