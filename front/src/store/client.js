export default {
  state: {
    clientID: null
  },
  getters: {
    getClientID (state) {
      return state.clientID
    }
  },
  mutations: {
    updateClientID (state, id) {
      state.clientID = id
    }
  },
  actions: {
    async getClient ({ dispatch }, clientID='') {
      return dispatch('fetch', {
        path: '/users/'+clientID+'/config',
        method: 'GET'
      })
    },
    async saveClient ({ dispatch }, data) {
      return dispatch('fetch', {
        path: '/users/'+data.clientID+'/config/'+data.certificateID,
        method: 'PUT',
        body: data
      })
    },
    async createClient ({ dispatch }, clientID) {
      return dispatch('fetch', {
        path: '/users/'+clientID+'/config',
        method: 'POST'
      })
    },
    async deleteClient ({ dispatch }, data) {
      return dispatch('fetch', {
        path: '/users/'+data.clientID+'/config/'+data.certificateID,
        method: 'DELETE'
      })
    },
    async revokeCert ({ dispatch }, data) {
      return dispatch('fetch', {
        path: '/certs/'+data.id,
        method: 'DELETE'
      })
    },
    async downloadConfig ({ dispatch }, data) {
      return dispatch('download', {
        path: '/users/'+data.clientID+'/config/'+data.certificateID,
        method: 'GET'
      })
    },
    async getUsers ({ dispatch }) {
      return dispatch('fetch', {
        path: '/users?all=true',
        method: 'GET'
      })
    },
    async getHosts ({ dispatch }) {
      return dispatch('fetch', {
        path: '/hosts',
        method: 'GET'
      })
    },
    async getCerts ({ dispatch }, username=null) {
      let path = '/certs'
      if (username !== null) {
        path += '?subject=' + username
      }

      return dispatch('fetch', {
        path: path,
        method: 'GET'
      })
    }
  }
}
