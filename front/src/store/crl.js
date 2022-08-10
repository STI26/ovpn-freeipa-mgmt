export default {
  actions: {
    async getCrl ({ dispatch }) {
      return dispatch('fetch', {
        path: '/crl',
        method: 'GET'
      })
    }
  }
}
