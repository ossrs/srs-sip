<script setup lang="ts">
</script>

<template>
  <div class="settings-view">
    <div class="settings-content">
      <router-view></router-view>
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
