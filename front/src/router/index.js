import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Create from '../views/Create.vue'
import Config from '../views/Config.vue'
import Login from '../views/Login.vue'
import store from '../store'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      meta: { onlyAuth: true },
      component: Home
    },
    {
      path: '/create',
      name: 'create',
      meta: { onlyAuth: true },
      component: Create
    },
    {
      path: '/config',
      name: 'config',
      meta: { onlyAuth: true },
      component: Config
    },
    {
      path: '/login',
      name: 'login',
      meta: { onlyAuth: false },
      component: Login
    }
  ]
})

router.beforeEach((to, from, next) => {
  const ifAuthenticated = store.getters.ifAuthenticated
  const onlyAuth = to.matched.some(r => r.meta.onlyAuth)

  if (onlyAuth && !ifAuthenticated) {
    router.push('/login')
  } else {
    next()
  }
})

export default router
