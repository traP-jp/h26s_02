import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Top',
    component: () => import('@/pages/timeline/TimelineView.vue'),
  },
  {
    path: '/camera',
    name: 'Camera',
    component: () => import('@/pages/camera/CameraView.vue'),
  },
  {
    path: '/camera/yakudo',
    name: 'CameraYakudo',
    component: () => import('@/pages/camera/YakudoView.vue'),
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('@/pages/grid/SearchView.vue'),
  },
  {
    path: '/user/:userId',
    name: 'User',
    component: () => import('@/pages/grid/UserView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
