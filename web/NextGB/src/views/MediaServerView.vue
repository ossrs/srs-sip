<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import MediaServer from '../components/mediaserver/MediaServer.vue'
import type { MediaServer as MediaServerType } from '@/types/api'

// 扩展 MediaServer 类型
interface ExtendedMediaServer extends MediaServerType {
  isDefault?: boolean
}

const mediaServers = ref<ExtendedMediaServer[]>([])

// 表单校验规则
const rules = {
  name: [
    { required: true, message: '请输入节点名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' },
  ],
  ip: [
    { required: true, message: '请输入IP地址', trigger: 'blur' },
    { pattern: /^(\d{1,3}\.){3}\d{1,3}$/, message: '请输入正确的IP地址格式', trigger: 'blur' },
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' },
    {
      type: 'number' as const,
      min: 1,
      max: 65535,
      message: '端口号范围为1-65535',
      trigger: 'blur',
    },
  ],
}

const formRef = ref()
const dialogVisible = ref(false)
const newServer = ref<ExtendedMediaServer>({
  id: '',
  name: '',
  status: 'offline',
  ip: '',
  port: 1985,
  streams: 0,
  clients: 0,
  isDefault: false,
  type: 'SRS', // 设置默认类型
})

const handleAdd = () => {
  dialogVisible.value = true
  // 重置表单
  newServer.value = {
    id: '',
    name: '',
    status: 'offline',
    ip: '',
    port: 1985,
    streams: 0,
    clients: 0,
    isDefault: false,
    type: 'SRS', // 设置默认类型
  }
}

const handleDelete = (server: ExtendedMediaServer) => {
  ElMessage.success('删除成功')
  mediaServers.value = mediaServers.value.filter((s) => s.id !== server.id)
}

const handleSetDefault = (server: ExtendedMediaServer) => {
  // 先将所有服务器设为非默认
  mediaServers.value.forEach((s) => {
    s.isDefault = false
  })
  // 将选中的服务器设为默认
  const targetServer = mediaServers.value.find((s) => s.id === server.id)
  if (targetServer) {
    targetServer.isDefault = true
    ElMessage.success('已设为默认节点')
  }
}

const submitForm = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    // 如果新节点设为默认，先将其他节点设为非默认
    if (newServer.value.isDefault) {
      mediaServers.value.forEach((s) => {
        s.isDefault = false
      })
    }

    mediaServers.value.push({
      ...newServer.value,
      id: Date.now().toString(),
    })
    dialogVisible.value = false
    ElMessage.success('添加成功')
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

const handleClose = () => {
  formRef.value?.resetFields()
  dialogVisible.value = false
}
</script>

<template>
  <div class="media-view">
    <!-- 工具栏 -->
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>新增节点
      </el-button>
    </div>

    <!-- 节点卡片列表 -->
    <div class="server-grid">
      <MediaServer
        v-for="server in mediaServers"
        :key="server.id"
        :server="server"
        @delete="handleDelete"
        @set-default="handleSetDefault"
      />
    </div>

    <!-- 优化后的添加节点对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="添加节点"
      width="500px"
      :close-on-click-modal="false"
      @close="handleClose"
    >
      <el-form ref="formRef" :model="newServer" :rules="rules" label-width="100px" status-icon>
        <el-form-item label="节点名称" prop="name">
          <el-input v-model="newServer.name" placeholder="请输入节点名称" clearable />
        </el-form-item>

        <el-form-item label="IP地址" prop="ip">
          <el-input v-model="newServer.ip" placeholder="请输入IP地址" clearable />
        </el-form-item>

        <el-form-item label="端口" prop="port">
          <el-input-number
            v-model="newServer.port"
            :min="1"
            :max="65535"
            placeholder="请输入端口号"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="服务器类型">
          <el-radio-group v-model="newServer.type">
            <el-radio label="SRS">SRS</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="设为默认">
          <el-switch v-model="newServer.isDefault" active-text="是" inactive-text="否" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="handleClose">取 消</el-button>
        <el-button type="primary" @click="submitForm">确 定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.media-view {
  padding: 20px;
}

.toolbar {
  margin-bottom: 20px;
}

.server-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-input-number) {
  width: 100%;
}
</style>
