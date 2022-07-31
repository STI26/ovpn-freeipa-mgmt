<script setup>
import { RouterView, useRouter } from 'vue-router'
import { ref, watch, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import Toast from '@/components/Toast.vue';
import Modal from '@/components/Modal.vue';

const store = useStore()
const query = ref('')
const router = useRouter()

const mainbg = computed(() => {
  return router.currentRoute.value.name === 'login'
    ? "bg-light"
    : ""
})

const logout = () => {
  store.commit('clearAuth')
  router.push('/login')
}

watch(() => query.value, (newQ, oldQ) => {
  store.commit('updateSearchString', newQ)
})

watch(() => store.getters.token, (newT, oldT) => {
  if (newT === null) logout()
})

onMounted(() => {
  store.dispatch('autoLogin')
})
</script>

<template>
  <header v-if="!mainbg">
    <nav class="navbar fixed-top navbar-light bg-light">
      <div class="container-fluid">

        <div class="d-flex">
          <router-link class="btn btn-outline-dark" to="/">/</router-link>
        </div>

        <div class="d-flex ps-3 search-form">
          <input
            class="form-control me-2"
            type="search"
            v-model.trim="query"
            placeholder="Search"
            aria-label="Search"
          />
        </div>

        <div class="d-flex flex-fill">
          <router-link class="btn btn-outline-success" to="/create">Create new</router-link>
        </div>

        <div class="dropdown">
          <button
            class="btn btn-outline-secondary dropdown-toggle"
            type="button"
            id="userMenu"
            data-bs-toggle="dropdown"
            aria-expanded="false"
          >
            {{ store.getters.username }}
          </button>
          <ul class="dropdown-menu dropdown-menu-end" aria-labelledby="userMenu">
            <li><button class="dropdown-item" @click="logout" type="button">Logout</button></li>
          </ul>
        </div>

      </div>
    </nav>
  </header>

  <main aria-live="polite" aria-atomic="true" class="position-relative" :class="mainbg">
    <div class="container">
      <RouterView />
    </div>
    <Modal />
    <Toast />
  </main>
</template>

<style>
.search-form {
  min-width: 50%;
}
main {
  padding-top: 54px;
}
</style>
