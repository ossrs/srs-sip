<script setup lang="ts">
import SystemForm from '@/views/setting/SystemForm.vue'
</script>

<template>
  <div class="settings-view">
    <h1>系统设置</h1>
    <div class="settings-content">
      <SystemForm />
    </div>
  </div>
</template>

<style scoped>
.settings-view {
  padding: 20px;
}

.settings-content {
  margin-top: 20px;
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}
</style>
