import { ref } from 'vue'
import type { MediaServer } from '@/types/api'
import { mediaServerApi } from '@/api'
import { ElMessage } from 'element-plus'

// 所有媒体服务器列表
const mediaServers = ref<MediaServer[]>([])
// 默认媒体服务器
const defaultMediaServer = ref<MediaServer | null>(null)
// 加载状态
const loading = ref(false)

// 检查服务器状态
export const checkServersStatus = async () => {
  for (const server of mediaServers.value) {
    try {
      await mediaServerApi.checkStatus(server)
      server.status = 1 // 在线
    } catch (error) {
      server.status = 0 // 离线
    }
  }
}

// 获取媒体服务器列表
export const fetchMediaServers = async () => {
  try {
    loading.value = true
    const response = await mediaServerApi.getMediaServers()
    // 确保 mediaServers 始终是数组，并将 is_default 映射为 isDefault
    mediaServers.value = Array.isArray(response.data) ? response.data.map((server: any) => ({
      ...server,
      isDefault: server.is_default
    })) : []
    
    if (mediaServers.value.length > 0) {
      await checkServersStatus()
      // 找到默认服务器并更新 defaultMediaServer
      const defaultServer = mediaServers.value.find(server => server.isDefault === 1)
      if (defaultServer) {
        defaultMediaServer.value = defaultServer
      }
    }
  } catch (error) {
    console.error('获取媒体服务器列表失败:', error)
    mediaServers.value = [] // 出错时也清空列表
  } finally {
    loading.value = false
  }
}

// 设置默认媒体服务器
export const setDefaultMediaServer = async (server: MediaServer) => {
  try {
    // 调用后端API设置默认服务器
    await mediaServerApi.setDefaultMediaServer(server.id)
    
    // 更新前端状态
    mediaServers.value.forEach((s) => {
      s.isDefault = s.id === server.id ? 1 : 0
    })
    
    // 更新默认服务器引用
    defaultMediaServer.value = mediaServers.value.find(s => s.id === server.id) || null
    
    ElMessage.success('已设为默认节点')
  } catch (error) {
    console.error('设置默认服务器失败:', error)
    ElMessage.error('设置默认节点失败')
    throw error
  }
}

// 删除媒体服务器
export const deleteMediaServer = async (server: MediaServer) => {
  try {
    // 如果要删除的是默认服务器，先清除默认服务器状态
    if (server.isDefault) {
      defaultMediaServer.value = null
    }
    
    await mediaServerApi.deleteMediaServer(server.id)
    ElMessage.success('删除成功')
    await fetchMediaServers()
  } catch (error) {
    console.error('删除失败:', error)
    throw error
  }
}


// 获取所有媒体服务器列表
export const useMediaServers = () => mediaServers
// 获取默认媒体服务器
export const useDefaultMediaServer = () => defaultMediaServer
// 获取加载状态
export const useMediaServersLoading = () => loading 