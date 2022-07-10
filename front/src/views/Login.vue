<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

const store = useStore()
const router = useRouter()

const loginForm = reactive({
  username: '',
  password: ''
})

const login = () => {
  const body = {
    username: loginForm.username,
    password: loginForm.password
  }
  
  store
    .dispatch('login', body)
    .then(res => {
      if (res.token) {
        store.commit('auth', {token: res.token, username: body.username})
        router.push('/')
      } else {
        throw res
      }
    })
    .catch(e => {
      store.commit('updateToast', {color: 'danger', text: e})
      store.dispatch('showToast')
    })
}

</script>

<template>
<div class="container container-login">
  <form @submit.prevent="login" class="p-5">

    <div class="form-group pt-2">
      <label for="user_login">Username</label>
      <input
        type="text"
        name="log"
        id="user_login"
        v-model.trim="loginForm.username"
        class="form-control"
        size="20"
        required
        autofocus
      />
    </div>
    <div class="form-group pt-2">
      <label for="user_pass">Password</label>
      <input
        type="password"
        name="pwd"
        id="user_pass"
        v-model.trim="loginForm.password"
        class="form-control"
        size="20"
        required
      />
    </div>

    <div class="form-group pt-2">
      <input
        type="submit"
        name="wp-submit"
        id="wp-submit"
        class="btn btn-primary mb-4"
        value="Sign In"
      />
    </div>

  </form>
</div>
</template>

<style scoped>
.container-login{
  width: 28rem;
  height: calc(100vh - 54px);
}
</style>