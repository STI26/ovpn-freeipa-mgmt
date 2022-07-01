<script setup>
import { ref, computed } from 'vue'

defineProps({
  msg: {
    type: String,
    required: true
  }
})

const clients = ref([
  {id: 0, name: 'client 0', numberOfCertificates: 1, active: false},
  {id: 1, name: 'client 1', numberOfCertificates: 0, active: true},
  {id: 2, name: 'client 2', numberOfCertificates: 2, active: false},
  {id: 3, name: 'eclient 3', numberOfCertificates: 1, active: false},
  {id: 4, name: 'client 4', numberOfCertificates: 3, active: false},
  {id: 11, name: 'client 11', numberOfCertificates: 1, active: false},
  {id: 12, name: 'j.dou', numberOfCertificates: 1, active: false}
])

const q = 'cli  '.trim()
const filteredClients = computed(() => {
  return q === ''
    ? clients.value
    : clients.value.filter((item) => item.name.includes(q))
})

const showClient = (client)=> {
  clients.value.forEach((c) => {
    c.active = false
  })
  client.active = true
} 
</script>

<template>
  <div class="list-group">
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
