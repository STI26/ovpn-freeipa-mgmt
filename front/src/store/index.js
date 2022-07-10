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
    showToast (context) {
      setTimeout(() => {
        const toast = context.getters.toast.obj
        toast.show()
      }, 0)
    },
    showModal (context) {
      setTimeout(() => {
        const modal = context.getters.modal.obj
        modal.show()
      }, 0)
    },
    hideModal (context) {
      setTimeout(() => {
        const modal = context.getters.modal.obj
        modal.hide()
      }, 0)
    }
  },
  modules: {
    auth,
    client
  }
})
