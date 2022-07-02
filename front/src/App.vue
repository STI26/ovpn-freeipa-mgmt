<script setup>
import { RouterView } from 'vue-router'
import { ref, watch } from 'vue'
import { useStore } from 'vuex'
import Toast from '@/components/Toast.vue';
import Modal from '@/components/Modal.vue';

const store = useStore()
const query = ref('')

watch(() => query.value, (newQ, oldQ) => {
  store.commit('updateSearchString', newQ)
});
</script>

<template>
  <header>
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
            Admin
          </button>
          <ul class="dropdown-menu" aria-labelledby="userMenu">
            <li><a class="dropdown-item" href="#">Logout</a></li>
          </ul>
        </div>

      </div>
    </nav>
  </header>

  <main aria-live="polite" aria-atomic="true" class="position-relative">
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
