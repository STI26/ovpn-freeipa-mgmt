export default {
  state: {
    token: null,
    userid: null,
    username: null,
  },
  getters: {
    token (state) {
      return state.token
    },
    username (state) {
      return state.username
    },
    ifAuthenticated (state) {
      return state.token !== null
    }
  },
  mutations: {
    auth (state, userData) {
      state.token = userData.token
      state.username = userData.username
      // Set sessionStorage
      sessionStorage.setItem('token', state.token)
      sessionStorage.setItem('username', state.username)
    },
    clearAuth (state) {
      state.token = null
      state.username = null
      // Clear sessionStorage
      sessionStorage.removeItem('token')
      sessionStorage.removeItem('username')
    }
  },
  actions: {
    async login ({ getters }, creds) {
      const url = getters.backendURL+'/login'

      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(creds)
      })

      if (!response.ok) {
        throw `${url}: ${response.status} (${response.statusText})`
      }

      return response.json()
    },
    autoLogin ({ commit }) {
      const token = sessionStorage.getItem('token')
      if (!token) {
        return undefined
      }
      commit('auth', {
        token: sessionStorage.getItem('token'),
        username: sessionStorage.getItem('username')
      })
    }
  }
}