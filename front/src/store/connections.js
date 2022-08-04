export default {
  actions: {
    async getConnections ({ dispatch }) {
      return dispatch('fetch', {
        path: '/status',
        method: 'GET'
      })
    },
  }
}
