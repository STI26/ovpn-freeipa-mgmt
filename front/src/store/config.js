export default {
  actions: {
    async getServerConfig ({ dispatch }) {
      return dispatch('fetch', {
        path: '/config',
        method: 'GET'
      })
    },
  }
}
