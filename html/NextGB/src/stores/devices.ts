import { ref } from 'vue'
import type { Device, ChannelInfo } from '@/api/types'
import { deviceApi } from '@/api'

// 设备列表
const devices = ref<Device[]>([])
// 通道列表
const channels = ref<ChannelInfo[]>([])
// 加载状态
const loading = ref(false)

const formatDeviceData = (device: any): Device => {
  return {
    device_id: device.device_id,
    source_addr: device.source_addr,
    network_type: device.network_type,
    status: 1,
    name: device.device_id,
  }
}

// 获取设备和通道列表
export const fetchDevicesAndChannels = async () => {
  try {
    loading.value = true
    // 获取设备列表
    const response = await deviceApi.getDevices()
    const deviceList = Array.isArray(response.data) ? response.data : []
    devices.value = deviceList.map(formatDeviceData)

    // 获取所有设备的通道
    const allChannels: ChannelInfo[] = []
    for (const device of devices.value) {
      try {
        const response = await deviceApi.getDeviceChannels(device.device_id)
        if (Array.isArray(response.data)) {
          // 确保每个通道都有正确的设备ID和其他必要属性
          const deviceChannels = response.data.map((channel: any) => ({
            ...channel,
            device_id: channel.device_id,
            status: channel.status || 'OFF',
            name: channel.name || '未命名',
            parent_id: device.device_id,
            info: {
              ...channel.info,
              ptz_type: channel.info?.ptz_type || 0,
              resolution: channel.info?.resolution || '',
              download_speed: channel.info?.download_speed || '',
            },
          }))
          allChannels.push(...deviceChannels)
        }
      } catch (error) {
        console.error(`获取设备 ${device.device_id} 的通道失败:`, error)
      }
    }
    channels.value = allChannels
  } catch (error) {
    console.error('获取设备列表失败:', error)
    throw error
  } finally {
    loading.value = false
  }
}

// 获取设备列表
export const useDevices = () => devices

// 获取通道列表
export const useChannels = () => channels

// 获取加载状态
export const useDevicesLoading = () => loading

// 更新设备列表
export const updateDevices = (newDevices: Device[]) => {
  devices.value = newDevices
}

// 更新通道列表
export const updateChannels = (newChannels: ChannelInfo[]) => {
  channels.value = newChannels
}

// 根据设备ID获取设备
export const getDeviceById = (deviceId: string) => {
  return devices.value.find((device) => device.device_id === deviceId)
}

// 根据设备ID获取该设备的所有通道
export const getChannelsByDeviceId = (deviceId: string) => {
  return channels.value.filter((channel) => channel.device_id === deviceId)
}

// 清空数据
export const clearDevicesStore = () => {
  devices.value = []
  channels.value = []
}
