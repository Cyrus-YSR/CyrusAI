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
// Theme logic removed
</script>

<script>
export default {
  name: 'App'
}
</script>

<style>
:root {
  /* CareerCompass Inspired Theme - Animated & Professional */
  
  /* Backgrounds */
  --bg-primary: #F0F2F5; /* Light Grayish Blue */
  --bg-secondary: #FFFFFF;
  --bg-input: #F9FAFB;
  
  /* Text */
  --text-primary: #1F2937; /* Gray-900 */
  --text-secondary: #6B7280; /* Gray-500 */
  --text-regular: #9CA3AF;
  --text-inverse: #FFFFFF;
  
  /* Accents - Modern Indigo/Blue Gradient feel */
  --accent-color: #4F46E5; /* Indigo-600 */
  --accent-hover: #4338CA; /* Indigo-700 */
  --accent-light: #E0E7FF;
  
  /* Borders */
  --border-color: #E5E7EB;
  --border-focus: #4F46E5;
  
  /* Shadows */
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  --card-shadow: var(--shadow-lg);
  
  /* Element Plus Overrides */
  --el-color-primary: var(--accent-color);
  --el-text-color-primary: var(--text-primary);
  --el-bg-color: var(--bg-secondary);
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: var(--bg-primary);
  color: var(--text-primary);
  overflow-x: hidden;
}

/* Animated Background */
body::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    radial-gradient(at 0% 0%, rgba(79, 70, 229, 0.15) 0px, transparent 50%),
    radial-gradient(at 100% 0%, rgba(236, 72, 153, 0.15) 0px, transparent 50%),
    radial-gradient(at 100% 100%, rgba(59, 130, 246, 0.15) 0px, transparent 50%),
    radial-gradient(at 0% 100%, rgba(16, 185, 129, 0.15) 0px, transparent 50%);
  z-index: -1;
  animation: bgAnimation 20s ease-in-out infinite alternate;
  background-size: 200% 200%;
}

@keyframes bgAnimation {
  0% { background-position: 0% 0%; }
  100% { background-position: 100% 100%; }
}

#app {
  height: 100%;
}

/* Page Transitions */
.page-enter-active,
.page-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Scrollbar */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: #D1D5DB;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #9CA3AF;
}

/* Element Plus Overrides */
.el-button {
  font-weight: 500;
  border-radius: 8px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.el-button--primary {
  background-color: var(--accent-color) !important;
  border-color: var(--accent-color) !important;
  color: var(--text-inverse) !important;
  box-shadow: 0 4px 6px -1px rgba(79, 70, 229, 0.2);
}

.el-button--primary:hover {
  background-color: var(--accent-hover) !important;
  border-color: var(--accent-hover) !important;
  transform: translateY(-1px);
  box-shadow: 0 6px 8px -1px rgba(79, 70, 229, 0.3);
}

.el-button--primary:active {
  transform: translateY(0);
}

.el-input__wrapper {
  background-color: var(--bg-input) !important;
  box-shadow: 0 0 0 1px var(--border-color) inset !important;
  border-radius: 8px !important;
  padding: 8px 12px !important;
  transition: all 0.2s ease !important;
}

.el-input__wrapper.is-focus {
  box-shadow: 0 0 0 2px var(--accent-color) inset !important;
  background-color: #fff !important;
}

.el-card {
  background-color: var(--bg-secondary) !important;
  border: 1px solid var(--border-color) !important;
  color: var(--text-primary) !important;
  box-shadow: var(--card-shadow) !important;
  border-radius: 16px !important;
  overflow: visible !important; /* Allow hover effects to show */
}

/* Responsive */
@media (max-width: 768px) {
  .page-enter-active,
  .page-leave-active {
    transition: none;
  }
}
</style>