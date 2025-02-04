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
        <span>NextGB</span>
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
          <el-menu-item index="settings" @click="$router.push('/settings/system')">基本设置</el-menu-item>
          <el-menu-item index="about" @click="$router.push('/settings/about')">关于</el-menu-item>
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
          <a href="https://github.com/ossrs/srs-sip" target="_blank" class="github-link">
            <el-icon><svg viewBox="0 0 1024 1024" width="1em" height="1em"><path fill="currentColor" d="M511.6 76.3C264.3 76.2 64 276.4 64 523.5 64 718.9 189.3 885 363.8 946c23.5 5.9 19.9-10.8 19.9-22.2v-77.5c-135.7 15.9-141.2-73.9-150.3-88.9C215 726 171.5 718 184.5 703c30.9-15.9 62.4 4 98.9 57.9 26.4 39.1 77.9 32.5 104 26 5.7-23.5 17.9-44.5 34.7-60.8-140.6-25.2-199.2-111-199.2-213 0-49.5 16.3-95 48.3-131.7-20.4-60.5 1.9-112.3 4.9-120 58.1-5.2 118.5 41.6 123.2 45.3 33-8.9 70.7-13.6 112.9-13.6 42.4 0 80.2 4.9 113.5 13.9 11.3-8.6 67.3-48.8 121.3-43.9 2.9 7.7 24.7 58.3 5.5 118 32.4 36.8 48.9 82.7 48.9 132.3 0 102.2-59 188.1-200 212.9a127.5 127.5 0 0 1 38.1 91v112.5c.8 9 0 17.9 15 17.9 177.1-59.7 304.6-227 304.6-424.1 0-247.2-200.4-447.3-447.5-447.3z"></path></svg></el-icon>
          </a>
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

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.github-link {
  color: #606266;
  font-size: 20px;
  display: flex;
  align-items: center;
  text-decoration: none;
  transition: color 0.3s;
}

.github-link:hover {
  color: #1890ff;
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
