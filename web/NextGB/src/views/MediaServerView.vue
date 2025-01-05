<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import MediaServerCard from '../components/mediaserver/MediaServerCard.vue'
import type { MediaServer } from '@/types/api'
import { mediaServerApi } from '@/api'
import { 
  useMediaServers,
  useDefaultMediaServer,
  fetchMediaServers,
  setDefaultMediaServer,
  deleteMediaServer,
  checkServersStatus
} from '@/stores/mediaServer'

const mediaServers = useMediaServers()
const defaultMediaServer = useDefaultMediaServer()
let statusCheckTimer: number | null = null

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
  type: [
    { required: true, message: '请选择服务器类型', trigger: 'change' },
  ],
}

const formRef = ref()
const dialogVisible = ref(false)
const newServer = ref<Pick<MediaServer, 'name' | 'ip' | 'port' | 'type' | 'username' | 'password' | 'isDefault'>>({
  name: '',
  ip: '',
  port: 1985,
  type: 'SRS',
  username: '',
  password: '',
  isDefault: 0,
})

const handleAdd = () => {
  dialogVisible.value = true
  // 重置表单
  newServer.value = {
    name: '',
    ip: '',
    port: 1985,
    type: 'SRS',
    username: '',
    password: '',
    isDefault: 0,
  }
}

const handleDelete = async (server: MediaServer) => {
  try {
    await deleteMediaServer(server)
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleSetDefault = (server: MediaServer) => {
  setDefaultMediaServer(server)
}

const submitForm = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    await mediaServerApi.addMediaServer({
      name: newServer.value.name,
      ip: newServer.value.ip,
      port: newServer.value.port,
      type: newServer.value.type,
      username: newServer.value.username,
      password: newServer.value.password,
    })
    
    dialogVisible.value = false
    ElMessage.success('添加成功')
    await fetchMediaServers()
  } catch (error) {
    console.error('添加失败:', error)
  }
}

const handleClose = () => {
  formRef.value?.resetFields()
  dialogVisible.value = false
}

// 组件挂载时获取数据
onMounted(() => {
  // 延迟3秒后获取服务器状态
  setTimeout(() => {
    checkServersStatus()
  }, 3000)
})
</script>

<template>
  <div class="media-view">
    <!-- 工具栏 -->
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>新增节点
      </el-button>
      <el-button @click="fetchMediaServers">
        <el-icon><Refresh /></el-icon>刷新
      </el-button>
    </div>

    <!-- 节点卡片列表 -->
    <div class="server-grid">
      <MediaServerCard
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
            <el-radio value="SRS">SRS</el-radio>
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
