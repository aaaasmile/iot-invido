import Navbar from './components/Navbar.js'
import store from './store/index.js'
import routes from './routes.js'
import API from './apicaller.js'



export const app = new Vue({
	el: '#app',
	router: new VueRouter({ routes }),
	components: { Navbar },
	vuetify: new Vuetify(),
	store,
	data() {
		return {
			Buildnr: "",
			links: routes,
			AppTitle: "Iot Invido",
			drawer: false,
			connection: null,
		}
	},
	computed: {
		...Vuex.mapState({

		})
	},
	created() {
		// keep in mind that all that is comming from index.html is a string. Boolean or numerics need to be parsed.
		this.Buildnr = window.myapp.buildnr
		let port = location.port;
		let prefix = (window.location.protocol.match(/https/) ? 'wss' : 'ws')
		let socketUrl = prefix + "://" + location.hostname + (port ? ':' + port : '') + "/websocket";
		this.connection = new WebSocket(socketUrl)
		console.log("WS socket created")

		this.check_credential()

		this.connection.onmessage = (event) => {
			console.log(event)
			let dataMsg = JSON.parse(event.data)
			if (dataMsg.type === "status") {
				console.log('Socket msg type: status')
				this.$store.commit('playerstate', dataMsg)
			} else {
				console.warn('Socket message type not recognized ', dataMsg.type)
			}
		}

		this.connection.onopen = (event) => {
			console.log(event)
			console.log("Socket connection success")
		}

		this.connection.onclose = (event) => {
			console.log(event)
			console.log("Socket closed")
			this.connection = null
		}
	},
	methods: {
		check_credential(){
			console.log('Check credential')
      const tk = localStorage.getItem('tkcred')
			const req = {token: tk}
			API.CheckAPIToken(this, req, (res) => {
        console.log('token validity  check: ', res.Valid)
        if (!res.Valid){
          let path = '/login'
          if (this.$route.path !== path){
            this.$router.replace(path)
          } 
        }
      })
		},
	},
	template: `
  <v-app class="grey lighten-4">
    <Navbar />
    <v-content class="mx-4 mb-4">
      <router-view></router-view>
    </v-content>
    <v-footer>
      <div class="caption">
        {{ new Date().getFullYear() }} —
        <span>Buildnr: {{Buildnr}}</span>
      </div>
    </v-footer>
  </v-app>
`
})

console.log('Main is here!')