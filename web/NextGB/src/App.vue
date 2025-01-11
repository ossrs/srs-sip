<script setup lang="ts">
import { ref, provide, onMounted } from 'vue'
import {
  Monitor,
  Setting,
  Tools,
  Fold,
  VideoCamera,
  User,
  VideoPlay,
  DataLine,
} from '@element-plus/icons-vue'
import { useDefaultMediaServer } from '@/stores/mediaServer'
import { fetchDevicesAndChannels } from '@/stores/devices'
import { fetchMediaServers } from '@/stores/mediaServer'

const isCollapse = ref(false)

const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

// 提供默认媒体服务器
const defaultMediaServer = useDefaultMediaServer()
provide('defaultMediaServer', defaultMediaServer)

// 初始化数据
const initializeData = async () => {
  try {
    // 并行获取设备列表和媒体服务器列表
    await Promise.all([fetchDevicesAndChannels(), fetchMediaServers()])
  } catch (error) {
    console.error('初始化数据失败:', error)
  }
}

onMounted(() => {
  initializeData()
})
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
        <el-menu-item index="dashboard" @click="$router.push('/dashboard')">
          <el-icon><DataLine /></el-icon>
          <span>系统概览</span>
        </el-menu-item>

        <el-menu-item index="realplay" @click="$router.push('/realplay')">
          <el-icon><Monitor /></el-icon>
          <span>实时监控</span>
        </el-menu-item>

        <el-menu-item index="playback" @click="$router.push('/playback')">
          <el-icon><VideoPlay /></el-icon>
          <span>录像回放</span>
        </el-menu-item>

        <el-menu-item index="media" @click="$router.push('/media')">
          <el-icon><VideoCamera /></el-icon>
          <span>流媒体服务</span>
        </el-menu-item>

        <el-sub-menu index="device">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>设备管理</span>
          </template>
          <el-menu-item index="device-list" @click="$router.push('/devices')">
            设备列表
          </el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="system">
          <template #title>
            <el-icon><Tools /></el-icon>
            <span>系统设置</span>
          </template>
          <el-menu-item index="settings" @click="$router.push('/settings')">基本设置</el-menu-item>
        </el-sub-menu>
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
              <el-icon class="avatar-icon"><User /></el-icon>
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
        <router-view v-slot="{ Component }">
          <keep-alive :include="['RealplayView', 'PlaybackView']">
            <component :is="Component" />
          </keep-alive>
        </router-view>
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
  box-shadow: 2px 0 6px rgba(0, 21, 41, 0.15);
  z-index: 10;
  display: flex;
  flex-direction: column;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  overflow: hidden;
  background-color: #2b3a4d;
  border-bottom: 1px solid #1f2d3d;
}

.logo img {
  width: 32px;
  height: 32px;
  transition: margin 0.3s;
}

.logo span {
  margin-left: 12px;
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  transition: opacity 0.3s;
  white-space: nowrap;
}

.sidebar-menu {
  border-right: none !important;
  background-color: transparent;
}

/* 自定义菜单样式 */
:deep(.el-menu) {
  border-right: none;
}

:deep(.el-menu-item) {
  height: 50px;
  line-height: 50px;
  color: #bfcbd9;

  &:hover {
    background-color: #263445 !important;
  }

  &.is-active {
    background-color: #1890ff !important;
    color: #fff;
  }
}

:deep(.el-sub-menu__title) {
  height: 50px;
  line-height: 50px;
  color: #bfcbd9;

  &:hover {
    background-color: #263445 !important;
  }
}

:deep(.el-menu--collapse) {
  width: 64px;

  .el-sub-menu__title span {
    display: none;
  }

  .el-sub-menu__title .el-sub-menu__icon-arrow {
    display: none;
  }
}

/* 折叠状态下的样式 */
.is-collapse {
  width: 64px;

  .logo {
    padding: 0 16px;

    img {
      margin: 0;
    }

    span {
      opacity: 0;
      display: none;
    }
  }
}

/* 图标样式 */
:deep(.el-menu-item .el-icon),
:deep(.el-sub-menu__title .el-icon) {
  font-size: 18px;
  margin-right: 12px;
  vertical-align: middle;
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
  color: #333;

  .avatar-icon {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    margin-right: 8px;
    background-color: #e6e6e6;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
  }
}

.main-content {
  flex: 1;
  padding: 20px;
  background: #f0f2f5;
  overflow-y: auto;
}

/* 修改二级菜单样式 */
:deep(.el-sub-menu) {
  .el-menu {
    background-color: #1f2d3d !important;
  }

  .el-menu-item {
    padding-left: 54px !important;
    height: 44px;
    line-height: 44px;
    font-size: 13px;

    &:hover {
      background-color: #001528 !important;
    }

    &.is-active {
      background-color: #1890ff !important;
      &::before {
        content: '';
        position: absolute;
        left: 0;
        top: 0;
        bottom: 0;
        width: 3px;
        background-color: #fff;
      }
    }
  }
}

/* 优化子菜单标题样式 */
:deep(.el-sub-menu__title) {
  &:hover {
    background-color: #263445 !important;
  }

  .el-sub-menu__icon-arrow {
    right: 15px;
    margin-top: -4px;
    font-size: 12px;
    transition: transform 0.3s;
  }
}

/* 展开状态的箭头动画 */
:deep(.el-sub-menu.is-opened) {
  > .el-sub-menu__title {
    color: #f4f4f5;

    .el-sub-menu__icon-arrow {
      transform: rotateZ(180deg);
    }
  }
}

/* 折叠状态下的弹出菜单样式 */
:deep(.el-menu--popup) {
  background-color: #1f2d3d !important;
  padding: 0;

  .el-menu-item {
    height: 44px;
    line-height: 44px;
    font-size: 13px;
    padding: 0 20px !important;
    color: #bfcbd9;

    &:hover {
      background-color: #001528 !important;
    }

    &.is-active {
      background-color: #1890ff !important;
      color: #fff;
    }
  }
}

/* 修改菜单过渡动画 */
:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  transition:
    background-color 0.3s,
    color 0.3s,
    border-color 0.3s;
}
</style>
