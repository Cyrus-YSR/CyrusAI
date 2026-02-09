# CyrusAI: 全栈 AI 应用服务平台

CyrusAI 是一个基于 Go 语言构建的现代化 AI 应用平台，集成了大模型对话、图像识别与分析、文件管理等功能。项目采用前后端分离架构，后端基于 Gin 框架和 CloudWeGo/Eino 生态，前端使用 Vue.js 3。

> ⚠️ **重要提示 / Important Note**
>
> 本项目依赖 **ONNX Runtime** 进行本地图像推理，使用的动态链接库 (`libonnxruntime.so`) 仅支持 **Linux 环境**。
>
> - **Windows/macOS 用户**：请务必使用 **Docker** 或 **WSL2** 运行本项目。
> - **Linux 用户**：可直接在本地运行。

---

## 🚀 功能特性 (Features)

*   **🤖 智能对话 (AI Chat)**
    *   集成 OpenAI 兼容接口，支持多轮对话。
    *   支持 MCP (Model Context Protocol) 扩展，增强模型能力。
    *   基于 CloudWeGo/Eino 的智能体编排。
    *   会话历史管理与持久化。

*   **👁️ 图像识别与分析 (Image Analysis)**
    *   **本地推理**：内置 ResNet50 (v2) 模型，基于 ONNX Runtime 实现高性能本地图像分类。
    *   **深度解析**：结合大语言模型对识别结果进行详细的中文解读与场景分析。

*   **📁 文件管理 (File Management)**
    *   支持文件上传与管理。
    *   结合 RAG (检索增强生成) 技术的文件内容解析（开发中）。

*   **🔐 用户系统 (User System)**
    *   完整的注册、登录流程。
    *   基于 JWT 的安全认证机制。

---

## 🛠️ 技术栈 (Tech Stack)

*   **后端 (Backend)**: Go 1.25+, Gin, GORM, CloudWeGo/Eino
*   **前端 (Frontend)**: Vue.js 3, Element Plus
*   **数据库 (Database)**: MySQL 8.0
*   **缓存 (Cache)**: Redis
*   **消息队列 (Message Queue)**: RabbitMQ
*   **AI 推理 (Inference)**: ONNX Runtime (Go bindings)

---

## 🐳 快速开始 (Docker 运行 - 推荐)

这是在 Windows 和 macOS 上运行本项目的**唯一推荐方式**。

### 1. 前置准备
确保本地已安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/)。

### 2. 配置环境变量
修改 `docker-compose.yml` 中的环境变量，填入你的 OpenAI API Key：

```yaml
    environment:
      - OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxx  # 替换为你的 Key
      - OPENAI_MODEL_NAME=qwen-plus                 # 或 gpt-4, gpt-3.5-turbo
      - OPENAI_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1 # 替换为对应的 Base URL
```

### 3. 启动服务
在项目根目录下运行：

```bash
docker-compose up -d
```

### 4. 访问应用
服务启动后，可以通过浏览器访问：

*   **前端页面**: [http://localhost:8081](http://localhost:8081)
*   **后端接口**: [http://localhost:9090](http://localhost:9090)

---

## 🐧 Linux 本地运行 (开发指南)

如果你是在 Linux 环境下开发，或者使用 WSL2，可以按照以下步骤手动运行。

### 前置依赖
*   Go 1.25+
*   MySQL, Redis, RabbitMQ (建议通过 Docker 启动这些中间件)
*   `libonnxruntime.so` (项目 `libs/` 目录下已包含，需配置路径)

### 运行步骤

1.  **启动依赖服务**：
    ```bash
    # 仅启动数据库和中间件
    docker-compose up -d mysql redis rabbitmq
    ```

2.  **配置动态链接库**：
    确保系统能找到 `libonnxruntime.so`。
    ```bash
    export ORT_DYLIB_PATH=$(pwd)/libs/libonnxruntime.so
    export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$(pwd)/libs
    ```

3.  **运行后端**：
    ```bash
    go mod tidy
    go run main.go
    ```

4.  **运行前端**：
    ```bash
    cd vue-frontend
    npm install
    npm run serve
    ```

---

## 📂 目录结构 (Project Structure)

```
CyrusAI/
CyrusAI/
├── common/             # 公共组件 (AI模型, 图像识别, 消息队列等)
├── config/             # 配置文件
├── controller/         # 控制器层 (HTTP 接口处理)
├── dao/                # 数据访问层 (Database Access)
├── libs/               # 第三方动态库 (libonnxruntime.so)
├── middleware/         # 中间件 (JWT, CORS)
├── model/              # 数据模型定义
├── models/             # AI 模型文件 (.onnx)
├── router/             # 路由定义
├── service/            # 业务逻辑层
├── vue-frontend/       # Vue 前端源码
├── docker-compose.yml  # Docker 编排文件
└── main.go             # 程序入口
```

## 📝 常见问题 (FAQ)

**Q: 为什么在 Windows 上直接运行 `go run main.go` 会报错？**
A: 因为项目依赖的 `onnxruntime-go` 库需要加载 `libonnxruntime.so`，这是一个 Linux 格式的动态链接库，无法在 Windows 上运行。请使用 Docker 或 WSL2。

**Q: 图像识别报错 "Invalid output name"？**
A: 请确保使用了正确的 ONNX 模型文件，并且代码中的 Input/Output 节点名称与模型一致。当前项目适配的是 ResNet50 v2-7 模型。
