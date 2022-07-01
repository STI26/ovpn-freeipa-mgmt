<script setup>
import { RouterView } from 'vue-router'
import { ref, watch } from 'vue'
import { useStore } from 'vuex'

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

        <div class="d-flex ps-5 search-form">
          <input
            class="form-control me-2"
            type="search"
            v-model.trim="query"
            placeholder="Search"
            aria-label="Search"
          />
        </div>

        <div class="d-flex flex-fill">
          <button class="btn btn-outline-success" type="button">
            Create new
          </button>
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

  <main class="container">
    <RouterView />
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
