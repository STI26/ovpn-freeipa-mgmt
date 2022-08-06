<script setup>
import { computed, watch, onMounted, reactive } from 'vue'
import { useStore } from 'vuex'

const store = useStore()
const clients = reactive({
  users: [],
  hosts: []
})
const spinner = reactive({
  users: true,
  hosts: true
})

const isInclude = (clientItem, str) => {
  return clientItem.name.toLocaleLowerCase().includes(str.toLocaleLowerCase())
}

const filteredUsers = computed(() => {
  return store.getters.searchString === ''
    ? clients.users
    : clients.users.filter((item) => isInclude(item, store.getters.searchString))
})

const filteredHosts = computed(() => {
  return store.getters.searchString === ''
    ? clients.hosts
    : clients.hosts.filter((item) => isInclude(item, store.getters.searchString))
})

watch(() => store.getters.getClientID, (newID, oldID) => {
  clients.users.forEach((c) => {
    if (c.id === newID) {
      c.active = true
    } else {
      c.active = false
    }
  })

  clients.hosts.forEach((c) => {
    if (c.id === newID) {
      c.active = true
    } else {
      c.active = false
    }
  })
})

const showClient = (client)=> {
  store.commit('updateClientID', client.name)
}

onMounted(() => {
  store
    .dispatch('getUsers')
    .then(res => {
      if (!res.users) {
        throw 'can\'t get users object'
      } else {
        spinner.users = false
        clients.users = res.users.filter((item) => item.numberOfCertificates > 0)
      }
    })
    .catch(e => {
      store.commit('updateToast', {color: 'danger', text: e})
      store.dispatch('showToast')
    })

  store
    .dispatch('getHosts')
    .then(res => {
      if (!res.hosts) {
        throw 'can\'t get hosts object'
      } else {
        spinner.hosts = false
        clients.hosts = res.hosts.filter((item) => item.numberOfCertificates > 0)
      }
    })
    .catch(e => {
      store.commit('updateToast', {color: 'danger', text: e})
      store.dispatch('showToast')
    })
})
</script>

<template>

  <div class="card text-center">
    <h5 class="card-header text-muted">Users</h5>
  </div>

  <div v-if="spinner.users" class="d-flex justify-content-center">
    <div class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
  <div v-else-if="filteredUsers.length === 0">
    <h5>users not found</h5>
  </div>
  <div v-else class="list-group">
    <button
      type="button"
      v-for="user in filteredUsers" :key="user.id"
      class="list-group-item d-flex justify-content-between align-items-center list-group-item-action"
      :class="{ active: user.active }"
      @click="showClient(user)"
    >
      {{ user.name }}
      <span class="badge bg-secondary rounded-pill">{{ user.numberOfCertificates }}</span>
    </button>
  </div>

  <div class="card text-center">
    <h5 class="card-header text-muted">Hosts</h5>
  </div>

  <div v-if="spinner.hosts" class="d-flex justify-content-center">
    <div class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
  <div v-else-if="filteredHosts.length === 0">
    <h5>hosts not found</h5>
  </div>
  <div v-else class="list-group">
    <button
      type="button"
      v-for="host in filteredHosts" :key="host.id"
      class="list-group-item d-flex justify-content-between align-items-center list-group-item-action"
      :class="{ active: host.active }"
      @click="showClient(host)"
    >
      {{ host.name }}
      <span class="badge bg-secondary rounded-pill">{{ host.numberOfCertificates }}</span>
    </button>
  </div>

</template>
