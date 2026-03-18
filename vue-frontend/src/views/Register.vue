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
          :password-length="registerForm.password.length"
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
          ref="registerFormRef"
          :model="registerForm"
          :rules="registerRules"
          label-position="top"
          class="login-form"
          @submit.prevent="handleRegister"
        >
          <div class="form-header">
            <h2>Create an account</h2>
            <p>Please enter your details to sign up</p>
          </div>
          
          <el-form-item label="Username" prop="username">
            <el-input
              v-model="registerForm.username"
              placeholder="请输入用户名"
              size="large"
              @focus="isTyping = true"
              @blur="isTyping = false"
            />
          </el-form-item>
          
          <el-form-item label="Email" prop="email">
            <el-input
              v-model="registerForm.email"
              placeholder="you@example.com"
              type="email"
              size="large"
              @focus="isTyping = true"
              @blur="isTyping = false"
            />
          </el-form-item>

          <el-form-item label="Verification Code" prop="captcha">
            <el-row :gutter="12" class="w-full">
              <el-col :span="16">
                <el-input
                  v-model="registerForm.captcha"
                  placeholder="请输入验证码"
                  size="large"
                  @focus="isTyping = true"
                  @blur="isTyping = false"
                />
              </el-col>
              <el-col :span="8">
                <el-button
                  type="default"
                  :loading="codeLoading"
                  :disabled="countdown > 0"
                  @click="sendCode"
                  class="captcha-btn"
                >
                  {{ countdown > 0 ? `${countdown}s` : 'Send' }}
                </el-button>
              </el-col>
            </el-row>
          </el-form-item>
          
          <el-form-item label="Password" prop="password">
            <el-input
              v-model="registerForm.password"
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
          
          <el-form-item label="Confirm Password" prop="confirmPassword">
            <el-input
              v-model="registerForm.confirmPassword"
              placeholder="••••••••"
              :type="showPassword ? 'text' : 'password'"
              size="large"
            />
          </el-form-item>

          <div class="form-actions">
            <el-button
              type="primary"
              :loading="loading"
              @click="handleRegister"
              class="submit-btn"
            >
              Sign up
            </el-button>
            <el-button class="google-btn" plain>
              <svg viewBox="0 0 24 24" width="20" height="20" xmlns="http://www.w3.org/2000/svg" class="google-icon">
                <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
                <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
                <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
                <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
              </svg>
              Sign up with Google
            </el-button>
          </div>
          
          <div class="form-footer">
            <span class="text-gray">Already have an account? </span>
            <el-button
              type="primary"
              link
              @click="$router.push('/login')"
            >
              Log in
            </el-button>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { View, Hide } from '@element-plus/icons-vue'
import api from '../utils/api'
import AnimatedCharacters from '../components/AnimatedCharacters.vue'

export default {
  name: 'RegisterView',
  components: {
    AnimatedCharacters,
    View,
    Hide
  },
  setup() {
    const router = useRouter()
    const registerFormRef = ref()
    const loading = ref(false)
    const codeLoading = ref(false)
    const countdown = ref(0)

    const registerForm = reactive({
      username: '',
      email: '',
      captcha: '',
      password: '',
      confirmPassword: ''
    })

    const isTyping = ref(false)
    const showPassword = ref(false)

    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== registerForm.password) {
        callback(new Error('两次输入密码不一致'))
      } else {
        callback()
      }
    }

    const registerRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
      ],
      email: [
        { required: true, message: '请输入邮箱', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
      ],
      captcha: [
        { required: true, message: '请输入验证码', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
      ],
      confirmPassword: [
        { required: true, message: '请确认密码', trigger: 'blur' },
        { validator: validateConfirmPassword, trigger: 'blur' }
      ]
    }

    const sendCode = async () => {
      if (!registerForm.email) {
        ElMessage.warning('请先输入邮箱')
        return
      }
      try {
        codeLoading.value = true
        const response = await api.post('/user/captcha', { email: registerForm.email })
        if (response.data.status_code === 1000) {
          ElMessage.success('验证码发送成功')
          countdown.value = 60
          const timer = setInterval(() => {
            countdown.value--
            if (countdown.value <= 0) {
              clearInterval(timer)
            }
          }, 1000)
        } else {
          ElMessage.error(response.data.status_msg || '验证码发送失败')
        }
      } catch (error) {
        console.error('Send code error:', error)
        ElMessage.error('验证码发送失败，请重试')
      } finally {
        codeLoading.value = false
      }
    }

    const handleRegister = async () => {
      try {
        await registerFormRef.value.validate()
        loading.value = true
        const response = await api.post('/user/register', {
          username: registerForm.username,
          email: registerForm.email,
          captcha: registerForm.captcha,
          password: registerForm.password
        })
        if (response.data.status_code === 1000) {
          ElMessage.success('注册成功，请登录')
          router.push('/login')
        } else {
          ElMessage.error(response.data.status_msg || '注册失败')
        }
      } catch (error) {
        console.error('Register error:', error)
        ElMessage.error('注册失败，请重试')
      } finally {
        loading.value = false
      }
    }

    return {
      registerFormRef,
      loading,
      codeLoading,
      countdown,
      registerForm,
      registerRules,
      sendCode,
      handleRegister,
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
  margin-bottom: 20px;
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

.w-full {
  width: 100%;
}

.captcha-btn {
  width: 100%;
  height: 52px;
  border-radius: 24px !important;
  font-weight: 600;
  background: #ffffff !important;
  color: #111827 !important;
  border: 1px solid #e5e7eb !important;
}

.captcha-btn:hover {
  background: #f9fafb !important;
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 32px;
  margin-top: 32px;
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

.google-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: #ffffff !important;
  color: #111827 !important;
  border: 1px solid #e5e7eb !important;
}

.google-btn:hover {
  background: #f9fafb !important;
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
