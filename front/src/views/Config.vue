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
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">Address:</div>
            {{ config.value.Local }}
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">Port:</div>
            {{ config.value.Port }}
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">Proto:</div>
            {{ config.value.Proto }}
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">TunMtu:</div>
            {{ config.value.TunMtu }}
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">DataCiphers:</div>
            {{ config.value.DataCiphers }}
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">Auth:</div>
            {{ config.value.Auth }}
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">CompLzo:</div>
            {{ config.value.CompLzo }}
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">CA:</div>
            {{ config.value.CA }} ({{ config.status.CA }})
          </div>
          <span class="badge bg-primary btn">Update</span>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">Crl:</div>
            {{ config.value.Crl }} ({{ config.status.Crl }})
          </div>
          <span class="badge bg-primary btn">Update</span>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">TlsAuth:</div>
            {{ config.value.TlsAuth }} ({{ config.status.TlsAuth }})
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">Cert:</div>
            {{ config.value.Cert }} ({{ config.status.Cert }})
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">Key:</div>
            {{ config.value.Key }} ({{ config.status.Key }})
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">Ipp:</div>
            {{ config.value.Ipp }} ({{ config.status.Ipp }})
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">Ccd:</div>
            {{ config.value.Ccd }} ({{ config.status.Ccd }})
          </div>
        </li>
        <li
          class="list-group-item d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto">
            <div class="fw-bold">Status:</div>
            {{ config.value.Status }} ({{ config.status.Status }})
          </div>
        </li>
      </ol>
    </div>
  </div>
</template>

<style scoped></style>
