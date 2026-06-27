import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Top',
    component: () => import('@/pages/TopView.vue'),
  },
  {
    path: '/timeline',
    name: 'Timeline',
    component: () => import('@/pages/timeline/TimelineView.vue'),
  },
  {
    path: '/samples',
    name: 'Samples',
    component: () => import('@/samples/SamplesView.vue'),
  },
  {
    path: '/samples/camera',
    name: 'Camera',
    component: () => import('@/samples/CameraView.vue'),
  },
  {
    path: '/samples/motion',
    name: 'Motion',
    component: () => import('@/samples/MotionView.vue'),
  },
  {
    path: '/samples/scroll',
    name: 'Scroll',
    component: () => import('@/samples/ScrollView.vue'),
  },
  {
    path: '/samples/yakudo',
    name: 'Yakudo',
    component: () => import('@/samples/YakudoView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
