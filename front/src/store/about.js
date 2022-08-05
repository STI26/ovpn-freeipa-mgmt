export default {
  actions: {
    async getApiVerion ({ dispatch }) {
      return dispatch('fetch', {
        path: '/version',
        method: 'GET'
      })
    },
  }
}
