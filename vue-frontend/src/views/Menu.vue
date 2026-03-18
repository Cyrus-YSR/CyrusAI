<template>
  <div class="menu-container">
    <nav class="glass-nav">
      <div class="nav-content">
        <div class="brand">
          <LogoEye :size="40" />
          <span class="brand-text">CyrusAI</span>
        </div>
        <div class="nav-actions">
          <div class="user-profile">
            <div class="avatar">U</div>
            <span class="username">User</span>
          </div>
          <button class="logout-btn" @click="handleLogout">
            <el-icon><SwitchButton /></el-icon>
          </button>
        </div>
      </div>
    </nav>

    <main class="dashboard-main">
      <!-- Hero Section with Characters -->
      <div class="hero-section">
        <div class="hero-text">
          <h1>Hello, Traveler!</h1>
          <p>Ready to explore the digital universe?</p>
        </div>
        <div class="hero-visual">
          <!-- Reusing the characters component for brand consistency -->
          <div class="mini-characters">
            <AnimatedCharacters :is-typing="false" :show-password="false" :password-length="0" />
          </div>
        </div>
      </div>

      <!-- Bento Grid Menu -->
      <div class="bento-grid">
        <!-- AI Chat Card (Large) -->
        <div class="bento-card large chat-card" @click="$router.push('/ai-chat')">
          <div class="card-bg-decoration"></div>
          <div class="card-content">
            <div class="icon-box chat-icon-box">
              <el-icon><ChatDotRound /></el-icon>
            </div>
            <div class="card-text">
              <h2>AI Assistant</h2>
              <p>Chat with our smart AI. It can see, hear, and help.</p>
            </div>
            <div class="hover-arrow">→</div>
          </div>
        </div>

        <!-- Image Rec Card (Medium) -->
        <div class="bento-card medium camera-card" @click="$router.push('/image-recognition')">
          <div class="card-content">
            <div class="icon-box camera-icon-box">
              <el-icon><Camera /></el-icon>
            </div>
            <div class="card-text">
              <h2>Vision</h2>
              <p>Analyze images instantly.</p>
            </div>
          </div>
        </div>

        <!-- Placeholder Card (Small - Fun) -->
        <div class="bento-card small fun-card">
          <div class="fun-content">
            <span class="emoji">🚀</span>
            <span>More coming soon...</span>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ChatDotRound, Camera, SwitchButton } from '@element-plus/icons-vue'
import AnimatedCharacters from '../components/AnimatedCharacters.vue'
import LogoEye from '../components/LogoEye.vue'

export default {
  name: 'MenuView',
  components: {
    ChatDotRound,
    Camera,
    SwitchButton,
    AnimatedCharacters,
    LogoEye
  },
  setup() {
    const router = useRouter()

    const handleLogout = async () => {
      try {
        await ElMessageBox.confirm('Ready to disconnect?', 'Logout', {
          confirmButtonText: 'Yes, bye!',
          cancelButtonText: 'Stay',
          type: 'warning',
          customClass: 'fun-confirm'
        })
        localStorage.removeItem('token')
        ElMessage.success('See you next time!')
        router.push('/login')
      } catch {
        // User stayed
      }
    }

    return {
      handleLogout
    }
  }
}
</script>

<style scoped>
.menu-container {
  min-height: 100vh;
  /* background-color: #F3F4F6; Removed to show App.vue global background */
  font-family: 'Inter', sans-serif;
  overflow-x: hidden;
  padding-bottom: 40px;
}

/* Glass Navigation */
.glass-nav {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  width: 90%;
  max-width: 1200px;
  height: 70px;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.4);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.05);
  z-index: 100;
  transition: all 0.3s ease;
}

.glass-nav:hover {
  background: rgba(255, 255, 255, 0.8);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.08);
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  padding: 0 24px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Removed inline Logo Eye styles as it is now a component */

.brand-text {
  font-weight: 800;
  font-size: 1.25rem;
  color: #111827;
  letter-spacing: -0.02em;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.5);
  padding: 4px 12px 4px 4px;
  border-radius: 30px;
  border: 1px solid rgba(255, 255, 255, 0.6);
}

.avatar {
  width: 32px;
  height: 32px;
  background: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: #111827;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.username {
  font-size: 0.9rem;
  font-weight: 600;
  color: #374151;
}

.logout-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: none;
  background: rgba(254, 226, 226, 0.8);
  color: #EF4444;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  backdrop-filter: blur(4px);
}

.logout-btn:hover {
  background: #EF4444;
  color: white;
  transform: rotate(90deg);
}

/* Dashboard Main */
.dashboard-main {
  padding-top: 120px;
  max-width: 1200px;
  margin: 0 auto;
  padding-left: 20px;
  padding-right: 20px;
  padding-bottom: 40px;
}

.hero-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 60px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(16px);
  border-radius: 32px;
  padding: 40px 60px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.6);
}

.hero-text h1 {
  font-size: 3.5rem;
  font-weight: 900;
  color: #111827;
  margin: 0 0 16px 0;
  line-height: 1.1;
  background: linear-gradient(135deg, #111827 0%, #4B5563 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-text p {
  font-size: 1.25rem;
  color: #6B7280;
  margin: 0;
}

.hero-visual {
  position: relative;
  width: 300px;
  height: 200px;
}

.mini-characters {
  transform: scale(0.6);
  transform-origin: center center;
  position: absolute;
  top: -80px;
  right: -80px;
}

/* Bento Grid */
.bento-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  grid-template-rows: 280px;
  gap: 24px;
}

.bento-card {
  background: rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(12px);
  border-radius: 32px;
  padding: 32px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.6);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.02);
}

.bento-card:hover {
  transform: translateY(-8px) scale(1.02);
  background: rgba(255, 255, 255, 0.85);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  border-color: rgba(255, 255, 255, 0.9);
}

/* Chat Card */
.chat-card {
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.8) 0%, rgba(224, 231, 255, 0.4) 100%);
  border: 1px solid rgba(255, 255, 255, 0.6);
}

.chat-card:hover .icon-box {
  transform: scale(1.1) rotate(-5deg);
  background: #4F46E5;
  color: white;
  box-shadow: 0 10px 20px rgba(79, 70, 229, 0.3);
}

.chat-card:hover .hover-arrow {
  opacity: 1;
  transform: translateX(0);
}

.icon-box {
  width: 64px;
  height: 64px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  margin-bottom: 24px;
  transition: all 0.3s ease;
}

.chat-icon-box {
  background: rgba(224, 231, 255, 0.8);
  color: #4F46E5;
}

.card-text h2 {
  font-size: 2rem;
  font-weight: 800;
  margin: 0 0 8px 0;
  color: #111827;
}

.card-text p {
  color: #4B5563;
  font-size: 1.1rem;
  font-weight: 500;
}

.hover-arrow {
  position: absolute;
  bottom: 32px;
  right: 32px;
  font-size: 2rem;
  color: #4F46E5;
  opacity: 0;
  transform: translateX(-20px);
  transition: all 0.3s ease;
}

/* Camera Card */
.camera-card {
  background: linear-gradient(145deg, rgba(255, 255, 255, 0.8) 0%, rgba(253, 231, 243, 0.4) 100%);
  border: 1px solid rgba(255, 255, 255, 0.6);
}

.camera-card:hover .icon-box {
  transform: scale(1.1) rotate(5deg);
  background: #DB2777;
  color: white;
  box-shadow: 0 10px 20px rgba(219, 39, 119, 0.3);
}

.camera-icon-box {
  background: rgba(252, 231, 243, 0.8);
  color: #DB2777;
}

.camera-card h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #111827;
  margin-bottom: 4px;
}

/* Fun Card */
.fun-card {
  background: rgba(255, 255, 255, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px dashed rgba(209, 213, 219, 0.6);
}

.fun-content {
  text-align: center;
  color: #6B7280;
  font-weight: 600;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.emoji {
  font-size: 2rem;
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

/* Responsive */
@media (max-width: 900px) {
  .bento-grid {
    grid-template-columns: 1fr;
    grid-template-rows: auto;
  }
  
  .hero-section {
    flex-direction: column;
    text-align: center;
    padding: 32px;
  }
  
  .hero-visual {
    margin-top: 40px;
    height: 150px;
  }
  
  .mini-characters {
    top: -120px;
    right: 50%;
    transform: translateX(50%) scale(0.4);
  }
}
</style>