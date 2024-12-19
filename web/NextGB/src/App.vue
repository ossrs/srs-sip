<script setup lang="ts">
import { ref } from 'vue'
import { Monitor, Setting, Tools, Fold, VideoCamera } from '@element-plus/icons-vue'

const isCollapse = ref(false)

const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}
</script>

<template>
  <div class="app-container">
    <!-- 左侧菜单 -->
    <div class="sidebar" :class="{ 'is-collapse': isCollapse }">
      <div class="logo">
        <img src="./assets/logo.svg" alt="Logo" />
        <span>demo</span>
      </div>
      <el-menu :collapse="isCollapse" default-active="1" class="sidebar-menu">
        <el-menu-item index="1" @click="$router.push('/')">
          <el-icon><Monitor /></el-icon>
          <span>实时监控</span>
        </el-menu-item>

        <el-sub-menu index="2">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>设备管理</span>
          </template>
          <el-menu-item index="2-1" @click="$router.push('/devices')">设备列表</el-menu-item>
          <el-menu-item index="2-2">设备状态</el-menu-item>
        </el-sub-menu>

        <el-menu-item index="3" @click="$router.push('/settings')">
          <el-icon><Tools /></el-icon>
          <span>系统设置</span>
        </el-menu-item>
      </el-menu>
    </div>

    <!-- 右侧内容区 -->
    <div class="main-container">
      <!-- 顶部导航 -->
      <div class="header">
        <div class="header-left">
          <el-button @click="toggleSidebar">
            <el-icon><Fold /></el-icon>
          </el-button>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              <img class="avatar" src="./assets/avatar.png" />
              <span>管理员</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人信息</el-dropdown-item>
                <el-dropdown-item>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- 主要内容区域 -->
      <div class="main-content">
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  height: 100vh;
  width: 100%;
  overflow-x: hidden; /* 防止横向溢出 */
}

.sidebar {
  width: 240px;
  background-color: #304156;
  color: #fff;
  transition: width 0.3s;
}

.sidebar.is-collapse {
  width: 64px;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  overflow: hidden;
}

.logo img {
  width: 32px;
  height: 32px;
}

.logo span {
  margin-left: 12px;
  font-size: 16px;
  font-weight: 600;
  transition: opacity 0.3s;
}

.is-collapse .logo span {
  opacity: 0;
  display: none;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.3s;
  width: 100%;
}

.sidebar.is-collapse + .main-container {
}

.header {
  height: 60px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  margin-right: 8px;
}

.main-content {
  flex: 1;
  padding: 20px;
  background: #f0f2f5;
  overflow-y: auto;
}
</style>