<script setup>
defineProps({
  msg: {
    type: String,
    required: true
  }
})
</script>

<template>
  <form>
    <div class="row">
      <div class="col-md-8 mb-3">
        <label for="clientCerts" class="form-label">Certificate</label>
        <select class="form-select" id="clientCerts" aria-label="Certificate">
          <option value="1" selected>12</option>
          <option value="2">14</option>
          <option value="3">41</option>
        </select>
      </div>
      <div class="col-md mb-3">
        <label for="revokeCert" class="form-label">Revoke</label>
        <button type="button" class="form-control btn btn-outline-danger" id="revokeCert">revoke</button>
      </div>
    </div>
    <div class="row">
      <div class="col-md mb-3">
        <label for="clientCertStatus" class="form-label">Status</label>
        <input type="text" readonly class="form-control" id="clientCertStatus" placeholder="revoked">
      </div>
      <div class="col-md mb-3">
        <label for="clientCertExpiresOn" class="form-label">Expires On</label>
        <input type="text" readonly class="form-control" id="clientCertExpiresOn" placeholder="Fri Jan 12 00:04:26 2024 UTC">
      </div>
    </div>
    <div class="row p-3">
      <hr class="dropdown-divider" />
    </div>
    <div class="row">
      <div class="col-lg-8 mb-3">
        <label for="clientName" class="form-label">Subject</label>
        <input type="text" readonly class="form-control" id="clientName" placeholder="l.name">
      </div>
      <div class="col-lg-4 mb-3">
        <label for="clientConfig" class="form-label">Download config</label>
        <button type="button" class="form-control btn btn-outline-dark" id="clientConfig">config.ovpn</button>
      </div>
    </div>
    <div class="row">
      <div class="col mb-3">
        <label for="clientIP" class="form-label">IP</label>
        <input type="text" class="form-control" id="clientIP" placeholder="10.10.10.1">
      </div>
    </div>
    <div class="row">
      <div class="col mb-3">
        <label for="clientRoutes" class="form-label">Routes</label>
        <textarea class="form-control" id="clientRoutes" rows="3"></textarea>
      </div>
    </div>
    <div class="row pt-2">
      <div class="col mb-3">
        <button type="button" class="btn btn-danger float-start">Delete</button>
      </div>
      <div class="col mb-3">
        <button type="submit" class="btn btn-primary float-end">Save</button>
      </div>
    </div>
  </form>
</template>

<style scoped>
h1 {
  font-weight: 500;
  font-size: 2.6rem;
  top: -10px;
}

h3 {
  font-size: 1.2rem;
}

.greetings h1,
.greetings h3 {
  text-align: center;
}

@media (min-width: 1024px) {
  .greetings h1,
  .greetings h3 {
    text-align: left;
  }
}
</style>
