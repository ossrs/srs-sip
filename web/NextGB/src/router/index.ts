import { createRouter, createWebHistory } from 'vue-router'
import MonitorView from '@/views/MonitorView.vue'
import DevicesView from '@/views/DevicesView.vue'
import SettingsView from '@/views/SettingsView.vue'
import PlaybackView from '@/views/PlaybackView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'monitor',
      component: MonitorView
    },
    {
      path: '/devices',
      name: 'devices',
      component: DevicesView
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView
    },
    {
      path: '/playback',
      name: 'playback',
      component: PlaybackView
    }
  ]
})

export default router
