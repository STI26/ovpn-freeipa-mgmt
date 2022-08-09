<script setup>
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router';
import { useStore } from 'vuex'

const router = useRouter()
const store = useStore()

const config = ref(null)
const configError = ref('')

const isInclude = (cfgItem, str) => {
  const lstr = str.toLocaleLowerCase()
  return (
    cfgItem.name.toLocaleLowerCase().includes(lstr) ||
    cfgItem.value.toLocaleLowerCase().includes(lstr) ||
    cfgItem.message.toLocaleLowerCase().includes(lstr)
  )
}

const filteredConfig = computed(() => {
  return store.getters.searchString === ''
    ? config.value
    : config.value.filter((item) => isInclude(item, store.getters.searchString))
})

const showUpdateBtn = (key) => {
  switch (key) {
    case 'ca':
    case 'crl':
    case 'cert':
    case 'key':
      return true
    default:
      return false
  }
}

const update = (key) => {
  let action = ''

  switch (key) {
    case 'ca':
      action = 'updateCA'
      break
    case 'crl':
      action = 'updateCrl'
      break
    case 'cert':
    case 'key':
      router.push('/server/cert')
      return
    default:
      return
  }

  store
    .dispatch(action)
    .then(() => {
      loadConfig()
    })
    .catch(e => {
      store.commit('updateToast', { color: 'danger', text: e })
      store.dispatch('showToast')
    })
}

const cfgToArray = (obj) => {
  const list = []

  for (const [key, val] of Object.entries(obj.value)) {
    list.push({
      key: key,
      name: obj.name[key],
      value: val,
      status: obj.status[key],
      message: obj.message[key]
    })
  }

  return list
}

const loadConfig = () => {
  store
    .dispatch('getServerConfig')
    .then((data) => {
      configError.value = data.error
      config.value = cfgToArray(data.config)
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
  <div class="card mt-3">
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
          class="btn btn-outline-primary float-end"
          to="/config/create"
          >Create</router-link
        >
      </div>
      <ol v-else class="list-group list-group-numbered">
        <li
          v-for="item of filteredConfig"
          :key="item.name"
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto w-100">
            <div class="fw-bold">
              <span v-if="item.status === 'true'" class="text-success"
                >&#10003;</span
              >
              <span v-else-if="item.status === 'false'" class="text-danger"
                >&#10007;</span
              >
              {{ item.name || item.key }}:
              <span
                v-if="showUpdateBtn(item.key)"
                class="badge bg-primary btn float-end"
                v-on:click="update(item.key)"
                >Update</span
              >
            </div>
            {{ item.value }}
            <span v-if="item.message"> ({{ item.message }})</span>
          </div>
        </li>
      </ol>
    </div>
  </div>
</template>

<style scoped></style>
