<script setup>
import { onMounted, ref, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

const router = useRouter()
const store = useStore()

const config = ref(null)
const configError = ref('')
const updateStatus = reactive({
  ca: false,
  crl_verify: false,
  dh: false,
  tls_auth: false,
  cert: false,
  key: false
})

const isInclude = (cfgItem, str) => {
  const lstr = str.toLocaleLowerCase()
  return (
    cfgItem.key.toLocaleLowerCase().includes(lstr) ||
    cfgItem.value.toLocaleLowerCase().includes(lstr) ||
    cfgItem.description.toLocaleLowerCase().includes(lstr)
  )
}

const filteredConfig = computed(() => {
  return store.getters.searchString === ''
    ? config.value
    : config.value.filter((item) => isInclude(item, store.getters.searchString))
})

const showUpdateBtn = (key) => {
  key = key.replace('-', '_')
  return Object.keys(updateStatus).includes(key)
}

const update = (key) => {
  key = key.replace('-', '_')

  if (updateStatus[key]) {
    return
  }

  let action = ''

  switch (key) {
    case 'ca':
      action = 'updateCA'
      break
    case 'crl_verify':
      action = 'updateCrl'
      break
    case 'dh':
      action = 'updateDH'
      break
    case 'tls_auth':
      action = 'updateTlsAuth'
      break
    case 'cert':
    case 'key':
      router.push('/server/cert')
      return
    default:
      return
  }

  updateStatus[key] = true

  store
    .dispatch(action)
    .then(() => {
      updateStatus[key] = false
      loadConfig()
    })
    .catch((e) => {
      store.commit('updateToast', { color: 'danger', text: e })
      store.dispatch('showToast')
    })
}

const loadConfig = () => {
  store
    .dispatch('getServerConfig')
    .then((data) => {
      configError.value = data.error
      config.value = data.config
    })
    .catch((e) => {
      store.commit('updateToast', { color: 'danger', text: e })
      store.dispatch('showToast')
    })
}

onMounted(() => {
  loadConfig()
});
</script>

<template>
  <div class="card text-dark bg-light bg-gradient mt-3">
    <div class="card-header">Server Config</div>
    <div class="card-body">
      <div v-if="!config" class="d-flex justify-content-center">
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      <div v-else-if="configError">
        {{ configError }}
        <router-link
          class="btn btn-outline-danger float-end"
          to="/config/create"
          >Create</router-link
        >
      </div>
      <ol v-else class="list-group list-group-numbered">
        <li
          v-for="(item, i) of filteredConfig"
          :key="item.key + i"
          class="list-group-item list-group-item-light d-flex justify-content-between align-items-start position-relative"
        >
          <div class="ms-2 me-auto w-100">
            <div class="fw-bold">
              <span v-if="item.status" class="text-success bg-light rounded-pill">&#10003;</span>
              <span v-else-if="('status' in item) && !item.status" class="text-danger bg-light rounded-pill">&#10007;</span>
              {{ item.key }}:
              <span
                v-if="showUpdateBtn(item.key)"
                class="badge bg-primary btn float-end"
                v-on:click="update(item.key)"
              >
                <span
                  :class="
                    updateStatus[item.key]
                      ? 'spinner-border spinner-border-sm'
                      : ''
                  "
                  role="status"
                  aria-hidden="true"
                ></span>
                <span v-if="!updateStatus[item.key]">Update</span></span
              >
            </div>
            {{ item.value }}
            <span class="fw-light" v-if="item.note"> ({{ item.note }})</span>
          </div>
          <span
            class="btn position-absolute top-10 start-0 translate-middle badge rounded-pill bg-info"
            data-bs-toggle="tooltip"
            data-bs-placement="top"
            :title="item.description"
          >
            ?
          </span>
        </li>
      </ol>
    </div>
  </div>
</template>

<style scoped></style>
