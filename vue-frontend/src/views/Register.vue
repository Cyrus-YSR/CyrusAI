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
            <InteractiveHoverButton
              type="submit"
              :text="loading ? 'Signing up...' : 'Sign up'"
              :disabled="loading"
              class="submit-btn"
              @click="handleRegister"
            />
            <div style="margin-top: 16px;">
              <InteractiveHoverButton
                type="button"
                text="Sign up with Google"
                class="google-btn"
                @click="handleGoogleLogin"
              >
                <template #icon>
                  <svg style="width: 20px; height: 20px;" aria-hidden="true" focusable="false" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 488 512">
                    <path fill="currentColor" d="M488 261.8C488 403.3 391.1 504 248 504 110.8 504 0 393.2 0 256S110.8 8 248 8c66.8 0 126 23.4 172.9 61.9l-76.2 76.2C322.3 113.2 289.4 96 248 96c-88.8 0-160.1 71.9-160.1 160.1s71.3 160.1 160.1 160.1c98.2 0 135-70.4 140.8-106.9H248v-85.3h236.1c2.3 12.7 3.9 26.9 3.9 41.4z"></path>
                  </svg>
                </template>
              </InteractiveHoverButton>
            </div>
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
import InteractiveHoverButton from '../components/InteractiveHoverButton.vue'
import { googleTokenLogin } from 'vue3-google-login'

export default {
  name: 'RegisterView',
  components: {
    AnimatedCharacters,
    InteractiveHoverButton,
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

    const handleGoogleLogin = async () => {
      try {
        loading.value = true
        const response = await googleTokenLogin()
        if (!response.access_token) {
          throw new Error('No access token received')
        }

        const res = await api.post('/user/google-login', {
          credential: response.access_token,
          is_access_token: true
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
      registerFormRef,
      loading,
      codeLoading,
      countdown,
      registerForm,
      registerRules,
      sendCode,
      handleRegister,
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
  height: 48px;
  border-radius: 24px !important;
  font-weight: 600 !important;
  font-size: 15px !important;
  background: #ffffff !important;
  color: #111827 !important;
  border: 1px solid #e5e7eb !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  padding: 0 !important;
  cursor: pointer !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
}

.captcha-btn:hover {
  background: #f9fafb !important;
  border-color: #d1d5db !important;
  transform: translateY(-2px) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08) !important;
}

.form-actions {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 8px;
  margin-bottom: 32px;
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
