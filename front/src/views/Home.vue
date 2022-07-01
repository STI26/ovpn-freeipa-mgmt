<script setup>
import Client from '@/components/Client.vue';
import Clients from '@/components/Clients.vue';
</script>

<template>
  <div class="pt-2">
    <div class="row">
      <div class="col scrolled">
        <Clients msg="The clients" />
      </div>
      <div class="col">
        <Client msg="The client" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.scrolled {
  max-height: 100vh;
  overflow-y: auto;
}
</style>>