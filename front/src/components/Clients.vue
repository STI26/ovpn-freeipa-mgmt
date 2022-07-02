<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useStore } from 'vuex'

const store = useStore()
const clients = ref([])

const filteredClients = computed(() => {
  return store.getters.searchString === ''
    ? clients.value
    : clients.value.filter((item) => item.name.includes(store.getters.searchString))
})

watch(() => store.getters.getClientID, (newID, oldID) => {
  clients.value.forEach((c) => {
    if (c.id === newID) {
      c.active = true
    } else {
      c.active = false
    }
  })
})

const showClient = (client)=> {
  store.commit('updateClientID', client.id)
}

onMounted(() => {
  store
    .dispatch('getClients')
    .then(data => {
      clients.value = data
    })
})
</script>

<template>
  <div v-if="filteredClients.length === 0">
    <h5>clients not found</h5>
  </div>
  <div v-else class="list-group">
    <button
      type="button"
      v-for="client in filteredClients" :key="client.id"
      class="list-group-item d-flex justify-content-between align-items-center list-group-item-action"
      :class="{ active: client.active }"
      @click="showClient(client)"
    >
      {{ client.name }}
      <span class="badge bg-secondary rounded-pill">{{ client.numberOfCertificates }}</span>
    </button>
  </div>
</template>
