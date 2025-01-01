import { createRouter, createWebHistory } from 'vue-router'
import MonitorView from '@/views/MonitorView.vue'
import SettingsView from '@/views/SettingsView.vue'
import PlaybackView from '@/views/PlaybackView.vue'
import MediaServerView from '@/views/MediaServerView.vue'
import DashboardView from '@/views/DashboardView.vue'
import DeviceList from '@/views/DeviceList.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
    },
    {
      path: '/monitor',
      name: 'monitor',
      component: MonitorView,
    },
    {
      path: '/devices',
      name: 'devices',
      component: DeviceList,
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView,
    },
    {
      path: '/playback',
      name: 'playback',
      component: PlaybackView,
    },
    {
      path: '/media',
      name: 'media',
      component: MediaServerView,
    },
  ],
})

export default router
