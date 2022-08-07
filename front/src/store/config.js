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
    }
  }
}
