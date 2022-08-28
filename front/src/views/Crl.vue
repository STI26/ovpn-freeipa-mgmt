<script setup>
import { onMounted, ref, computed } from 'vue'
import { useStore } from 'vuex'

const store = useStore()

const crl = ref(null)

const isInclude = (crlItem, str) => {
  const lstr = str.toLowerCase()
  return (
    crlItem.serial_number.toString().includes(lstr) ||
    crlItem.revocation_time.toString().toLowerCase().includes(lstr)
  )
}

const filteredCrl = computed(() => {
  return store.getters.searchString === ''
    ? crl.value.revoked_certificates
    : crl.value.revoked_certificates.filter((item) => isInclude(item, store.getters.searchString))
})

const loadCrl = () => {
  store
    .dispatch('getCrl')
    .then((data) => {
      crl.value = data.crl
      crl.value.revoked_certificates.sort((a, b) => (b.serial_number - a.serial_number))
    })
    .catch((e) => {
      store.commit('updateToast', { color: 'danger', text: e })
      store.dispatch('showToast')
    })
}

const update = () => {
  store
    .dispatch('updateCrl')
    .then(() => {
      loadCrl()
    })
    .catch(e => {
      store.commit('updateToast', { color: 'danger', text: e })
      store.dispatch('showToast')
    })
}

onMounted(() => {
  loadCrl()
});
</script>

<template>
  <div v-if="!crl" class="d-flex justify-content-center">
    <div class="spinner-border" role="status">
      <span class="visually-hidden">Loading...</span>
    </div>
  </div>
  <div v-else class="card text-dark bg-light mt-3">
    <div class="card-header">
      <span class="lh-lg">Certificate Revocation List</span>
      <span class="badge bg-primary btn float-end" v-on:click="update">Update</span>
      <div class="fs-09"><span class="fw-light">Issuer:</span> {{ crl.issuer }}</div>
      <div class="fs-09"><span class="fw-light">This Time:</span> {{ crl.this_time }}</div>
      <div class="fs-09"><span class="fw-light">Next Time:</span> {{ crl.next_time }}</div>
    </div>
    <div class="card-body">
      <ol class="list-group list-group-numbered">
        <li
          v-for="item of filteredCrl"
          :key="item.serial_number"
          class="list-group-item list-group-item-light d-flex justify-content-between align-items-start"
        >
          <div class="ms-2 me-auto w-100">
            <div><span class="fw-light">Serial Number:</span> {{ item.serial_number }}</div>
            <div><span class="fw-light">Revocation Time:</span> {{ item.revocation_time }}</div>
          </div>
        </li>
      </ol>
    </div>
  </div>
</template>

<style scoped>
.fs-09 {
  font-size: .9rem;
}
</style>
