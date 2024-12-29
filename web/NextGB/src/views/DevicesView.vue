<script setup lang="ts">
import DeviceList from '@/components/devices/DeviceList.vue'
</script>

<template>
  <div class="devices-view">
    <h1>设备管理</h1>
    <DeviceList />
  </div>
</template>

<style scoped>
.devices-view {
  height: 100%;
}

h1 {
  margin-bottom: 20px;
}
</style>
