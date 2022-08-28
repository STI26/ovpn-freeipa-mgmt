export default {
  state: {
    darkmode: false
  },
  getters: {
    darkmode (state) {
      return state.darkmode
    }
  },
  mutations: {
    setDarkmode (state, enable) {
      state.darkmode = enable
      // Set localStorage
      localStorage.setItem('darkmode', String(enable))
      
      if (enable) {
        document.body.classList.add('dark')
      }
      else {
        document.body.classList.remove('dark')
      }
    }
  },
  actions: {
    loadTheme ({ commit }) {
      const darkmode = localStorage.getItem('darkmode')
      if (!darkmode) {
        commit('setDarkmode', false)
      } else {
        commit('setDarkmode', darkmode === 'true' ? true : false)
      }
    }
  }
}
