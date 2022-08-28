<script setup>
import { RouterView, useRouter } from 'vue-router'
import { ref, watch, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import Toast from '@/components/Toast.vue'
import Modal from '@/components/Modal.vue'

const store = useStore()
const query = ref('')
const router = useRouter()

const enableNavBar = computed(() => {
  return !router.currentRoute.value.meta.navbar
})

const darkmode = computed(() => {
  return store.getters.darkmode
})

const logout = () => {
  store.commit('clearAuth')
  router.push('/login')
}

watch(
  () => router.currentRoute.value,
  () => {
    query.value = ''
  }
)

watch(
  () => query.value,
  (newQ, oldQ) => {
    store.commit('updateSearchString', newQ)
  }
)

watch(
  () => store.getters.token,
  (newT, oldT) => {
    if (newT === null) logout()
  }
)

const swTheme = () => {
  console.log(darkmode.value)
  store.commit('setDarkmode', !darkmode.value)
}

onMounted(() => {
  store.dispatch('autoLogin')
  store.dispatch('loadTheme')
});
</script>

<template>
  <header v-if="!enableNavBar">
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
          <router-link class="btn btn-outline-success" to="/create"
            >Create new</router-link
          >
        </div>

        <div class="dropdown">
          <button
            class="btn btn-outline-secondary dropdown-toggle"
            type="button"
            id="userMenu"
            data-bs-toggle="dropdown"
            aria-expanded="false"
          >
            <span class="material-icons align-top"> account_circle </span>
            {{ store.getters.username }}
          </button>
          <ul
            class="dropdown-menu dropdown-menu-end"
            aria-labelledby="userMenu"
          >
            <li>
              <button class="dropdown-item lh-lg" @click="swTheme"
                ><span class="material-icons opacity-75 align-text-bottom pe-2">
                  {{ darkmode ? 'light_mode' : 'dark_mode' }} </span
                >{{ darkmode ? 'Light mode' : 'Dark mode' }}</button
              >
            </li>
            <li>
              <router-link class="dropdown-item lh-lg" to="/connections"
                ><span class="material-icons opacity-75 align-text-bottom pe-2">
                  people </span
                >Connections</router-link
              >
            </li>
            <li>
              <router-link class="dropdown-item lh-lg" to="/config"
                ><span class="material-icons opacity-75 align-text-bottom pe-2">
                  settings </span
                >Config</router-link
              >
            </li>
            <li>
              <router-link class="dropdown-item lh-lg" to="/crl"
                ><span class="material-icons opacity-75 align-text-bottom pe-2">
                  group_remove </span
                >CRL</router-link
              >
            </li>
            <li>
              <router-link class="dropdown-item lh-lg" to="/about"
                ><span class="material-icons opacity-75 align-text-bottom pe-2">
                  info </span
                >About</router-link
              >
            </li>
            <li><hr class="dropdown-divider" /></li>
            <li>
              <button class="dropdown-item lh-lg" @click="logout" type="button">
                <span class="material-icons opacity-75 align-text-bottom pe-2">
                  logout </span
                >Logout
              </button>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  </header>

  <main
    aria-live="polite"
    aria-atomic="true"
    class="position-relative"
  >
    <div class="container">
      <RouterView />
    </div>
    <Modal />
    <Toast />
  </main>
</template>

<style lang="scss">
@import "bootstrap/scss/functions";
@import "bootstrap/scss/variables";
@import "bootstrap/scss/mixins";

.dark {
    
    $enable-gradients: true;

    /* redefine theme colors for dark theme */
    $primary: #012345;
    $secondary: #111111;
    $success: #222222;
    $dark: #000;
    
    $theme-colors: (
        "primary": $primary,
        "secondary": $secondary,
        "success": $success,
        "danger": $danger,
        "info": $indigo,
        "dark": $dark,
        "light": #aaa,
    );

    /* redefine theme color variables */
    @each $color, $value in $theme-colors {
        --#{$variable-prefix}#{$color}: #{$value};
    }
    
    $theme-colors-rgb: map-loop($theme-colors, to-rgb, "$value");
    
    @each $color, $value in $theme-colors-rgb {
        --#{$variable-prefix}#{$color}-rgb: #{$value};
    }

    $body-color: #eeeeee;
    $body-bg: #263C5C;
    
    --#{$variable-prefix}body-color: #{$body-color};
    --#{$variable-prefix}body-bg: #{$body-bg};
      
    @import "bootstrap/scss/bootstrap";
}

.search-form {
  min-width: 50%;
}
main {
  padding-top: 54px;
  min-height: 100vh;
}
</style>
