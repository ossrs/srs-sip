import { createRouter, createWebHistory } from 'vue-router'
import RealplayView from '@/views/realplay/RealplayView.vue'
import SettingsView from '@/views/setting/SettingsView.vue'
import PlaybackView from '@/views/playback/PlaybackView.vue'
import MediaServerView from '@/views/mediaserver/MediaServerView.vue'
import DashboardView from '@/views/DashboardView.vue'
import DeviceListView from '@/views/DeviceListView.vue'
import SystemForm from '@/views/setting/SystemForm.vue'
import About from '@/views/setting/About.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
    },
    {
      path: '/realplay',
      name: 'realplay',
      component: RealplayView,
    },
    {
      path: '/devices',
      name: 'devices',
      component: DeviceListView,
    },
    {
      path: '/settings',
      component: SettingsView,
      children: [
        {
          path: '',
          redirect: 'system'
        },
        {
          path: 'system',
          name: 'settings-system',
          component: SystemForm
        },
        {
          path: 'about',
          name: 'settings-about',
          component: About
        }
      ]
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
