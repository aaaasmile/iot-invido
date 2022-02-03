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
        return state.sen.last.sensorid
      },
      SensorPlace: state => {
        return state.sen.last.place
      },
      Measures: state => {
        return state.sen.measures
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
    <v-toolbar>
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
            <p class="ml-3">Sensor: {{ SensorName }}</p>
            <p class="ml-3">Place: {{ SensorPlace }}</p>
          </v-card-title>
          <div class="font-weight-bold ml-8 mb-2">Timeline</div>

          <v-timeline align-top dense>
            <v-timeline-item
              v-for="measure in Measures"
              :key="measure.time"
              :color="measure.color"
              small
            >
              <div>
                <div class="font-weight-normal">
                  IAQ <strong>{{ measure.iaq }}</strong>, acc {{ measure.iaqacc }}  @{{ measure.timeStamp }}
                </div>
                <div>Temperature: {{ measure.temp }}</div>
                <div>Humidiy: {{ measure.humy }}</div>
                <div>Pressure: {{ measure.press }}</div>
                <div>CO2: {{ measure.co2 }}</div>
                <div><strong>{{ measure.iaqclass }}</strong></div>
              </div>
            </v-timeline-item>
          </v-timeline>
        </v-card>
      </v-row>
    </v-container>
  </v-card>`
}