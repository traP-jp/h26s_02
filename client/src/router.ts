import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Top',
    component: () => import('@/pages/TopView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
