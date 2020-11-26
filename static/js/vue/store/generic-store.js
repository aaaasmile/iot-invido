export default {
  state: {
    errorText: '',
    msgText: '',
    lastmsgText: '',
  },
  mutations: {
    errorText(state, msg) {
      state.errorText = msg
      state.lastmsgText = msg
    },
    msgText(state, msg) {
      state.msgText = msg
      state.lastmsgText = msg
    },
    clearErrorText(state) {
      if (state.errorText !== '') {
        state.errorText = ''
      }
    },
    clearMsgText(state) {
      if (state.msgText !== '') {
        state.msgText = ''
      }
    }
  }
}