import API from '../apicaller.js'

export default {
  data() {
    return {
      loading: false,
    }
  },
  computed: {
    ...Vuex.mapState({

    })
  },
  methods: {
    fetchData() {
      console.log('fetchData ')
      API.FetchData(this, req)
    }
  },
  template: `
  <v-card color="grey lighten-4" flat tile>
    <v-toolbar flat dense>
      <v-toolbar-title class="subheading grey--text"
        >Iot Invido</v-toolbar-title
      >
      <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-btn icon @click="syncRepo" :loading="fetchData" v-on="on">
              <v-icon>mdi-sync</v-icon>
            </v-btn>
          </template>
          <span>Update data</span>
        </v-tooltip>
      <v-spacer></v-spacer>
    </v-toolbar>
  </v-card>`
}