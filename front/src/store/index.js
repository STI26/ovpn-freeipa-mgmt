import { createStore } from 'vuex'
// import auth from '@/store/auth'
import client from '@/store/client'

export default createStore({
  state: {
    backendURL: 'process.env.VUE_APP_BASEURL',
    searchString: ''
  },
  getters: {
    backendURL (state) {
      return state.backendURL || 'http://localhost:8000'
    },
    searchString (state) {
      return state.searchString
    }
  },
  mutations: {
    updateSearchString (state, q) {
      state.searchString = q
    }
  },
  actions: {},
  modules: {
    // auth,
    client
  }
})
