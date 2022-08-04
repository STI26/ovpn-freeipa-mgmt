<script setup>
import { onMounted, ref, watch, computed } from 'vue'
import { useStore } from 'vuex'

const store = useStore()

const config = ref(null)

const loadConfig = () => {
  store
    .dispatch('getServerConfig')
    .then((data) => {
      if (data.error !== '') {
        throw data.error
      }

      config.value = data.config
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
        <li v-for="[key, val] of Object.entries(config.value)" :key="key"
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">{{ key }}:</div>
            {{ val }} <span v-if="config.status[key]"> ({{ config.status[key] }})</span>
          </div>
        </li>
      </ol>
    </div>
  </div>
</template>

<style scoped></style>
