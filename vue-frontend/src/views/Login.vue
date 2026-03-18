<template>
  <div class="login-layout">
    <!-- 左侧插画区域 (CareerCompass 风格) -->
    <div class="illustration-section">
      <div class="brand-logo">
        <svg viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round" class="logo-icon">
          <circle cx="12" cy="12" r="10"></circle>
          <polygon points="16.24 7.76 14.12 14.12 7.76 16.24 9.88 9.88 16.24 7.76"></polygon>
        </svg>
        <span class="logo-text">CyrusAI</span>
      </div>

      <div class="characters-container">
        <AnimatedCharacters 
          :is-typing="isTyping"
          :show-password="showPassword"
          :password-length="loginForm.password.length"
        />
      </div>

      <div class="footer-links">
        <a href="#">Privacy Policy</a>
        <a href="#">Terms of Service</a>
      </div>
    </div>

    <!-- 右侧表单区域 -->
    <div class="form-section">
      <div class="login-content">
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          label-position="top"
          class="login-form"
          @submit.prevent="handleLogin"
        >
          <div class="form-header">
            <h2>Welcome back!</h2>
            <p>Please enter your details</p>
          </div>
          
          <el-form-item label="Username" prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              size="large"
              @focus="isTyping = true"
              @blur="isTyping = false"
            />
          </el-form-item>
          <el-form-item label="Password" prop="password">
            <el-input
              v-model="loginForm.password"
              placeholder="••••••••"
              :type="showPassword ? 'text' : 'password'"
              size="large"
            >
              <template #suffix>
                <el-icon class="cursor-pointer" @click="showPassword = !showPassword">
                  <View v-if="showPassword" />
                  <Hide v-else />
                </el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <div class="form-options">
            <el-checkbox v-model="rememberMe">Remember for 30 days</el-checkbox>
            <a href="#" class="forgot-link">Forgot password?</a>
          </div>

          <div class="form-actions">
            <el-button
              type="primary"
              :loading="loading"
              @click="handleLogin"
              class="submit-btn"
            >
              Log in
            </el-button>
            <div class="google-btn-wrapper">
              <GoogleLogin :callback="handleGoogleLogin" />
            </div>
          </div>
          
          <div class="form-footer">
            <span class="text-gray">Don't have an account? </span>
            <el-button
              type="primary"
              link
              @click="$router.push('/register')"
            >
              Sign Up
            </el-button>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { View, Hide } from '@element-plus/icons-vue'
import api from '../utils/api'
import AnimatedCharacters from '../components/AnimatedCharacters.vue'
import { GoogleLogin } from 'vue3-google-login'

export default {
  name: 'LoginView',
  components: {
    AnimatedCharacters,
    View,
    Hide,
    GoogleLogin
  },
  setup() {
    const router = useRouter()
    const loginFormRef = ref()
    const loading = ref(false)
    const rememberMe = ref(false)
    const loginForm = ref({
      username: '',
      password: ''
    })

    const isTyping = ref(false)
    const showPassword = ref(false)

    const loginRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
      ]
    }

    const handleLogin = async () => {
      try {
        await loginFormRef.value.validate()
        loading.value = true
        const response = await api.post('/user/login', {
          username: loginForm.value.username,
          password: loginForm.value.password
        })
        if (response.data.status_code === 1000) {
          localStorage.setItem('token', response.data.token)
          ElMessage.success('登录成功')
          router.push('/menu')
        } else {
          ElMessage.error(response.data.status_msg || '登录失败')
        }
      } catch (error) {
        console.error('Login error:', error)
        ElMessage.error('登录失败，请重试')
      } finally {
        loading.value = false
      }
    }

    const handleGoogleLogin = async (response) => {
      try {
        loading.value = true
        // api 的 baseURL 已经配置为 /api，vue.config.js 中代理会将 /api 重写为 /api/v1
        // 所以这里应该请求 /user/google-login，最终会被代理为 /api/v1/user/google-login
        const res = await api.post('/user/google-login', {
          credential: response.credential
        })
        if (res.data.status_code === 1000) {
          localStorage.setItem('token', res.data.token)
          ElMessage.success('Google登录成功')
          router.push('/menu')
        } else {
          ElMessage.error(res.data.status_msg || 'Google登录失败')
        }
      } catch (error) {
        console.error('Google login error:', error)
        ElMessage.error('Google登录失败，请重试')
      } finally {
        loading.value = false
      }
    }

    return {
      loginFormRef,
      loading,
      loginForm,
      loginRules,
      rememberMe,
      handleLogin,
      handleGoogleLogin,
      isTyping,
      showPassword
    }
  }
}
</script>

<style scoped>
.login-layout {
  display: flex;
  min-height: 100vh;
  width: 100%;
  background-color: #ffffff;
}

/* 左侧插画区域 */
.illustration-section {
  flex: 1;
  background: linear-gradient(135deg, #8a93a1 0%, #4a5568 100%);
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 40px;
  overflow: hidden;
}

@media (max-width: 900px) {
  .illustration-section {
    display: none;
  }
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  color: white;
  font-size: 1.25rem;
  font-weight: 700;
  z-index: 10;
}

.logo-icon {
  background: #1a1a1a;
  border-radius: 50%;
  padding: 4px;
}

/* 底部链接 */
.footer-links {
  display: flex;
  gap: 20px;
  z-index: 10;
}

.footer-links a {
  color: rgba(255, 255, 255, 0.5);
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.2s;
}

.footer-links a:hover {
  color: white;
}

/* 角色动画容器 */
.characters-container {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 右侧表单区域 */
.form-section {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  background: #ffffff;
}

.login-content {
  width: 100%;
  max-width: 400px;
}

.form-header {
  text-align: center;
  margin-bottom: 40px;
}

.form-header h2 {
  margin: 0 0 8px 0;
  color: #111827;
  font-size: 2rem;
  font-weight: 800;
}

.form-header p {
  margin: 0;
  color: #6b7280;
  font-size: 1rem;
}

.el-form-item {
  margin-bottom: 24px;
}

:deep(.el-form-item__label) {
  padding-bottom: 8px;
  font-weight: 600;
  color: #374151 !important;
  font-size: 0.9rem;
}

:deep(.el-input__wrapper) {
  border-radius: 24px !important;
  padding: 4px 16px !important;
  background-color: #ffffff !important;
  box-shadow: 0 0 0 1px #e5e7eb inset !important;
  transition: all 0.2s ease;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px #f3f4f6 inset !important;
}

:deep(.el-input__inner) {
  height: 44px !important;
  font-size: 1rem;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  font-size: 0.9rem;
}

.forgot-link {
  color: #4f46e5;
  text-decoration: none;
  font-weight: 500;
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 32px;
}

.submit-btn, .google-btn {
  width: 100%;
  height: 48px;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 24px !important;
  transition: all 0.2s ease;
}

.submit-btn {
  background: #ffffff !important;
  color: #111827 !important;
  border: 1px solid #e5e7eb !important;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05) !important;
}

.submit-btn:hover {
  background: #f9fafb !important;
}

.google-btn-wrapper {
  display: flex;
  justify-content: center;
  width: 100%;
}

.google-btn-wrapper > div {
  width: 100%;
}

.form-footer {
  text-align: center;
  font-size: 0.95rem;
}

.text-gray {
  color: #6b7280;
}

.form-footer .el-button {
  font-weight: 600;
  font-size: 0.95rem;
  color: #111827;
}
</style>