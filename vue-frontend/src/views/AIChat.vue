<template>
  <div class="ai-chat-container">
    <!-- Glass Panel Wrapper -->
    <div class="glass-wrapper">
      <!-- 左侧会话列表 -->
      <div class="session-list">
        <div class="session-list-header">
          <span>
            <LogoEye :size="24" class="header-logo" />
            会话列表
          </span>
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
        <div class="chat-header">
          <div class="header-left">
            <button class="icon-btn back-btn" @click="$router.push('/menu')">
              <el-icon><ArrowLeft /></el-icon>
            </button>
            <div class="header-info">
              <h2>AI Assistant</h2>
              <span class="status-dot"></span>
            </div>
          </div>
          
          <div class="header-right">
            <el-select v-model="selectedModel" class="model-select" placeholder="Select Model">
              <el-option label="✨ Qwen Max" value="1" />
              <el-option label="🧠 Qwen RAG" value="2" />
              <el-option label="🔌 Qwen MCP" value="3" />
            </el-select>
            <button class="icon-btn" @click="syncHistory" title="Sync History">
              <el-icon><Refresh /></el-icon>
            </button>
          </div>
        </div>

        <div class="chat-messages" ref="messagesRef">
          <!-- Empty State -->
          <div v-if="currentMessages.length === 0" class="empty-state">
            <div class="empty-icon">
              <LogoEye :size="60" />
            </div>
            <h3>How can I help you today?</h3>
            <div class="starter-pills">
              <button class="pill" @click="setInputAndSend('Explain quantum physics like I\'m 5')">⚛️ Explain Physics</button>
              <button class="pill" @click="setInputAndSend('Write a python script to parse CSV')">🐍 Python Script</button>
              <button class="pill" @click="setInputAndSend('Tell me a fun fact about space')">🌌 Space Fact</button>
            </div>
          </div>

          <div
            v-for="(message, index) in currentMessages"
            :key="index"
            :class="['message-row', message.role === 'user' ? 'user-row' : 'ai-row']"
          >
            <div class="avatar">
              <span v-if="message.role === 'user'">U</span>
              <LogoEye v-else :size="24" />
            </div>
            <div class="message-bubble">
              <div class="bubble-content" v-html="renderMarkdown(message.content)"></div>
              <div v-if="message.role === 'assistant'" class="message-actions">
                <button class="action-btn" @click="playTTS(message.content)" title="Read Aloud">
                  <el-icon><Microphone /></el-icon>
                </button>
              </div>
            </div>
          </div>
          
          <div v-if="loading" class="message-row ai-row">
            <div class="avatar">
              <LogoEye :size="24" />
            </div>
            <div class="message-bubble loading-bubble">
              <div class="typing-dots">
                <span></span><span></span><span></span>
              </div>
            </div>
          </div>
        </div>
        
        <div class="chat-input-wrapper">
          <!-- File Preview -->
          <div v-if="uploadedFiles.length > 0" class="file-preview-bar">
            <div v-for="(file, index) in uploadedFiles" :key="index" class="file-chip">
              <el-icon><Document /></el-icon>
              <span class="filename">{{ file.name }}</span>
              <button class="close-file" @click="removeFile(index)">×</button>
            </div>
          </div>

          <div class="input-container">
            <button class="attach-btn" @click="triggerFileUpload" :disabled="uploading">
              <el-icon><Paperclip /></el-icon>
            </button>
            <input
              ref="fileInput"
              type="file"
              accept=".md,.txt,text/markdown,text/plain"
              style="display: none"
              @change="handleFileUpload"
            />
            <textarea
              v-model="inputMessage"
              placeholder="Type a message..."
              @keydown.enter.exact.prevent="sendMessage"
              :disabled="loading"
              ref="messageInput"
              rows="1"
              class="main-textarea"
            ></textarea>
            <button
              class="send-btn-round"
              :disabled="!inputMessage.trim() || loading"
              @click="sendMessage"
            >
              <el-icon><Position /></el-icon>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, nextTick, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Refresh, Microphone, Paperclip, Position, Document } from '@element-plus/icons-vue'
import api from '../utils/api'
import { marked } from 'marked'
import LogoEye from '../components/LogoEye.vue'

export default {
  name: 'AIChat',
  components: {
    ArrowLeft, Refresh, Microphone, Paperclip, Position, Document, LogoEye
  },
  setup() {
    const sessions = ref({})
    // ... (keep existing state)
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

    // Helper for starters
    const setInputAndSend = (text) => {
      inputMessage.value = text
      sendMessage()
    }

    // ... (keep existing functions: escapeHtml, mdRenderer, renderMarkdown, playTTS, loadSessions, etc.)
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
        return html
      }
      return `<pre class="plain-text">${escapeHtml(src)}</pre>`
    }

    const playTTS = async (text) => {
      // ... (keep existing implementation)
      try {
        const createResponse = await api.post('/AI/chat/tts', { text })
        if (createResponse.data && createResponse.data.status_code === 1000 && createResponse.data.task_id) {
          const taskId = createResponse.data.task_id
          // Shorten wait for better UX, maybe 2s
          await new Promise(resolve => setTimeout(resolve, 2000))
          
          const maxAttempts = 30
          const pollInterval = 2000
          let attempts = 0
          
          const pollResult = async () => {
            const queryResponse = await api.get('/AI/chat/tts/query', { params: { task_id: taskId } })
            
            if (queryResponse.data && queryResponse.data.status_code === 1000) {
              const taskStatus = queryResponse.data.task_status
                
              if (taskStatus === 'Success' && queryResponse.data.task_result) {
                let audioUrl = queryResponse.data.task_result
                // 强制使用 https，避免 Mixed Content 拦截
                if (audioUrl.startsWith('http://')) {
                  audioUrl = audioUrl.replace('http://', 'https://')
                }
                const audio = new Audio(audioUrl)
                audio.play().catch(err => {
                  console.error('Audio play error:', err)
                  ElMessage.error('播放失败: 浏览器可能拦截了自动播放，或地址无法访问')
                })
                return true
              } else if (taskStatus === 'Running' ||taskStatus === 'Created' ) {
                attempts++
                if (attempts < maxAttempts) {
                  await new Promise(resolve => setTimeout(resolve, pollInterval))
                  return await pollResult()
                } else {
                  ElMessage.error('TTS Timeout')
                  return true
                }
              } else {
                ElMessage.error('TTS Failed')
                return true
              }
            }
            attempts++
            if (attempts < maxAttempts) {
              await new Promise(resolve => setTimeout(resolve, pollInterval))
              return await pollResult()
            }
          }
          await pollResult()
        } else {
          ElMessage.error('Could not create TTS task')
        }
      } catch (error) {
        console.error('TTS error:', error)
      }
    }

    const loadSessions = async () => {
        // ... (keep existing)
      try {
        const response = await api.get('/AI/chat/sessions')
        if (response.data && response.data.status_code === 1000 && Array.isArray(response.data.sessions)) {
          const sessionMap = {}
          response.data.sessions.forEach(s => {
            const sid = String(s.sessionId)
            sessionMap[sid] = {
              id: sid,
              name: s.name || `Session ${sid}`,
              messages: []
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
      nextTick(() => {
        if (messageInput.value) messageInput.value.focus()
      })
    }

    const switchSession = async (sessionId) => {
      // ... (keep existing)
      if (!sessionId) return
      currentSessionId.value = String(sessionId)
      tempSession.value = false

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
      // ... (keep existing)
      if (!currentSessionId.value || tempSession.value) return
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
        }
      } catch (err) {
        console.error('Sync history error:', err)
      }
    }

    const deleteSession = async (sessionId) => {
        // ... (keep existing)
      if (!sessionId) return
      const confirmed = window.confirm('Delete this session?')
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
          ElMessage.success('Session deleted')
        }
      } catch (err) {
        console.error('Delete session error:', err)
      }
    }

    const sendMessage = async () => {
      // ... (keep existing logic mostly)
      if (!inputMessage.value || !inputMessage.value.trim()) return

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
        await handleStreaming(currentInput)
      } catch (err) {
        console.error('Send message error:', err)
        ElMessage.error('Failed to send')
        currentMessages.value.pop()
      } finally {
        loading.value = false
        uploadedFiles.value = []
        await nextTick()
        scrollToBottom()
      }
    }

    async function handleStreaming(question) {
        // ... (keep existing streaming logic)
      const aiMessage = {
        role: 'assistant',
        content: '',
        meta: { status: 'streaming' } 
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
        const response = await fetch(url, { method: 'POST', headers, body: JSON.stringify(body) })
        if (!response.ok) throw new Error('Network response was not ok')

        const reader = response.body.getReader()
        const decoder = new TextDecoder()
        let buffer = ''

        // eslint-disable-next-line no-constant-condition
        while (true) {
          const { done, value } = await reader.read()
          if (done) break
          const chunk = decoder.decode(value, { stream: true })
          buffer += chunk
          const lines = buffer.split('\n')
          buffer = lines.pop() || '' 

          for (const line of lines) {
            const normalizedLine = line.endsWith('\r') ? line.slice(0, -1) : line
            if (!normalizedLine) continue
            if (normalizedLine.startsWith('data:')) {
              let data = normalizedLine.slice(5)
              if (data.startsWith(' ')) data = data.slice(1)
              
              if (data === '[DONE]') {
                loading.value = false
                currentMessages.value[aiMessageIndex].meta = { status: 'done' }
              } else if (data.startsWith('{')) {
                try {
                  const parsed = JSON.parse(data)
                  if (parsed.sessionId) {
                    const newSid = String(parsed.sessionId)
                    if (tempSession.value) {
                      let title = (question || '').trim().slice(0, 30) || `Session ${newSid}`
                      sessions.value[newSid] = {
                        id: newSid,
                        name: title,
                        messages: [...currentMessages.value]
                      }
                      currentSessionId.value = newSid
                      tempSession.value = false
                    }
                  } else if (parsed.type === 'delta') {
                    currentMessages.value[aiMessageIndex].content += parsed.content
                  }
                } catch (e) {
                  currentMessages.value[aiMessageIndex].content += data
                }
              } else {
                currentMessages.value[aiMessageIndex].content += data
              }
              currentMessages.value = [...currentMessages.value]
              await nextTick()
              scrollToBottom()
            }
          }
        }
        
        if (!tempSession.value && currentSessionId.value && sessions.value[currentSessionId.value]) {
             // Sync back to session store
             const sessMsgs = sessions.value[currentSessionId.value].messages
             if (sessMsgs && sessMsgs.length) {
                 sessMsgs[sessMsgs.length - 1].content = currentMessages.value[aiMessageIndex].content
             }
        }

      } catch (err) {
        console.error(err)
        currentMessages.value[aiMessageIndex].meta = { status: 'error' }
      }
    }

    const scrollToBottom = () => {
      if (messagesRef.value) {
        messagesRef.value.scrollTop = messagesRef.value.scrollHeight
      }
    }

    const triggerFileUpload = () => fileInput.value && fileInput.value.click()
    
    const handleFileUpload = async (event) => {
        // ... (keep existing)
      const file = event.target.files[0]
      if (!file) return
      
      try {
        uploading.value = true
        const formData = new FormData()
        formData.append('file', file)
        const response = await api.post('/file/upload', formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })
        if (response.data && response.data.status_code === 1000) {
          uploadedFiles.value.push({ name: file.name })
          ElMessage.success('Uploaded')
        }
      } catch (e) {
        ElMessage.error('Upload failed')
      } finally {
        uploading.value = false
        fileInput.value.value = ''
      }
    }

    const removeFile = (index) => uploadedFiles.value.splice(index, 1)

    onMounted(() => loadSessions())

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
      uploading,
      fileInput,
      uploadedFiles,
      renderMarkdown,
      playTTS,
      createNewSession,
      switchSession,
      syncHistory,
      deleteSession,
      sendMessage,
      triggerFileUpload,
      handleFileUpload,
      removeFile,
      setInputAndSend
    }
  }
}
</script>

<style scoped>
.ai-chat-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  font-family: 'Inter', sans-serif;
  color: #111827;
  padding: 40px;
}

.glass-wrapper {
  display: flex;
  width: 100%;
  max-width: 1400px;
  height: 90vh;
  background: rgba(255, 255, 255, 0.75);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.6);
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

/* Sidebar */
.session-list {
  width: 280px;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: rgba(243, 244, 246, 0.6);
  border-right: 1px solid rgba(229, 231, 235, 0.6);
}

.session-list-header {
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.session-list-header span {
  font-weight: 700;
  font-size: 1.1rem;
  color: #374151;
  display: flex;
  align-items: center;
  gap: 8px;
}

.new-chat-btn {
  background: #E0E7FF;
  color: #4F46E5;
  border: none;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.new-chat-btn:hover {
  background: #4F46E5;
  color: white;
}

.session-list-ul {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  list-style: none;
}

.session-item {
  padding: 12px;
  border-radius: 12px;
  margin-bottom: 4px;
  cursor: pointer;
  color: #4B5563;
  font-size: 0.95rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all 0.2s;
}

.session-item:hover {
  background: rgba(255, 255, 255, 0.5);
}

.session-item.active {
  background: #FFFFFF;
  color: #111827;
  font-weight: 600;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

.session-delete-btn {
  opacity: 0;
  border: none;
  background: transparent;
  color: #9CA3AF;
  cursor: pointer;
  font-size: 1.2rem;
}

.session-item:hover .session-delete-btn {
  opacity: 1;
}

.session-delete-btn:hover {
  color: #EF4444;
}

/* Chat Section */
.chat-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
  background: rgba(255, 255, 255, 0.4);
}

.chat-header {
  height: 64px;
  border-bottom: 1px solid rgba(243, 244, 246, 0.6);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(12px);
  z-index: 10;
}

.header-left, .header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-info h2 {
  font-size: 1rem;
  margin: 0;
  font-weight: 600;
}

.status-dot {
  width: 8px;
  height: 8px;
  background: #10B981;
  border-radius: 50%;
  display: inline-block;
}

.icon-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  color: #6B7280;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.icon-btn:hover {
  background: #F3F4F6;
  color: #111827;
}

.model-select {
  width: 140px;
}

/* Messages */
.chat-messages {
  flex: 1;
  padding: 24px 15%;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #9CA3AF;
  gap: 24px;
}

.empty-icon {
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-state h3 {
  font-size: 1.5rem;
  color: #111827;
  font-weight: 700;
}

.starter-pills {
  display: flex;
  gap: 12px;
}

.pill {
  background: #FFFFFF;
  border: 1px solid #E5E7EB;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 0.9rem;
  color: #4B5563;
  cursor: pointer;
  transition: all 0.2s;
}

.pill:hover {
  border-color: #4F46E5;
  color: #4F46E5;
  background: #EEF2FF;
}

/* Message Rows */
.message-row {
  display: flex;
  gap: 16px;
  max-width: 100%;
}

.user-row {
  flex-direction: row-reverse;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.8);
  color: #4F46E5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.9rem;
  flex-shrink: 0;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.user-row .avatar {
  background: #FCE7F3;
  color: #DB2777;
}

.message-bubble {
  background: rgba(255, 255, 255, 0.85);
  padding: 12px 16px;
  border-radius: 16px;
  border-top-left-radius: 4px;
  max-width: 80%;
  line-height: 1.6;
  font-size: 1rem;
  color: #1F2937;
  position: relative;
  box-shadow: 0 2px 4px rgba(0,0,0,0.02);
}

.user-row .message-bubble {
  background: #4F46E5;
  color: white;
  border-radius: 16px;
  border-top-right-radius: 4px;
}

.message-actions {
  position: absolute;
  bottom: -24px;
  left: 0;
  opacity: 0;
  transition: opacity 0.2s;
}

.message-row:hover .message-actions {
  opacity: 1;
}

.action-btn {
  background: transparent;
  border: none;
  color: #6B7280;
  cursor: pointer;
  padding: 4px;
}

.action-btn:hover {
  color: #4F46E5;
}

/* Input Area */
.chat-input-wrapper {
  padding: 20px 15%;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(8px);
  position: relative;
}

.input-container {
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid #E5E7EB;
  border-radius: 24px;
  padding: 8px 16px;
  display: flex;
  align-items: flex-end;
  gap: 8px;
  transition: border-color 0.2s;
  box-shadow: 0 4px 6px rgba(0,0,0,0.02);
}

.input-container:focus-within {
  border-color: #4F46E5;
  background: #FFFFFF;
  box-shadow: 0 0 0 4px rgba(224, 231, 255, 0.6);
}

.attach-btn {
  background: transparent;
  border: none;
  color: #9CA3AF;
  padding: 8px;
  cursor: pointer;
  border-radius: 50%;
}

.attach-btn:hover {
  background: #E5E7EB;
  color: #4B5563;
}

.main-textarea {
  flex: 1;
  background: transparent;
  border: none;
  resize: none;
  max-height: 120px;
  padding: 10px 0;
  outline: none;
  font-size: 1rem;
  color: #1F2937;
}

.send-btn-round {
  background: #4F46E5;
  color: white;
  border: none;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.send-btn-round:hover:not(:disabled) {
  background: #4338CA;
  transform: scale(1.05);
}

.send-btn-round:disabled {
  background: #E5E7EB;
  cursor: default;
}

.file-preview-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.file-chip {
  background: #F3F4F6;
  border: 1px solid #E5E7EB;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  gap: 6px;
}

.close-file {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 1.1rem;
  line-height: 1;
  color: #9CA3AF;
}

/* Typing Animation */
.loading-bubble {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 60px;
  height: 40px;
}

.typing-dots span {
  display: inline-block;
  width: 6px;
  height: 6px;
  background: #9CA3AF;
  border-radius: 50%;
  margin: 0 2px;
  animation: bounce 1.4s infinite ease-in-out both;
}

.typing-dots span:nth-child(1) { animation-delay: -0.32s; }
.typing-dots span:nth-child(2) { animation-delay: -0.16s; }

@keyframes bounce {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}

/* Markdown Overrides */
:deep(.bubble-content p) {
  margin: 0 0 8px 0;
}
:deep(.bubble-content p:last-child) {
  margin: 0;
}
:deep(.user-row .bubble-content) {
  color: white;
}
:deep(.user-row .bubble-content a) {
  color: #E0E7FF;
}
</style>