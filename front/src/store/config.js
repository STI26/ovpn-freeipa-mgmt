export default {
  actions: {
    async getServerConfig ({ dispatch }) {
      return dispatch('fetch', {
        path: '/config',
        method: 'GET'
      })
    },
    async createServer ({ dispatch }, data) {
      return dispatch('fetch', {
        path: '/config',
        method: 'POST',
        body: data
      })
    },
    async updateCert ({ dispatch }, id) {
      return dispatch('fetch', {
        path: '/config/cert/' + id,
        method: 'POST'
      })
    },
    async updateCA ({ dispatch }) {
      return dispatch('fetch', {
        path: '/config/ca',
        method: 'POST'
      })
    },
    async updateCrl ({ dispatch }) {
      return dispatch('fetch', {
        path: '/config/crl',
        method: 'POST'
      })
    }
  }
}
