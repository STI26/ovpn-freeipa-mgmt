import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Config from '../views/Config.vue'
import ServerConfigCreate from '../views/ServerConfigCreate.vue'
import ClientConfigCreate from '../views/ClientConfigCreate.vue'
import Connections from '../views/Connections.vue'
import About from '../views/About.vue'
import Login from '../views/Login.vue'
import NotFound from '../views/404.vue'
import store from '../store'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      meta: { onlyAuth: true, navbar: true },
      component: Home
    },
    {
      path: '/create',
      name: 'create',
      meta: { onlyAuth: true, navbar: true },
      component: ClientConfigCreate
    },
    {
      path: '/config',
      name: 'config',
      meta: { onlyAuth: true, navbar: true },
      component: Config
    },
    {
      path: '/config/create',
      name: 'config-create',
      meta: { onlyAuth: true, navbar: true },
      component: ServerConfigCreate
    },
    {
      path: '/connections',
      name: 'connections',
      meta: { onlyAuth: true, navbar: true },
      component: Connections
    },
    {
      path: '/about',
      name: 'about',
      meta: { onlyAuth: true, navbar: true },
      component: About
    },
    {
      path: '/login',
      name: 'login',
      meta: { onlyAuth: false, navbar: false },
      component: Login
    },
    {
      path: '/:pathMatch(.*)*',
      name: '404',
      meta: { onlyAuth: false, navbar: false },
      component: NotFound
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
