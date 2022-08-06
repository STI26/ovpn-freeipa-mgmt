<script setup>
import { onMounted, ref, computed } from 'vue'
import { useStore } from 'vuex'

const store = useStore()

const config = ref(null)

const isInclude = (cfgItem, str) => {
  const lstr = str.toLocaleLowerCase()
  return (cfgItem.name.toLocaleLowerCase().includes(lstr) ||
          cfgItem.value.toLocaleLowerCase().includes(lstr) ||
          cfgItem.message.toLocaleLowerCase().includes(lstr))
}

const filteredConfig = computed(() => {
  return store.getters.searchString === ''
    ? config.value
    : config.value.filter((item) => isInclude(item, store.getters.searchString))
})

const cfgToArray = (obj) => {
  const list = []

  for (const [key, val] of Object.entries(obj.value)) {
    list.push({
      name: key,
      value: val,
      status: false,
      message: obj.status[key]
    })
  }
  
  return list
} 

const loadConfig = () => {
  store
    .dispatch('getServerConfig')
    .then((data) => {
      if (data.error !== '') {
        throw data.error
      }

      config.value = cfgToArray(data.config)
    })
    .catch((e) => {
      store.commit('updateToast', { color: 'danger', text: e })
      store.dispatch('showToast')
    })
}

onMounted(() => {
  loadConfig()
})
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
      <ol v-else class="list-group list-group-numbered">
        <li v-for="item of filteredConfig" :key="item.name"
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">{{ item.name }}:</div>
            {{ item.value }} <span v-if="item.message"> ({{ item.message }})</span>
          </div>
        </li>
      </ol>
    </div>
  </div>
</template>

<style scoped></style>
