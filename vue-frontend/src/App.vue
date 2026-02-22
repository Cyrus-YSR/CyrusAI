<template>
  <div id="app">
    <router-view v-slot="{ Component }">
      <transition name="page" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  const savedTheme = localStorage.getItem('theme') || 'dark'
  document.documentElement.setAttribute('data-theme', savedTheme)
})
</script>

<script>
export default {
  name: 'App'
}
</script>

<style>
:root {
  /* Dark Mode (Default) - Black/Yellow */
  --bg-primary: #000000;
  --bg-secondary: #1a1a1a;
  --bg-input: #2a2a2a;
  --bg-cyber: #050505;
  
  --text-primary: #ffffff;
  --text-secondary: #f1c40f; /* Yellow */
  --text-regular: #cccccc;
  --text-inverse: #000000;
  
  --accent-color: #f1c40f;
  --accent-hover: #f39c12;
  --accent-dark: #d4ac0d;
  
  --border-color: #333333;
  --border-focus: #f1c40f;
  
  --scrollbar-track: #1a1a1a;
  --scrollbar-thumb: #f1c40f;
  
  --cyber-grid: rgba(241, 196, 15, 0.1);
  --cyber-wave: rgba(241, 196, 15, 0.15);
  --cyber-pulse: rgba(241, 196, 15, 0.08);

  --shadow-color: rgba(0, 0, 0, 0.5);
  --header-bg: rgba(26, 26, 26, 0.9);
  --bg-glass: rgba(26, 26, 26, 0.6);
  --bg-glass-light: rgba(255, 255, 255, 0.02);
  --icon-hover-color: #ffffff;
  
  /* Element Plus Overrides */
  --el-color-primary: var(--accent-color);
}

[data-theme="light"] {
  /* Light Mode - Gray/White/Black */
  --bg-primary: #ffffff;
  --bg-secondary: #f5f7fa;
  --bg-input: #ffffff;
  --bg-cyber: #f0f2f5;
  
  --text-primary: #303133;
  --text-secondary: #606266;
  --text-regular: #909399;
  --text-inverse: #ffffff;
  
  --accent-color: #303133; /* Dark Gray/Black */
  --accent-hover: #000000;
  --accent-dark: #000000;
  
  --border-color: #dcdfe6;
  --border-focus: #303133;
  
  --scrollbar-track: #f5f7fa;
  --scrollbar-thumb: #909399;
  
  --cyber-grid: rgba(0, 0, 0, 0.05);
  --cyber-wave: rgba(0, 0, 0, 0.05);
  --cyber-pulse: rgba(0, 0, 0, 0.03);

  --shadow-color: rgba(0, 0, 0, 0.1);
  --header-bg: rgba(255, 255, 255, 0.9);
  --bg-glass: rgba(255, 255, 255, 0.6);
  --bg-glass-light: rgba(0, 0, 0, 0.02);
  --icon-hover-color: #333333;

  /* Element Plus Overrides */
  --el-color-primary: var(--accent-color);
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: background 0.3s, color 0.3s;
}

#app {
  height: 100%;
}

/* 页面切换动画 */
.page-enter-active,
.page-leave-active {
  transition: all 0.4s cubic-bezier(0.55, 0, 0.1, 1);
}

.page-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

/* 全局滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: var(--scrollbar-track);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: var(--scrollbar-thumb);
  border-radius: 4px;
  transition: background 0.3s ease;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--accent-hover);
}

/* 输入框自动填充深色模式修复 */
input:-webkit-autofill,
input:-webkit-autofill:hover,
input:-webkit-autofill:focus,
input:-webkit-autofill:active {
  -webkit-box-shadow: 0 0 0 1000px var(--bg-input) inset !important;
  -webkit-text-fill-color: var(--text-primary) !important;
  transition: background-color 5000s ease-in-out 0s;
}

/* 全局动态背景动画类 */
.cyber-bg {
  background-color: var(--bg-cyber);
  background-image: 
    linear-gradient(var(--cyber-grid) 1px, transparent 1px),
    linear-gradient(90deg, var(--cyber-grid) 1px, transparent 1px);
  background-size: 40px 40px;
  position: relative;
  overflow: hidden;
  transition: background-color 0.3s;
}

/* 扫描波浪效果 */
.cyber-bg::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 50%;
  height: 100%;
  background: linear-gradient(90deg, transparent, var(--cyber-wave), transparent);
  animation: scanWave 8s linear infinite;
  pointer-events: none;
  z-index: 0;
}

@keyframes scanWave {
  0% { left: -50%; }
  100% { left: 150%; }
}

.cyber-bg::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at 50% 50%, var(--cyber-pulse) 0%, transparent 60%);
  animation: bgPulse 10s ease-in-out infinite;
  pointer-events: none;
}

@keyframes bgPulse {
  0%, 100% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.2); }
}

/* Element Plus 组件样式覆盖 */
:root {
  --el-color-primary: var(--accent-color);
  --el-color-primary-light-3: var(--accent-hover);
  --el-color-primary-dark-2: var(--accent-dark);
}

.el-button {
  font-weight: 600;
  border-radius: 8px;
}

.el-button--primary {
  background-color: var(--accent-color) !important;
  border-color: var(--accent-color) !important;
  color: var(--text-inverse) !important;
}

.el-button--primary:hover {
  background-color: var(--accent-hover) !important;
  border-color: var(--accent-hover) !important;
  color: var(--text-inverse) !important;
}

.el-button--danger {
  background-color: var(--bg-secondary) !important;
  border-color: var(--accent-color) !important;
  color: var(--accent-color) !important;
}

.el-button--danger:hover {
  background-color: var(--border-color) !important;
  border-color: var(--accent-hover) !important;
  color: var(--accent-hover) !important;
}

.el-input__wrapper {
  background-color: var(--bg-input) !important;
  box-shadow: 0 0 0 1px var(--border-color) inset !important;
}

.el-input__wrapper.is-focus {
  box-shadow: 0 0 0 1px var(--border-focus) inset !important;
}

.el-input__inner {
  color: var(--text-primary) !important;
}

.el-card {
  background-color: var(--bg-secondary) !important;
  border-color: var(--border-color) !important;
  color: var(--text-primary) !important;
}

.el-message {
  background: var(--bg-secondary) !important;
  border-color: var(--border-color) !important;
}

.el-message--success .el-message__content {
  color: var(--accent-color) !important;
}

.el-message--error .el-message__content {
  color: #e74c3c !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-enter-from,
  .page-leave-to {
    transform: translateX(0);
    opacity: 0;
  }

  .page-enter-active,
  .page-leave-active {
    transition: opacity 0.3s ease;
  }
}
</style>