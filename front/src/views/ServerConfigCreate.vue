<script setup>
import { ref, computed, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';

const store = useStore()
const router = useRouter()
const data = reactive({
  server: '',
  vpnIP: '',
  vpnMask: ''
})

const spinner = ref(false)
const spinnerClass = computed(() => {
  return spinner.value
    ? "spinner-border spinner-border-sm"
    : ""
})

const onSubmit = () => {
  store
    .dispatch('createServer', {
      server: {
        ip: data.vpnIP,
        mask: data.vpnMask
      },
      local: data.server
    })
    .then(() => {
      spinner.value = false

      store.commit('updateToast', {color: 'success', text: 'Successful Create'})
      store.dispatch('showToast')
      router.push('/config')
    })
    .catch(e => {
      store.commit('updateToast', {color: 'danger', text: e})
      store.dispatch('showToast')
    })
}
</script>

<template>
  <form @submit.prevent="onSubmit" class="card mt-3">
    <div class="card-header">Create Server Config</div>
    <div class="card-body">
      <div class="row px-5">
        <div class="input-group mb-3">
          <span class="input-group-text">Server Address</span>
          <input
            type="text"
            class="form-control"
            v-model="data.server"
          />
        </div>
        <div class="input-group mb-3">
          <span class="input-group-text">VPN IP</span>
          <input
            type="text"
            class="form-control"
            v-model="data.vpnIP"
            minlength="7"
            maxlength="15"
            placeholder="xxx.xxx.xxx.xxx"
            aria-label="xxx.xxx.xxx.xxx"
          />
        </div>
        <div class="input-group mb-3">
          <span class="input-group-text">VPN Mask</span>
          <input
            type="text"
            class="form-control"
            v-model="data.vpnMask"
            minlength="7"
            maxlength="15"
            placeholder="xxx.xxx.xxx.xxx"
            aria-label="xxx.xxx.xxx.xxx"
          />
        </div>
      </div>

      <router-link
        class="btn btn-secondary float-start"
        to="/"
      >Cancel</router-link>
      <button
        type="submit"
        class="btn btn-primary float-end"
        :disabled="spinner"
      >
      <span :class="spinnerClass" role="status" aria-hidden="true"></span>
      <span v-if="spinner"> Create...</span>
      <span v-else>Create</span>
      </button>
    </div>
  </form>
</template>

<style scoped></style>
