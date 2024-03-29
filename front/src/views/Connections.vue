<script setup>
import { onMounted, ref, computed } from 'vue'
import { useStore } from 'vuex'

const store = useStore()

const connections = ref(null)
const version = ref('')

const isInclude = (connItem, str) => {
  const lstr = str.toLocaleLowerCase()
  return (connItem.common_name.toLocaleLowerCase().includes(lstr) ||
          connItem.real_address.toLocaleLowerCase().includes(lstr) ||
          connItem.virtual_address.toLocaleLowerCase().includes(lstr))
}

const filteredConnections = computed(() => {
  return store.getters.searchString === ''
    ? connections.value
    : connections.value.filter((item) => isInclude(item, store.getters.searchString))
})

const loadConnections = () => {
  store
    .dispatch('getConnections')
    .then((data) => {
      if (data.error !== '') {
        throw data.error
      }

      connections.value = data.status.client_list
      version.value = data.status.title
    })
    .catch((e) => {
      store.commit('updateToast', { color: 'danger', text: e })
      store.dispatch('showToast')
    })
}

const update = () => {
  loadConnections()
}

onMounted(() => {
  loadConnections()
});
</script>

<template>
  <div v-if="!version" class="d-flex justify-content-center">
    <div class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
  <div v-else class="card text-dark bg-light bg-gradient mt-3">
    <div class="card-header">
      Active Clients <span class="fw-lighter">({{ version }})</span>
      <span class="badge bg-primary btn float-end" v-on:click="update">Update</span>
    </div>
    <div class="card-body">
      <table class="table table-ligth text-dark table-hover">
        <thead>
          <tr>
            <th scope="col">#</th>
            <th scope="col">Common Name</th>
            <th scope="col">Real Address</th>
            <th scope="col">Virtual Address</th>
            <th scope="col">Bytes Received</th>
            <th scope="col">Bytes Send</th>
            <th scope="col">Connected Since</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(client, idx) in filteredConnections"
            :key="client.common_name + client.connected_since"
          >
            <th scope="row">{{ idx + 1 }}</th>
            <td>{{ client.common_name }}</td>
            <td>{{ client.real_address }}</td>
            <td>{{ client.virtual_address }}</td>
            <td>{{ Number(client.bytes_received).toLocaleString() }}</td>
            <td>{{ Number(client.bytes_send).toLocaleString() }}</td>
            <td>{{ client.connected_since }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped></style>
