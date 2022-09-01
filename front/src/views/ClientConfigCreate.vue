<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'

const store = useStore()
const router = useRouter()

const options = ref(null)
const typeTarget = ref(null)
const selectedTarget = ref(null)

const spinner = ref(false)

const loadOptions = (action) => {
  store
    .dispatch(action)
    .then((data) => {
      if (data.error !== '') {
        throw data.error
      }

      const target = action === 'getUsers' ? 'users' : 'hosts'
      options.value = data[target]

      if (data[target].length > 0) {
        selectedTarget.value = data[target][0].name
      }
    })
    .catch((e) => {
      store.commit('updateToast', { color: 'danger', text: e })
      store.dispatch('showToast')
    })
}

const onSubmit = () => {
  let action = ''
  let next = ''

  if (router.currentRoute.value.name === 'server-cert') {
    action = 'updateCert'
    next = '/config'
  } else {
    action = 'createClient'
    next = '/'
  }

  spinner.value = true
  store
    .dispatch(action, selectedTarget.value)
    .then(() => {
      spinner.value = false

      store.commit('updateToast', {
        color: 'success',
        text: 'Successful Create'
      })
      store.dispatch('showToast')
      router.push(next)
    })
    .catch((e) => {
      spinner.value = false
      store.commit('updateToast', { color: 'danger', text: e })
      store.dispatch('showToast')
    })
}

watch(
  () => typeTarget.value,
  (newTarget, oldTarget) => {
    loadOptions(newTarget)
  }
)

onMounted(() => {
  typeTarget.value = 'getUsers'
  loadOptions(typeTarget.value)
});
</script>

<template>
  <div class="card text-dark bg-light bg-gradient mt-3">
    <div class="card-header">Select name</div>
    <div class="card-body">
      <div class="row px-5">
        <div class="col-9">
          <select
            v-model="selectedTarget"
            class="form-select form-select mb-3"
            aria-label=".form-select"
          >
            <option v-for="opt in options" :key="opt.id" :value="opt.name">
              {{ opt.name }}
            </option>
          </select>
        </div>

        <div class="col-3">
          <div class="form-check">
            <input
              class="form-check-input"
              v-model="typeTarget"
              type="radio"
              value="getUsers"
              name="flexRadioDefault"
              id="flexRadioDefault1"
              checked
            />
            <label class="form-check-label" for="flexRadioDefault1">
              User
            </label>
          </div>
          <div class="form-check">
            <input
              class="form-check-input"
              v-model="typeTarget"
              type="radio"
              value="getHosts"
              name="flexRadioDefault"
              id="flexRadioDefault2"
            />
            <label class="form-check-label" for="flexRadioDefault2">
              Host
            </label>
          </div>
        </div>
      </div>

      <router-link class="btn btn-secondary float-start" to="/"
        >Cancel</router-link
      >
      <button
        type="button"
        class="btn btn-primary float-end"
        @click="onSubmit"
        :disabled="spinner"
      >
        <span
          :class="spinner ? 'spinner-border spinner-border-sm' : ''"
          role="status"
          aria-hidden="true"
        ></span>
        <span v-if="spinner"> Create...</span>
        <span v-else>Create</span>
      </button>
    </div>
  </div>
</template>

<style scoped></style>
