<script setup>
import { onMounted, reactive } from 'vue'
import { useStore } from 'vuex'

const store = useStore()

const version = reactive({
  ipa: '',
  back: '',
  front: ''
})

onMounted(() => {
  version.front = store.getters.version
  
  store
    .dispatch('getApiVerion')
    .then((data) => {
      if (data.error !== '') {
        throw data.error
      }

      version.back = data.version
      version.ipa = data.ipa_version
    })
    .catch((e) => {
      store.commit('updateToast', { color: 'danger', text: e })
      store.dispatch('showToast')
    })
});
</script>

<template>
  <div class="card text-dark bg-light mt-3">
    <div class="card-header">About</div>
    <div class="card-body">
      <table class="table table-hover">
        <tbody>
          <tr>
            <th scope="row">IPA</th>
            <td>{{ version.ipa }}</td>
          </tr>
          <tr>
            <th scope="row">Backend</th>
            <td>{{ version.back }}</td>
          </tr>
          <tr>
            <th scope="row">Frontend</th>
            <td>{{ version.front }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped></style>
