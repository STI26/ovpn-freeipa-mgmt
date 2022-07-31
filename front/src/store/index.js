import { createStore } from 'vuex'
import auth from '@/store/auth'
import client from '@/store/client'
import { Toast, Modal } from 'bootstrap'


export default createStore({
  state: {
    backendURL: import.meta.env.VITE_API_URL,
    searchString: '',
    toast: {obj: null, text: '', color: ''},
    modal: {
      obj: null,
      title: '',
      text: '',
      action: '',
      data: {}
    }
  },
  getters: {
    backendURL (state) {
      return state.backendURL || 'http://localhost:8000'
    },
    searchString (state) {
      return state.searchString
    },
    toast (state) {
      return state.toast
    },
    modal (state) {
      return state.modal
    }
  },
  mutations: {
    updateSearchString (state, q) {
      state.searchString = q
    },
    updateToast (state, data) {
      if (state.toast.obj === null) {
        state.toast.obj = new Toast(document.getElementById('liveToast'))
      }
      state.toast.text = data.text
      state.toast.color = data.color
    },
    updateModal (state, data) {
      if (state.modal.obj === null) {
        state.modal.obj = new Modal(document.getElementById('modalApprovalForm'))
      }
      state.modal.title = data.title
      state.modal.text = data.text
      state.modal.action = data.action
      state.modal.data = data.data
    }
  },
  actions: {
    async fetch ({ getters, commit }, data) {
      const body = {
        method: data.method,
        headers: {
          'Authorization': getters.token,
          'Content-Type': 'application/json'
        }
      }

      if (data.body) {
        body.body = JSON.stringify(data.body)
      }

      const url = getters.backendURL+data.path
      const response = await fetch(url, body)
      
      if (response.status === 401) {
        commit('clearAuth')
      }
      
      if (!response.ok) {
        const r = await response.json()

        if (r.error) {
          throw r.error
        } else {
          throw `${url}: ${response.status} (${response.statusText})`
        }
      }
      
      return response.json()
    },
    async download ({ getters, commit }, data) {
      const body = {
        method: data.method,
        headers: {
          'Authorization': getters.token,
          'Content-Type': 'text/plain'
        }
      }

      if (data.body) {
        body.body = JSON.stringify(data.body)
      }

      const url = getters.backendURL+data.path
      const response = await fetch(url, body)
      
      if (response.status === 401) {
        commit('clearAuth')
      }
      
      if (!response.ok) {
        const r = await response.json()

        if (r.error) {
          throw r.error
        } else {
          throw `${url}: ${response.status} (${response.statusText})`
        }
      }
      
      return response.blob()
    },
    showToast ({ getters }) {
      setTimeout(() => {
        const toast = getters.toast.obj
        toast.show()
      }, 0)
    },
    showModal ({ getters }) {
      setTimeout(() => {
        const modal = getters.modal.obj
        modal.show()
      }, 0)
    },
    hideModal ({ getters }) {
      setTimeout(() => {
        const modal = getters.modal.obj
        modal.hide()
      }, 0)
    }
  },
  modules: {
    auth,
    client
  }
})
