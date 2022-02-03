import Generic from './generic-store.js'
import SensorStore from './sensor-store.js'
import UserStore from './user-store.js'

export default new Vuex.Store({
  modules: {
    gen: Generic,
    sen: SensorStore,
    user: UserStore,
  }
})
