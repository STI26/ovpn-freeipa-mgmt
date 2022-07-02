<script setup>
import { reactive, watch } from 'vue'
import { useStore } from 'vuex'

const store = useStore()
const client = reactive({
  id: null,
  subject: '',
  ip: '',
  routes: '',
  certificates: [],
  certStatus: '',
  certExpiresOn: '',
  selectedCertificate: null
})

const clearForm = () => {
  store.commit('updateClientID', null)
  
  client.id = null
  client.subject = ''
  client.ip = ''
  client.routes = ''
  client.certificates = []
  client.selectedCertificate = null
}

const downloadConfig = () => {
  if (client.selectedCertificate === null) {
    store.commit('updateToast', {color: 'warning', text: 'Certificate not selected.'})
    store.dispatch('showToast')
    return
  }

  store
    .dispatch('downloadConfig', {
      clientID: client.id,
      subject: client.subject,
      certificateID: client.selectedCertificate
    })
    .catch(e => {
      store.commit('updateToast', {color: 'danger', text: e})
      store.dispatch('showToast')
    })
}

const saveClient = () => {
  const data = {
    title: 'Save',
    text: 'These files will be modify:',
    action: 'saveClient',
    data: {
      clientID: client.id,
      subject: client.subject,
      ip: client.ip,
      routes: client.routes,
      certificateID: client.selectedCertificate
    }
  }
  store.commit('updateModal', data)
  store.dispatch('showModal')
}

const deleteClient = () => {
  const data = {
    title: 'Delete',
    text: 'These files will be removed:',
    action: 'deleteClient',
    data: {subject: client.subject}
  }
  store.commit('updateModal', data)
  store.dispatch('showModal')
}

const revokeCert = () => {
  const data = {
    title: 'Revoke',
    text: `Сertificate ${client.selectedCertificate} belonging to the ${client.subject} will be revoked.`,
    action: 'revokeCert',
    data: {id: client.selectedCertificate}
  }
  store.commit('updateModal', data)
  store.dispatch('showModal')
}

watch(() => store.getters.getClientID, (newID, oldID) => {
  store
    .dispatch('getClient', newID)
    .then(data => {
      client.id = data.id
      client.subject = data.subject
      client.ip = data.ip
      client.routes = data.routes
      client.certificates = data.certificates

      if (data.certificates.length > 0) {
        client.selectedCertificate = data.certificates[0].id
      } else {
        client.selectedCertificate = null
      }
    })
    .catch(err => {
      clearForm()
    })
})

watch(() => client.selectedCertificate, (newCertID, oldCertID) => {
  const cert = client.certificates.filter(c => c.id === newCertID)

  if (cert.length === 0) {
    client.certStatus = ''
    client.certExpiresOn = ''
  } else {
    client.certStatus = cert[0].revoked ? 'revoked' : 'valid' 
    client.certExpiresOn = cert[0].expiresOn
  }
})
</script>

<template>
  <div v-if="client.id === null">
    <h3>nothing selected</h3>
  </div>
  <form v-else @submit.prevent="saveClient">
    <div class="row">
      <div class="col mb-3">
        <button type="button" class="btn-close float-end" @click="clearForm" aria-label="Close"></button>
      </div>
    </div>
    <div class="row">
      <div class="col-md-8 mb-3">
        <label for="clientCerts" class="form-label">Certificate</label>
        <select class="form-select" id="clientCerts" v-model="client.selectedCertificate" aria-label="Certificate">
          <option
            v-for="certificate in client.certificates"
            :key="certificate.id"
            :value="certificate.id"
          >{{ certificate.id + ' (' + certificate.issuedOn + ')'}}</option>
        </select>
      </div>
      <div class="col-md mb-3">
        <label for="revokeCert" class="form-label">Revoke</label>
        <button
          type="button"
          class="form-control btn btn-outline-danger"
          id="revokeCert"
          @click="revokeCert"
          :disabled="client.selectedCertificate === null"
        >revoke</button>
      </div>
    </div>
    <div class="row">
      <div class="col-md mb-3">
        <label for="clientCertStatus" class="form-label">Status</label>
        <input
          type="text"
          readonly
          class="form-control"
          id="clientCertStatus"
          v-model.cm="client.certStatus"
        />
      </div>
      <div class="col-md mb-3">
        <label for="clientCertExpiresOn" class="form-label">Expires On</label>
        <input
          type="text"
          readonly
          class="form-control"
          id="clientCertExpiresOn"
          v-model.trim="client.certExpiresOn"
        />
      </div>
    </div>
    <div class="row p-3">
      <hr class="dropdown-divider" />
    </div>
    <div class="row">
      <div class="col-lg-8 mb-3">
        <label for="clientName" class="form-label">Subject</label>
        <input
          type="text"
          readonly
          class="form-control"
          id="clientName"
          v-model.trim="client.subject"
        />
      </div>
      <div class="col-lg-4 mb-3">
        <label for="clientConfig" class="form-label">Download config</label>
        <button
          type="button"
          class="form-control btn btn-outline-dark"
          id="clientConfig"
          @click="downloadConfig"
        >config.ovpn</button>
      </div>
    </div>
    <div class="row">
      <div class="col mb-3">
        <label for="clientIP" class="form-label">IP</label>
        <input
          type="text"
          class="form-control"
          id="clientIP"
          v-model.trim="client.ip"
          placeholder="10.10.10.1"
        />
      </div>
    </div>
    <div class="row">
      <div class="col mb-3">
        <label for="clientRoutes" class="form-label">Routes</label>
        <textarea
          class="form-control"
          id="clientRoutes"
          v-model.trim="client.routes"
          rows="3"
        ></textarea>
      </div>
    </div>
    <div class="row pt-2">
      <div class="col mb-3">
        <button
          type="button"
          class="btn btn-danger float-start"
          @click="deleteClient"
        >Delete</button>
      </div>
      <div class="col mb-3">
        <button
          type="submit"
          class="btn btn-primary float-end"
        >Save</button>
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
