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
    async getClients (context) {
      const data = [
        {id: 0, name: 'client 0', numberOfCertificates: 1},
        {id: 1, name: 'client 1', numberOfCertificates: 0},
        {id: 2, name: 'client 2', numberOfCertificates: 2},
        {id: 3, name: 'eclient 3', numberOfCertificates: 1},
        {id: 4, name: 'client 4', numberOfCertificates: 3},
        {id: 11, name: 'client 11', numberOfCertificates: 1},
        {id: 12, name: 'j.dou', numberOfCertificates: 1}
      ]
      return data
    },
    async getClient (context, id) {
      let certificates
      if (id !== 3) {
        certificates = [
          {id: 12+id, revoked: false, issuedOn: 'Tue Jan 11 00:04:26 2022 UTC', expiresOn: 'Fri Jan 12 00:04:26 2024 UTC'},
          {id: 14+id, revoked: true, issuedOn: 'Tue Jan 24 10:00:30 2020 UTC', expiresOn: 'Fri Feb 01 09:18:59 2030 UTC'}
        ]
      } else {
        certificates = []
      }
      const data = {
        id: id,
        subject: `client ${id}`,
        ip: '10.10.2.10',
        routes: 'push 192.168.0.0 dev\npush 192.168.1.0 dev',
        certificates
      }
      return data
    },
    async saveClient (context, data) {
      return 'ok'
    },
    async createClient (context, data) {
      return 'ok'
    },
    async deleteClient (context, data) {
      return 'ok'
    },
    async revokeCert (context, data) {
      return 'ok'
    },
    async downloadConfig (context, data) {
      return 'ok'
    },
    async getUsers (context) {
      const data = [
        {id: 0, name: 'user 0'},
        {id: 1, name: 'user 1'},
        {id: 2, name: 'user 2'},
        {id: 3, name: 'euser 3'},
        {id: 4, name: 'user 4'},
        {id: 11, name: 'user 11'},
        {id: 12, name: 'j.dou'}
      ]
      return data
    },
    async getHosts (context) {
      const data = [
        {id: 0, name: 'host 0'},
        {id: 1, name: 'host 1'},
        {id: 2, name: 'host 2'},
        {id: 3, name: 'ehost 3'},
        {id: 4, name: 'host 4'},
        {id: 11, name: 'host 11'},
        {id: 12, name: 'localhost'}
      ]
      return data
    }
  }
}