<script setup>
import { useStore } from 'vuex'
import { reactive, watch } from 'vue'

const store = useStore()

const mData = reactive({
  title: '',
  text: '',
  action: '',
  data: {}
})

watch(() => store.getters.modal.data, (newData, oldData) => {
  mData.title = store.getters.modal.title
  mData.text = store.getters.modal.text
  mData.action = store.getters.modal.action
  mData.data = newData
})

const onSubmit = () => {
  store
    .dispatch(mData.action, mData.data)
    .then(() => {
      store.commit('updateToast', {color: 'success', text: `Successful ${mData.title}.`})
      store.dispatch('showToast')
    })
    .catch(e => {
      store.commit('updateToast', {color: 'danger', text: e})
      store.dispatch('showToast')
    })

  store.dispatch('hideModal')
}
</script>

<template>
  <div
    class="modal fade"
    id="modalApprovalForm"
    data-bs-backdrop="static"
    data-bs-keyboard="false"
    tabindex="-1"
    aria-labelledby="modalApprovalLabelForm"
    aria-hidden="true"
  >
    <div class="modal-dialog">
      <div class="modal-content text-dark bg-light">
        <div class="modal-header">
          <h5 class="modal-title" id="modalApprovalLabelForm">{{ mData.title }}</h5>
          <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
        </div>
        <div class="modal-body">
            {{ mData.text }}
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          <button type="button" class="btn btn-primary" @click="onSubmit">Ok</button>
        </div>
      </div>
    </div>
  </div>
</template>