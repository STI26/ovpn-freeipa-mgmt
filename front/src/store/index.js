import { createStore } from 'vuex'
// import auth from '@/store/auth'
import client from '@/store/client'
import { Toast } from 'bootstrap'

export default createStore({
  state: {
    backendURL: 'process.env.VUE_APP_BASEURL',
    searchString: '',
    toast: {text: '', color: ''}
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
    }
  },
  mutations: {
    updateSearchString (state, q) {
      state.searchString = q
    },
    updateToast (state, data) {
      state.toast.text = data.text
      state.toast.color = data.color
    }
  },
  actions: {
    showToast (context) {
      const toastEl = document.getElementById('liveToast')
      const toast = new Toast(toastEl)
      setTimeout(() => {
        toast.show()
      }, 0)
    }
  },
  modules: {
    // auth,
    client
  }
})
