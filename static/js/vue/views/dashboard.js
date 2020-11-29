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
      SensorName: state => {
        return state.sen.sensorid
      },
      SensorPlace: state => {
        return state.sen.place
      },
    })
  },
  methods: {
    fetchData() {
      console.log('fetchData ')
      let req = { from: 'lastday' }
      API.FetchData(this, req)
    },
    insertData() {
      console.log('insertData')
      let req = { type: 'test' }
      API.InsertData(this, req)
    },
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
        <span>Check data</span>
      </v-tooltip>
      <v-spacer></v-spacer>
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="insertData" :loading="loading" v-on="on">
            <v-icon>mdi-chart-arc</v-icon>
          </v-btn>
        </template>
        <span>Insert some data</span>
      </v-tooltip>
    </v-toolbar>
   <v-container>
    <v-row justify="space-around">
      <v-card width="400">
          <v-card-title class="grey--text mt-8">
            <p class="ml-3">
              Sensor: {{ SensorName }}
            </p>
            <p class="ml-3">
              Place: {{ SensorPlace }}
            </p>
          </v-card-title>
      </v-card>
    </v-row>
   </v-container>
  </v-card>
`
}