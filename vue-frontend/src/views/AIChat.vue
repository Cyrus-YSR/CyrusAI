<template>
  <div class="ai-chat-container">
    <!-- 左侧会话列表 -->
    <div class="session-list">
      <div class="session-list-header">
        <span>会话列表</span>
        <button class="new-chat-btn" @click="createNewSession">＋ 新聊天</button>
      </div>
      <ul class="session-list-ul">
        <li
          v-for="session in sessions"
          :key="session.id"
          :class="['session-item', { active: currentSessionId === session.id }]"
          @click="switchSession(session.id)"
        >
          <span class="session-title">
            {{ session.name || `会话 ${session.id}` }}
          </span>
          <button class="session-delete-btn" @click.stop="deleteSession(session.id)">×</button>
        </li>
      </ul>
    </div>

    <!-- 右侧聊天区域 -->
    <div class="chat-section">
      <div class="top-bar">
        <button class="back-btn" @click="$router.push('/menu')">← 返回</button>
        
        <div class="top-bar-right">
          <button class="sync-btn" @click="syncHistory" :disabled="!currentSessionId || tempSession">同步历史数据</button>
          
          <div class="model-select-wrapper">
            <span class="label">模型：</span>
            <el-select v-model="selectedModel" class="custom-select" popper-class="custom-dropdown">
              <el-option label="阿里百炼" value="1" />
              <el-option label="阿里百炼 RAG" value="2" />
              <el-option label="阿里百炼 MCP" value="3" />
            </el-select>
          </div>

          <button class="upload-btn" @click="triggerFileUpload" :disabled="uploading">📎 上传文档(.md/.txt)</button>
          <input
            ref="fileInput"
            type="file"
            accept=".md,.txt,text/markdown,text/plain"
            style="display: none"
            @change="handleFileUpload"
          />
          <ThemeToggle />
        </div>
      </div>

      <div class="chat-messages" ref="messagesRef">
        <div
          v-for="(message, index) in currentMessages"
          :key="index"
          :class="['message', message.role === 'user' ? 'user-message' : 'ai-message']"
        >
          <div class="message-header">
            <b>{{ message.role === 'user' ? 'me' : 'CyrusAI' }}:</b>
            <button v-if="message.role === 'assistant'" class="tts-btn" @click="playTTS(message.content)">
              <svg class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" width="14" height="14">
                <path d="M836.56 382.89a344.17 344.17 0 0 0-70.19-110.61 45.79 45.79 0 0 0-32-14.72 51.33 51.33 0 0 0-37 13.95c-19.18 17.62-20.07 48.05-1.92 67.89l0.12 0.12a250.8 250.8 0 0 1 0 344.39l-0.52 0.58a53.34 53.34 0 0 0-10.42 19.65h-2.68l0.6 14a50.5 50.5 0 0 0 15.71 34.61 47.54 47.54 0 0 0 33.29 13.72h1.07a47.29 47.29 0 0 0 33.74-15.35 350 350 0 0 0 94.26-239.33 344.33 344.33 0 0 0-24.06-128.9z" fill="#FFB940"></path>
                <path d="M883.08 164.31a46.2 46.2 0 0 0-33.55-14.83 49.08 49.08 0 0 0-34.14 13l-0.35 0.33a48.56 48.56 0 0 0-0.15 69A404 404 0 0 1 927.37 512a399.6 399.6 0 0 1-113.72 280.2l-0.65 0.8a53.74 53.74 0 0 0-12.4 32.74v0.16a46.27 46.27 0 0 0 16.15 35.79 49.18 49.18 0 0 0 33.14 12.88 46.89 46.89 0 0 0 34.41-14.88C974.38 764.7 1024 641.25 1024 512c0-129.72-50-253.19-140.92-347.69zM541.93 74.73c-14.22 0.14-27 3.1-38.44 10.25l-2.93 1.46c-95 69.33-244.3 176.6-294 211.24a36.38 36.38 0 0 1-20.82 6.54L70 304.31a70.28 70.28 0 0 0-70 70.13l0.21 275.94a70.09 70.09 0 0 0 70.13 70l115.66-0.06a36.45 36.45 0 0 1 20.82 6.51c49.81 34.56 199.27 141.61 294.38 210.79l2.92 1.46a79.49 79.49 0 0 0 39.43 10.19c20.44 0 39.41-7.33 55.46-20.48 14.59-16.07 23.34-35.06 23.32-57l-0.54-717.64c-0.03-43.85-35.97-79.84-79.86-79.42z" fill="#FFB940"></path>
              </svg>
            </button>
            <span v-if="message.meta && message.meta.status === 'streaming'" class="streaming-indicator"> ··</span>
          </div>
          <div class="message-content" v-html="renderMarkdown(message.content)"></div>
        </div>
      </div>

      <div class="chat-input">
        <!-- 文件预览区域 -->
        <div v-if="uploadedFiles.length > 0" class="file-preview-area">
          <div v-for="(file, index) in uploadedFiles" :key="index" class="file-card">
            <span class="file-icon">📄</span>
            <span class="file-name">{{ file.name }}</span>
            <button class="remove-file-btn" @click="removeFile(index)">×</button>
          </div>
        </div>

        <textarea
          v-model="inputMessage"
          placeholder="请输入你的问题..."
          @keydown.enter.exact.prevent="sendMessage"
          :disabled="loading"
          ref="messageInput"
          rows="1"
        ></textarea>
        <button
          type="button"
          :disabled="!inputMessage.trim() || loading"
          @click="sendMessage"
          class="send-btn"
        >
          {{ loading ? '发送中...' : '发送' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, nextTick, computed, onMounted } from 'vue'
  // import { useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
import api from '../utils/api'
import { marked } from 'marked'
import ThemeToggle from '../components/ThemeToggle.vue'

export default {
  name: 'AIChat',
  components: { ThemeToggle },
  setup() {
    // const router = useRouter() // unused

    const sessions = ref({})
    const currentSessionId = ref(null)
    const tempSession = ref(false)
    const currentMessages = ref([])
    const inputMessage = ref('')
    const loading = ref(false)
    const messagesRef = ref(null)
    const messageInput = ref(null)
    const selectedModel = ref('1')
    const uploading = ref(false)
    const fileInput = ref(null)
    const uploadedFiles = ref([])

    const escapeHtml = (s) =>
      String(s)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;')

    const mdRenderer = new marked.Renderer()
    mdRenderer.html = () => ''
    mdRenderer.link = (href, title, text) => {
      const h = escapeHtml(href || '')
      const t = escapeHtml(text || '')
      const ttl = title ? ` title="${escapeHtml(title)}"` : ''
      return `<a href="${h}" target="_blank" rel="noopener noreferrer"${ttl}>${t}</a>`
    }
    marked.setOptions({
      renderer: mdRenderer,
      gfm: true,
      breaks: true,
      headerIds: false,
      mangle: false
    })


    const renderMarkdown = (text) => {
      if (!text && text !== '') return ''
      const src = String(text)
      const html = marked.parse(src)
      if (typeof html === 'string') {
        const wrapped = html.replace(/<p>搜索结果:<\/p>\s*<ul>([\s\S]*?)<\/ul>/, '<div class="search-results"><div class="search-title">搜索结果</div><ul>$1</ul></div>')
        return wrapped
      }
      return `<pre class="plain-text">${escapeHtml(src)}</pre>`
    }

    const playTTS = async (text) => {
      try {
        // 创建TTS任务
        const createResponse = await api.post('/AI/chat/tts', { text })
        if (createResponse.data && createResponse.data.status_code === 1000 && createResponse.data.task_id) {
          const taskId = createResponse.data.task_id
          
          // 先等待5秒钟再开始轮询
          await new Promise(resolve => setTimeout(resolve, 5000))
          
          // 轮询查询任务结果
          const maxAttempts = 30
          const pollInterval = 2000
          let attempts = 0
          
          const pollResult = async () => {
            const queryResponse = await api.get('/AI/chat/tts/query', { params: { task_id: taskId } })
            
            if (queryResponse.data && queryResponse.data.status_code === 1000) {
              const taskStatus = queryResponse.data.task_status
                
              if (taskStatus === 'Success' && queryResponse.data.task_result) {
                // 任务完成，播放音频
                // 后端返回的 task_result 是直接的 URL 字符串
                const audio = new Audio(queryResponse.data.task_result)
                audio.play()
                return true
              } else if (taskStatus === 'Running' ||taskStatus === 'Created' ) {
                // 任务进行中，继续轮询
                attempts++
                if (attempts < maxAttempts) {
                  await new Promise(resolve => setTimeout(resolve, pollInterval))
                  return await pollResult()
                } else {
                  ElMessage.error('语音合成超时')
                  return true
                }
              } else {
                // 其他状态（如失败）
                ElMessage.error('语音合成失败')
                return true
              }
            }
            
            attempts++
            if (attempts < maxAttempts) {
              await new Promise(resolve => setTimeout(resolve, pollInterval))
              return await pollResult()
            } else {
              ElMessage.error('语音合成超时')
              return true
            }
          }
          
          await pollResult()
        } else {
          ElMessage.error('无法创建语音合成任务')
        }
      } catch (error) {
        console.error('TTS error:', error)
        ElMessage.error('请求语音接口失败')
      }
    }

    const loadSessions = async () => {
      try {
        const response = await api.get('/AI/chat/sessions')
        if (response.data && response.data.status_code === 1000 && Array.isArray(response.data.sessions)) {
          const sessionMap = {}
          response.data.sessions.forEach(s => {
            const sid = String(s.sessionId)
            sessionMap[sid] = {
              id: sid,
              name: s.name || `会话 ${sid}`,
              messages: [] // lazy load
            }
          })
          sessions.value = sessionMap
        }
      } catch (error) {
        console.error('Load sessions error:', error)
      }
    }

    const createNewSession = () => {
      currentSessionId.value = 'temp'
      tempSession.value = true
      currentMessages.value = []
      // focus input
      nextTick(() => {
        if (messageInput.value) messageInput.value.focus()
      })
    }

    const switchSession = async (sessionId) => {
      if (!sessionId) return
      currentSessionId.value = String(sessionId)
      tempSession.value = false

      // lazy load history if not present
      if (!sessions.value[sessionId].messages || sessions.value[sessionId].messages.length === 0) {
        try {
          const response = await api.post('/AI/chat/history', { sessionId: currentSessionId.value })
          if (response.data && response.data.status_code === 1000 && Array.isArray(response.data.history)) {
            const messages = response.data.history.map(item => ({
              role: item.is_user ? 'user' : 'assistant',
              content: item.content
            }))
            sessions.value[sessionId].messages = messages
          }
        } catch (err) {
          console.error('Load history error:', err)
        }
      }


      currentMessages.value = [...(sessions.value[sessionId].messages || [])]
      await nextTick()
      scrollToBottom()
    }

    const syncHistory = async () => {
      if (!currentSessionId.value || tempSession.value) {
        ElMessage.warning('请选择已有会话进行同步')
        return
      }
      try {
        const response = await api.post('/AI/chat/history', { sessionId: currentSessionId.value })
        if (response.data && response.data.status_code === 1000 && Array.isArray(response.data.history)) {
          const messages = response.data.history.map(item => ({
            role: item.is_user ? 'user' : 'assistant',
            content: item.content
          }))
          sessions.value[currentSessionId.value].messages = messages
          currentMessages.value = [...messages]
          await nextTick()
          scrollToBottom()
        } else {
          ElMessage.error('无法获取历史数据')
        }
      } catch (err) {
        console.error('Sync history error:', err)
        ElMessage.error('请求历史数据失败')
      }
    }

    const deleteSession = async (sessionId) => {
      if (!sessionId) return
      const confirmed = window.confirm('确定要删除该会话吗？')
      if (!confirmed) return
      try {
        const response = await api.post('/AI/chat/session/delete', { sessionId })
        if (response.data && response.data.status_code === 1000) {
          const newSessions = { ...sessions.value }
          delete newSessions[sessionId]
          sessions.value = newSessions
          if (currentSessionId.value === sessionId) {
            currentSessionId.value = null
            currentMessages.value = []
            tempSession.value = false
          }
          ElMessage.success('会话已删除')
        } else {
          ElMessage.error('删除会话失败')
        }
      } catch (err) {
        console.error('Delete session error:', err)
        ElMessage.error('删除会话失败')
      }
    }


    const sendMessage = async () => {
      if (!inputMessage.value || !inputMessage.value.trim()) {
        ElMessage.warning('请输入消息内容')
        return
      }

      // 如果当前没有选中会话，且不是临时会话状态，则自动开启新会话
      if (!currentSessionId.value && !tempSession.value) {
        tempSession.value = true
        currentSessionId.value = 'temp'
      }

      const userMessage = {
        role: 'user',
        content: inputMessage.value
      }
      const currentInput = inputMessage.value
      inputMessage.value = ''


      currentMessages.value.push(userMessage)
      await nextTick()
      scrollToBottom()

      try {
        loading.value = true
        // Always use streaming
        await handleStreaming(currentInput)
      } catch (err) {
        console.error('Send message error:', err)
        ElMessage.error('发送失败，请重试')

        if (!tempSession.value && currentSessionId.value && sessions.value[currentSessionId.value] && sessions.value[currentSessionId.value].messages) {

          const sessionArr = sessions.value[currentSessionId.value].messages
          if (sessionArr && sessionArr.length) sessionArr.pop()
        }
        currentMessages.value.pop()
      } finally {
        // Always finish loading when done (streaming handles its own loading state during stream, but here we ensure it's off if error occurs)
        // Wait, handleStreaming sets loading=false when done.
        // But if error occurs in handleStreaming, it sets loading=false.
        // If error occurs before handleStreaming (unlikely), we should set it false.
        // Since we are not using isStreaming flag anymore, we can't check it.
        // But since we await handleStreaming, and handleStreaming manages loading.value internally for success/error cases,
        // we might not need to set it here unless we want to be safe.
        // However, if we set it here, it might overwrite the state if handleStreaming is still running (but we await it, so it's finished).
        // So it is safe to set loading.value = false here just in case.
        loading.value = false
        // 发送完成后清空上传列表
        uploadedFiles.value = []
        await nextTick()
        scrollToBottom()
      }
    }


    async function handleStreaming(question) {

      const aiMessage = {
        role: 'assistant',
        content: '',
        meta: { status: 'streaming' } // mark streaming
      }


      const aiMessageIndex = currentMessages.value.length
      currentMessages.value.push(aiMessage)

      if (!tempSession.value && currentSessionId.value && sessions.value[currentSessionId.value]) {
        if (!sessions.value[currentSessionId.value].messages) sessions.value[currentSessionId.value].messages = []
        sessions.value[currentSessionId.value].messages.push({ role: 'assistant', content: '' })
      }


      const isDev = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
      const backendBase = isDev
        ? `http://${window.location.hostname}:9090/api/v1/AI`
        : '/api/AI'

      const url = tempSession.value
        ? `${backendBase}/chat/send-stream-new-session`
        : `${backendBase}/chat/send-stream`

      const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
      }

      const body = tempSession.value
        ? { question: question, modelType: selectedModel.value }
        : { question: question, modelType: selectedModel.value, sessionId: currentSessionId.value }

      try {
        // 创建 fetch 连接读取 SSE 流
        const response = await fetch(url, {
          method: 'POST',
          headers,
          body: JSON.stringify(body)
        })

        if (!response.ok) {
          loading.value = false
          throw new Error('Network response was not ok')
        }

        const reader = response.body.getReader()
        const decoder = new TextDecoder()
        let buffer = ''

        // 读取流数据
        // eslint-disable-next-line no-constant-condition
        while (true) {
          const { done, value } = await reader.read()
          if (done) break

          const chunk = decoder.decode(value, { stream: true })
          buffer += chunk

          // 按行分割
          const lines = buffer.split('\n')
          buffer = lines.pop() || '' // 保留未完成的行

          for (const line of lines) {
            const normalizedLine = line.endsWith('\r') ? line.slice(0, -1) : line
            if (!normalizedLine) continue

            // 处理 SSE 格式：data: <content>
            if (normalizedLine.startsWith('data:')) {
              let data = normalizedLine.slice(5)
              if (data.startsWith(' ')) data = data.slice(1)
              console.log('[SSE] Received:', data) // 调试日志

              if (data === '[DONE]') {
                // 流结束
                console.log('[SSE] Stream done')
                loading.value = false
                currentMessages.value[aiMessageIndex].meta = { status: 'done' }
                currentMessages.value = [...currentMessages.value]
              } else if (data.startsWith('{')) {
                // 尝试解析 JSON（如 sessionId）
                try {
                  const parsed = JSON.parse(data)
                  if (parsed.sessionId) {
                    const newSid = String(parsed.sessionId)
                    console.log('[SSE] Session ID:', newSid)
                    if (tempSession.value) {
                      let title = (question || '').trim()
                      if (!title) {
                        title = `会话 ${newSid}`
                      } else if (title.length > 30) {
                        title = `${title.slice(0, 30)}...`
                      }
                      sessions.value[newSid] = {
                        id: newSid,
                        name: title,
                        messages: [...currentMessages.value]
                      }
                      currentSessionId.value = newSid
                      tempSession.value = false
                    }
                  } else if (parsed.type === 'delta' && typeof parsed.content === 'string') {
                    currentMessages.value[aiMessageIndex].content += parsed.content
                  }
                } catch (e) {
                  // 不是 JSON，当作普通文本处理
                  currentMessages.value[aiMessageIndex].content += data
                  console.log('[SSE] Content updated:', currentMessages.value[aiMessageIndex].content.length)
                }
              } else {
                // 普通文本数据，直接追加
                // 使用数组索引直接更新，强制 Vue 响应式系统检测变化
                currentMessages.value[aiMessageIndex].content += data
                console.log('[SSE] Content updated:', currentMessages.value[aiMessageIndex].content.length)
              }

              // 每收到一条数据就立即更新 DOM
              // 强制更新整个数组以触发响应式
              currentMessages.value = [...currentMessages.value]
              
              // 使用 requestAnimationFrame 强制浏览器重排
              await new Promise(resolve => {
                requestAnimationFrame(() => {
                  scrollToBottom()
                  resolve()
                })
              })
            }
          }
        }

        // 流读取完成后的处理
        loading.value = false
        currentMessages.value[aiMessageIndex].meta = { status: 'done' }
        currentMessages.value = [...currentMessages.value]

        // 同步到 sessions 存储
        if (!tempSession.value && currentSessionId.value && sessions.value[currentSessionId.value]) {
          const sessMsgs = sessions.value[currentSessionId.value].messages
          if (Array.isArray(sessMsgs) && sessMsgs.length) {
            const lastIndex = sessMsgs.length - 1
            if (sessMsgs[lastIndex] && sessMsgs[lastIndex].role === 'assistant') {
              sessMsgs[lastIndex].content = currentMessages.value[aiMessageIndex].content
            }
          }
        }
      } catch (err) {
        console.error('Stream error:', err)
        loading.value = false
        currentMessages.value[aiMessageIndex].meta = { status: 'error' }
        currentMessages.value = [...currentMessages.value]
        ElMessage.error('流式传输出错')
      }
    }


    const scrollToBottom = () => {
      if (messagesRef.value) {
        try {
          messagesRef.value.scrollTop = messagesRef.value.scrollHeight
        } catch (e) {
          // ignore
        }
      }
    }

    const triggerFileUpload = () => {
      if (fileInput.value) {
        fileInput.value.click()
      }
    }

    const handleFileUpload = async (event) => {
      const file = event.target.files[0]
      if (!file) return

      // 前端校验：只允许.md或.txt文件
      const fileName = file.name.toLowerCase()
      if (!fileName.endsWith('.md') && !fileName.endsWith('.txt')) {
        ElMessage.error('只允许上传 .md 或 .txt 文件')
        // 清空文件输入
        if (fileInput.value) {
          fileInput.value.value = ''
        }
        return
      }

      try {
        uploading.value = true

        const formData = new FormData()
        formData.append('file', file)

        const response = await api.post('/file/upload', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })

        if (response.data && response.data.status_code === 1000) {
          ElMessage.success(`文件上传成功`)
          // 添加到预览列表
          uploadedFiles.value.push({
            name: file.name,
            type: file.type
          })
        } else {
          ElMessage.error(response.data?.status_msg || '上传失败')
        }
      } catch (error) {
        console.error('File upload error:', error)
        ElMessage.error('文件上传失败')
      } finally {
        uploading.value = false
        // 清空文件输入
        if (fileInput.value) {
          fileInput.value.value = ''
        }
      }
    }

    const removeFile = (index) => {
      uploadedFiles.value.splice(index, 1)
    }

    onMounted(() => {
      loadSessions()
    })

    // expose to template
    return {
      sessions: computed(() => Object.values(sessions.value)),
      currentSessionId,
      tempSession,
      currentMessages,
      inputMessage,
      loading,
      messagesRef,
      messageInput,
      selectedModel,
      // isStreaming, // Removed
      uploading,
      fileInput,
      renderMarkdown,
      playTTS,
      createNewSession,
      switchSession,
      syncHistory,
      deleteSession,
      sendMessage,
      triggerFileUpload,
      handleFileUpload,
      uploadedFiles,
      removeFile
    }
  }
}
</script>

<style scoped>
.ai-chat-container {
  height: 100vh;
  display: flex;
  background-color: var(--bg-cyber);
  background-image: 
    linear-gradient(var(--cyber-grid) 1px, transparent 1px),
    linear-gradient(90deg, var(--cyber-grid) 1px, transparent 1px);
  background-size: 40px 40px;
  position: relative;
  overflow: hidden;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  color: var(--text-primary);
  padding: 20px;
  gap: 20px;
  box-sizing: border-box;
  transition: background-color 0.3s;
}

/* 扫描波浪效果 */
.ai-chat-container::after {
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

/* 动态背景光晕 */
.ai-chat-container::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle at 50% 50%, var(--cyber-pulse) 0%, transparent 60%);
  animation: bgPulse 20s ease-in-out infinite;
  pointer-events: none;
  z-index: 0;
}

@keyframes bgPulse {
  0%, 100% { transform: scale(1); opacity: 0.4; }
  50% { transform: scale(1.1); opacity: 0.6; }
}

/* 进场动画定义 */
@keyframes slideInLeft {
  from { opacity: 0; transform: translateX(-80px); }
  to { opacity: 1; transform: translateX(0); }
}

@keyframes slideInDown {
  from { opacity: 0; transform: translateY(-80px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes slideInUp {
  from { opacity: 0; transform: translateY(80px); }
  to { opacity: 1; transform: translateY(0); }
}

.session-list {
  width: 280px;
  height: 100%; /* 填满父容器 */
  overflow: hidden;
  display: flex;
  flex-direction: column;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: 24px;
  position: relative;
  z-index: 2;
  box-shadow: 0 10px 30px var(--shadow-color);
  animation: slideInLeft 1.0s ease-out;
}

.session-list-header {
  padding: 24px;
  text-align: center;
  font-weight: 600;
  background: transparent;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: center;
  color: #f1c40f;
}

.new-chat-btn {
  width: 100%;
  padding: 12px 0;
  cursor: pointer;
  background: #f1c40f;
  color: #000;
  border: none;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 700;
  transition: all 0.3s ease;
  text-transform: uppercase;
  letter-spacing: 1px;
  box-shadow: 0 4px 12px rgba(241, 196, 15, 0.2);
}

.new-chat-btn:hover {
  background: #f39c12;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(241, 196, 15, 0.4);
}

.session-list-ul {
  list-style: none;
  padding: 10px;
  margin: 0;
  flex: 1;
  overflow-y: auto;
}

.session-item {
  padding: 16px 20px;
  margin-bottom: 8px;
  cursor: pointer;
  border: 1px solid transparent;
  border-radius: 16px;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  position: relative;
  color: #aaa;
  background: rgba(255, 255, 255, 0.02);
  overflow: hidden;
}

/* 充能悬浮特效 */
.session-item:hover {
  background: rgba(241, 196, 15, 0.1);
  color: #fff;
  border-color: #f1c40f;
  box-shadow: 0 0 15px rgba(241, 196, 15, 0.2), inset 0 0 10px rgba(241, 196, 15, 0.1);
  transform: translateX(4px);
}

.session-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: #f1c40f;
  opacity: 0;
  transition: opacity 0.3s;
}

.session-item:hover::before {
  opacity: 1;
}

.session-item.active {
  background: rgba(241, 196, 15, 0.15);
  color: #f1c40f;
  font-weight: 600;
  border-color: rgba(241, 196, 15, 0.5);
  box-shadow: 0 0 20px rgba(241, 196, 15, 0.15);
}

.session-item.active::before {
  opacity: 1;
  box-shadow: 0 0 10px #f1c40f;
}

.session-title {
  display: inline-block;
  max-width: 200px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.session-delete-btn {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: none;
  background: rgba(0, 0, 0, 0.6);
  color: #f1c40f;
  font-size: 14px;
  line-height: 20px;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s ease, background 0.2s ease, transform 0.2s ease;
}

.session-item:hover .session-delete-btn {
  opacity: 1;
}

.session-delete-btn:hover {
  background: #f1c40f;
  color: #000;
  transform: translateY(-50%) scale(1.05);
}

/* chat section */
.chat-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 1;
  min-width: 0;
  min-height: 0;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 10px 30px var(--shadow-color);
  animation: slideInDown 1.0s ease-out;
}

.top-bar {
  background: var(--header-bg);
  color: var(--text-primary);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid var(--border-color);
  gap: 12px;
}

.top-bar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.back-btn {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #ddd;
  padding: 8px 14px;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s ease;
}

.back-btn:hover {
  border-color: #f1c40f;
  color: #f1c40f;
  background: rgba(241, 196, 15, 0.1);
}

.sync-btn {
  background: transparent;
  color: #f1c40f;
  padding: 8px 14px;
  border: 1px solid rgba(241, 196, 15, 0.5);
  border-radius: 12px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.sync-btn:hover:not(:disabled) {
  background: #f1c40f;
  color: #000;
  box-shadow: 0 0 15px rgba(241, 196, 15, 0.3);
}

.sync-btn:disabled {
  border-color: #444;
  color: #444;
  cursor: not-allowed;
}

.model-select-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: 10px;
}

.model-select-wrapper .label {
  font-size: 14px;
  font-weight: 600;
  color: #ccc;
}

.custom-select {
  width: 160px;
}

:deep(.el-input__wrapper),
:deep(.el-select__wrapper) {
  background-color: rgba(0, 0, 0, 0.3) !important;
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.1) inset !important;
  border-radius: 8px !important;
}

:deep(.el-input__wrapper.is-focus),
:deep(.el-select__wrapper.is-focused) {
  box-shadow: 0 0 0 1px #f1c40f inset !important;
  background-color: rgba(0, 0, 0, 0.3) !important;
}

:deep(.el-input__inner),
:deep(.el-select__selected-item) {
  color: #fff !important;
  font-weight: 600;
}

/* 修复选中时的蓝色字体和背景 */
:deep(.el-select .el-input.is-focus .el-input__inner),
:deep(.el-select__wrapper.is-focused .el-select__selected-item) {
  color: #f1c40f !important;
}

.upload-btn {
  background: transparent;
  color: #f1c40f;
  padding: 8px 14px;
  border: 1px solid rgba(241, 196, 15, 0.5);
  border-radius: 12px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.upload-btn:hover:not(:disabled) {
  background: #f1c40f;
  color: #000;
  box-shadow: 0 0 15px rgba(241, 196, 15, 0.3);
}

.chat-messages {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 30px;
  padding-bottom: 120px; /* 为悬浮输入框留出空间 */
  display: flex;
  flex-direction: column;
  gap: 20px;
  background: transparent;
}

.message {
  max-width: 70%;
  padding: 16px 20px;
  border-radius: 20px;
  line-height: 1.6;
  word-wrap: break-word;
  position: relative;
  font-size: 15px;
  box-sizing: border-box;
  box-shadow: 0 4px 10px var(--shadow-color);
}

.user-message {
  align-self: flex-end;
  background: var(--accent-color);
  color: var(--text-inverse);
  border-bottom-right-radius: 4px; /* 气泡角风格 */
  box-shadow: 0 5px 20px var(--cyber-wave);
}

.user-message::after {
  /* 移除以前的三角箭头，改用圆角风格 */
  content: none;
}

.ai-message {
  align-self: flex-start;
  background: var(--bg-input);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
  border-bottom-left-radius: 4px;
}

.ai-message::after {
  content: none;
}

/* Markdown 样式适配 */
:deep(.message-content) {
  line-height: 1.6;
  font-size: 15px;
  white-space: normal;
}

:deep(.message-content pre.plain-text) {
  margin: 0;
  white-space: pre-wrap;
  background: transparent;
  border: 0;
  padding: 0;
  font: inherit;
  color: inherit;
}

:deep(.message-content p) {
  margin: 0 0 10px 0;
}

:deep(.message-content p:last-child) {
  margin-bottom: 0;
}

:deep(.message-content pre) {
  background: #1e1e1e; /* 深色背景 */
  padding: 12px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 10px 0;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

:deep(.message-content code) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  background: rgba(241, 196, 15, 0.15); /* 黄色背景微调 */
  color: #f1c40f;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.9em;
}

:deep(.message-content pre code) {
  background: transparent;
  color: #dcdcdc; /* 代码块内颜色 */
  padding: 0;
  border-radius: 0;
  font-size: 14px;
}

:deep(.message-content ul), :deep(.message-content ol) {
  margin: 0 0 10px 20px;
  padding: 0;
}

:deep(.message-content li) {
  margin-bottom: 4px;
}

:deep(.message-content h1), :deep(.message-content h2), :deep(.message-content h3) {
  margin: 15px 0 10px 0;
  font-weight: 600;
  color: #f1c40f;
}

:deep(.message-content blockquote) {
  margin: 10px 0;
  padding-left: 15px;
  border-left: 4px solid #f1c40f;
  color: #999;
}

:deep(.message-content a) {
  color: #f1c40f;
  text-decoration: none;
  border-bottom: 1px dashed #f1c40f;
}

:deep(.message-content a:hover) {
  border-bottom-style: solid;
}

.search-results {
  background: rgba(241, 196, 15, 0.08);
  border: 1px solid rgba(241, 196, 15, 0.5);
  border-radius: 12px;
  padding: 12px 16px;
  margin: 12px 0;
}

.search-results .search-title {
  font-weight: 700;
  color: #f1c40f;
  margin-bottom: 8px;
}

:deep(.ai-message .message-header b) {
  color: #f1c40f;
}

:deep(.ai-message .message-content b),
:deep(.ai-message .message-content strong) {
  color: #f1c40f;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  font-size: 12px;
  opacity: 0.8;
}

.tts-btn {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.1);
  color: #f1c40f;
  border: none;
}

.tts-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

/* input area */
.chat-input-container {
  position: absolute;
  bottom: 20px;
  left: 320px;
  right: 20px;
  z-index: 10;
  animation: slideInUp 1.0s ease-out;
}

.chat-input {
  padding: 20px;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-color);
  border-radius: 24px;
  position: relative;
  box-shadow: 0 10px 30px var(--shadow-color);
}

.chat-input textarea {
  width: 100%;
  resize: none;
  border: 1px solid var(--border-color);
  border-radius: 16px;
  padding: 16px 18px;
  font-size: 15px;
  outline: none;
  background: var(--bg-input);
  color: var(--text-primary);
  transition: all 0.2s ease;
  min-height: 20px;
  max-height: 160px;
}

.chat-input textarea:focus {
  border-color: var(--accent-color);
  background: var(--bg-input);
  box-shadow: 0 0 15px var(--cyber-pulse);
}

.send-btn {
  position: absolute;
  right: 36px;
  bottom: 34px;
  padding: 10px 24px;
  border: none;
  border-radius: 50px;
  background: var(--accent-color);
  color: var(--text-inverse);
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.send-btn:hover:not(:disabled) {
  background: var(--accent-hover);
  transform: translateY(-2px);
  box-shadow: 0 4px 15px var(--cyber-wave);
}

.send-btn:disabled {
  background: #444;
  color: #888;
  cursor: not-allowed;
}

/* File Preview Area */
.file-preview-area {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 10px 18px 0;
}

.file-card {
  display: flex;
  align-items: center;
  background: rgba(241, 196, 15, 0.1);
  border: 1px solid rgba(241, 196, 15, 0.3);
  border-radius: 8px;
  padding: 6px 10px;
  gap: 8px;
  font-size: 13px;
  color: #f1c40f;
  animation: slideInUp 0.3s ease-out;
}

.file-icon {
  font-size: 16px;
}

.file-name {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.remove-file-btn {
  background: transparent;
  border: none;
  color: #f1c40f;
  cursor: pointer;
  font-size: 16px;
  line-height: 1;
  padding: 0 4px;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.remove-file-btn:hover {
  opacity: 1;
}
</style>

<style>
/* 全局覆盖下拉菜单样式 */
.custom-dropdown.el-popper {
  background-color: var(--bg-secondary) !important;
  border: 1px solid var(--border-color) !important;
  backdrop-filter: blur(10px);
}

.custom-dropdown .el-select-dropdown__item {
  color: var(--text-regular) !important;
  background-color: transparent !important;
}

.custom-dropdown .el-select-dropdown__item.hover,
.custom-dropdown .el-select-dropdown__item:hover {
  background-color: var(--cyber-grid) !important;
  color: var(--accent-color) !important;
}

.custom-dropdown .el-select-dropdown__item.is-selected {
  color: var(--accent-color) !important;
  font-weight: 700;
  background-color: var(--cyber-wave) !important;
}

/* 修复小箭头颜色 */
.custom-dropdown .el-popper__arrow::before {
  background-color: var(--bg-secondary) !important;
  border: 1px solid var(--border-color) !important;
}
</style>
