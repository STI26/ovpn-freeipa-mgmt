<script setup>
import { computed } from 'vue'
import { useStore } from 'vuex'

const store = useStore()

const bgClass = computed(() => {
  let classes = 'bg-' + store.getters.toast.color

  if (store.getters.toast.color !== 'warning' && store.getters.toast.color !== 'info') {
    classes += ' text-white'
  }

  return classes
})

const text = computed(() => {
  return store.getters.toast.text
})
</script>

<template>
<div class="toast-container position-absolute top-0 end-0 p-3 mt-5">
  <div
    id="liveToast"
    class="toast align-items-center"
    :class="bgClass"
    role="alert"
    aria-live="assertive"
    aria-atomic="true"
  >
    <div class="d-flex">
      <div class="toast-body">
        {{ text }}
      </div>
      <button type="button" class="btn-close me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
    </div>
  </div>
</div>
</template>