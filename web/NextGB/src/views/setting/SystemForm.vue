<script setup lang="ts">
import { ref } from 'vue'

const formData = ref({
  systemName: 'demo',
  storagePath: '/data/recordings',
  retention: 30,
  autoCleanup: true,
})
</script>

<template>
  <el-form :model="formData" label-width="120px">
    <el-form-item label="系统名称">
      <el-input v-model="formData.systemName" />
    </el-form-item>
    <el-form-item label="存储路径">
      <el-input v-model="formData.storagePath" />
    </el-form-item>
    <el-form-item label="保留天数">
      <el-input-number v-model="formData.retention" :min="1" :max="365" />
    </el-form-item>
    <el-form-item label="自动清理">
      <el-switch v-model="formData.autoCleanup" />
    </el-form-item>
    <el-form-item>
      <el-button type="primary">保存设置</el-button>
    </el-form-item>
  </el-form>
</template>

<style scoped>
.el-form {
  max-width: 600px;
}
</style>
