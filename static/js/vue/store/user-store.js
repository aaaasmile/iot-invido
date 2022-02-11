export default {
  state: {
    username: '',
    token: '',
    tokenchecked: false,
  },
  mutations: {
    setUsername(state, newVal) {
      state.username = newVal
    },
    setUserToken(state, newVal) {
      state.token = newVal
      state.tokenchecked = true
      localStorage.setItem('tkcred',newVal)
    }
  }
}