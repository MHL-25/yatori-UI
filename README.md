<div align="center">
  <h1>🎓 Yatori-UI</h1>
  <p><strong>智能网课助手 · 桌面端</strong></p>
  <p>一款基于 Wails v2 框架构建的跨平台桌面端智能网课辅助工具</p>
  
  <div style="margin: 20px 0;">
    <img src="https://img.shields.io/badge/Version-1.1.4-blue.svg?style=flat-square" alt="Version">
    <img src="https://img.shields.io/badge/Go-1.24.0-00ADD8.svg?style=flat-square&logo=go" alt="Go">
    <img src="https://img.shields.io/badge/Wails-v2.11.0-FF6B6B.svg?style=flat-square" alt="Wails">
    <img src="https://img.shields.io/badge/React-18.2-61DAFB.svg?style=flat-square&logo=react" alt="React">
    <img src="https://img.shields.io/badge/TypeScript-4.6-3178C6.svg?style=flat-square&logo=typescript" alt="TypeScript">
    <img src="https://img.shields.io/badge/License-MIT-green.svg?style=flat-square" alt="License">
  </div>
</div>

---

## 📖 项目简介

**Yatori-UI** 核心代码由 [yatori-dev/yatori-go-console](https://github.com/yatori-dev/yatori-go-console) 作者：Yatori-Dev 提供

**Yatori-UI** 是一款基于 [Wails v2](https://wails.io/) 框架构建的跨平台桌面端智能网课辅助工具。它将强大的 Go 后端与现代化的 React 前端融为一体，提供流畅的图形化操作体验，支持多个主流在线学习平台的自动学习、自动答题等功能。

### ✨ 核心优势

- 🚀 **开箱即用** - 无需浏览器，无需复杂配置，下载即用
- 🎨 **精美界面** - 暗色系 + 霓虹风格 UI，科技感十足
- 🔧 **简单易用** - 直观的 GUI 界面，可视化操作
- ⚡ **高效稳定** - Go 后端高性能，支持多账号并发
- 🤖 **智能答题** - 集成多种 AI 大模型，自动识别题目

---

## 🎯 功能特性

| 功能 | 状态 | 说明 |
|:---|:---:|:---|
| 独立桌面程序 | ✅ | 不依赖浏览器，原生桌面体验 |
| 图形化界面操作 | ✅ | 开箱即用，直观易用 |
| AI 自动识别验证码 | ✅ | 智能跳过验证码干扰 |
| 多账号同时刷课 | ✅ | 支持批量添加和管理账号 |
| 实时任务监控 | ✅ | 实时查看刷课进度和状态 |
| 支持状态邮箱通知 | ✅ | 任务完成自动通知 |
| 支持自动考试 | ✅ | 自动完成考试答题 |
| AI 大模型答题 | ✅ | 集成多种 AI 模型辅助答题 |
| 灵活配置文件 | ✅ | YAML 格式，支持热更新 |
| 自动续刷 | ✅ | 自动继续上次记录时长刷课 |
| 系统托盘运行 | ✅ | 最小化到托盘后台运行 |
| 本地数据持久化 | ✅ | SQLite 数据库存储学习记录 |
| 暴力模式 | ✅ | 部分平台支持，无视前置课程限制 |

---

## 🎯 支持平台

<table>
  <thead>
    <tr>
      <th>平台</th>
      <th>特性</th>
      <th>状态</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>📖 学习通</td>
      <td>支持绕过人脸认证，自动写章测、作业、考试，以及多课程、多任务点模式</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>🎓 英华学堂</td>
      <td>支持暴力模式（会被检测到）</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>🔧 仓辉实训</td>
      <td>支持暴力模式（套壳英华版本会被检测到）</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>💻 创能实训</td>
      <td>支持暴力模式（会被检测到）</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>📚 社会公益课</td>
      <td>支持暴力模式（会被检测到）</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>🏫 重庆工业学院 CQIE</td>
      <td>支持暴力模式（支持秒刷）</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>📚 学习公社（ENAEA）</td>
      <td>支持暴力模式（倍速刷）</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>🎓 大学生网络党校（ENAEA）</td>
      <td>支持暴力模式（倍速刷）</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>📖 中小学网络党校（ENAEA）</td>
      <td>支持暴力模式（倍速刷）</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>💻 码上研训</td>
      <td>默认秒刷</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>🌐 随行课堂</td>
      <td>支持秒刷完成度以及学时累计刷取</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>🎯 智慧职教（资源库）</td>
      <td>默认秒刷（目前只支持 Cookie 登录方式）</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>📗 青书学堂</td>
      <td>只支持普通模式</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>🌐 WeLearn</td>
      <td>支持自动学习</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>🚀 海旗科技</td>
      <td>支持自动学习</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>☁️ 工学云</td>
      <td>支持自动学习</td>
      <td>✅ 已完成</td>
    </tr>
    <tr>
      <td>🛡️ 安全微伴</td>
      <td>开发中</td>
      <td>🚧 进行中</td>
    </tr>
  </tbody>
</table>

---

## 🤖 AI 答题支持

支持多种主流 AI 大模型，智能识别题目并自动答题：

| AI 类型 | 名称 | 说明 |
|:---|:---|:---|
| `TONGYI` | 通义千问 | 阿里云 AI 模型 |
| `CHATGLM` | 智谱 ChatGLM | 智谱 AI 模型 |
| `XINGHUO` | 讯飞星火 | 科大讯飞 AI 模型 |
| `DOUBAO` | 豆包 | 字节跳动 AI 模型 |
| `OPENAI` | OpenAI | ChatGPT 系列模型 |
| `DEEPSEEK` | DeepSeek | 深度求索 AI 模型 |
| `SILICON` | 硅基流动 | 硅基流动 AI 模型 |
| `METAAI` | 秘塔 AI | 秘塔科技 AI 模型 |
| `OTHER` | 其他兼容 | 其他兼容 OpenAI 接口的模型 |

---

## 🛠️ 技术栈

### 后端（Go）

| 类别 | 技术 | 说明 |
|:---|:---|:---|
| 语言 | Go 1.24.0 | 核心后端语言 |
| 桌面框架 | Wails v2.11.0 | Go + Web 前端融合的桌面应用框架 |
| ORM | GORM v1.31.1 | Go 语言 ORM 库 |
| 数据库 | SQLite | 本地轻量级数据持久化 |
| 配置管理 | Viper v1.21.0 | 支持 YAML/JSON 等多格式配置读写 |
| 文件监听 | fsnotify v1.9.0 | 配置文件热更新监听 |
| 系统托盘 | getlantern/systray | 托盘图标与菜单支持 |
| 核心业务库 | yatori-go-core v1.9.1 | 多平台网课自动化核心逻辑 |
| 验证码识别 | ddddocr-go | ONNX 推理驱动的验证码 OCR |
| HTTP 框架 | Echo v4 | 轻量级高性能 Web 框架（间接依赖） |
| WebSocket | gorilla/websocket | 实时通信支持（间接依赖） |

### 前端（TypeScript + React）

| 类别 | 技术 | 说明 |
|:---|:---|:---|
| UI 框架 | React 18.2 | 声明式组件化界面开发 |
| 路由 | react-router-dom 6 | 单页应用路由管理 |
| 状态管理 | Zustand 5 | 轻量级响应式状态管理 |
| CSS 框架 | Tailwind CSS 3.4 | 原子化 CSS 工具类框架 |
| 动画 | Framer Motion 12 | 流畅的 UI 交互动画 |
| 图标 | Lucide React | 精美的开源图标库 |
| 构建工具 | Vite 3 | 极速前端构建与热更新 |
| 语言 | TypeScript 4.6 | 类型安全的 JavaScript 超集 |

### 架构特点

```
┌─────────────────────────────────────────────────────────────┐
│                    Yatori-UI 架构图                         │
├─────────────────────────────────────────────────────────────┤
│  前端层 (React + TypeScript)                                │
│  ├── 页面组件 (Pages)                                       │
│  ├── 状态管理 (Zustand)                                     │
│  └── UI 组件 (Tailwind + Framer Motion)                     │
├─────────────────────────────────────────────────────────────┤
│  通信层 (Wails Events)                                      │
│  ├── EventsEmit - 前端 → 后端                               │
│  └── EventsOn - 后端 → 前端                                 │
├─────────────────────────────────────────────────────────────┤
│  后端层 (Go)                                                │
│  ├── App 层 - 业务入口                                      │
│  ├── Service 层 - 业务逻辑                                  │
│  ├── DAO 层 - 数据访问                                      │
│  └── Entity 层 - 数据模型                                   │
├─────────────────────────────────────────────────────────────┤
│  核心层 (yatori-go-core)                                    │
│  ├── 平台 API 封装                                          │
│  ├── 账号认证                                               │
│  ├── 课程学习                                               │
│  └── 作业考试                                               │
└─────────────────────────────────────────────────────────────┘
```

- **分层架构**：后端采用 `Entity → DAO → Service → App` 标准分层，职责清晰
- **事件驱动**：通过 Wails `EventsEmit/EventsOn` 实现前后端实时通信与进度推送
- **并发模型**：每个账号的刷课任务在独立 goroutine 中运行，支持多任务并发
- **无边框窗口**：暗色系 + 霓虹风格 UI，科技感十足
- **配置热更新**：自动监听 `config.yaml` 变更，无需重启即可生效

---

## 🚀 快速开始

### 方式一：直接下载（推荐）

1. 前往 [Releases](https://github.com/MHL-25/yatori-UI/releases) 页面下载最新版本
2. 解压后运行 `yatori-UI.exe` 即可使用

### 方式二：从源码构建

#### 环境要求

- Go 1.24.0+
- Node.js 18+
- Wails CLI v2.11.0+

#### 构建步骤

```bash
# 1. 克隆项目
git clone https://github.com/MHL-25/yatori-UI.git
cd yatori-UI

# 2. 安装前端依赖
cd frontend
npm install
cd ..

# 3. 使用 Wails 构建
wails build

# 4. 运行程序
./build/bin/yatori-UI.exe
```

### 开发模式

```bash
# 启动开发服务器（支持热更新）
wails dev
```

---

## 📁 项目结构

```
yatori-UI/
├── app.go                      # Wails 应用入口
├── main.go                     # 主程序入口
├── wails.json                  # Wails 配置文件
├── go.mod                      # Go 依赖管理
├── go.sum                      # Go 依赖校验
├── internal/                   # 内部包
│   ├── activity/               # 平台活动实现
│   │   └── activity.go         # 所有平台的刷课逻辑
│   ├── config/                 # 配置管理
│   │   └── config.go           # 配置结构体定义
│   ├── service/                # 服务层
│   │   └── service.go          # 业务服务实现
│   ├── monitor/                # 监控模块
│   │   └── event_bus.go        # 事件总线实现
│   └── database/               # 数据库模块
│       ├── db.go               # 数据库初始化
│       └── entity/             # 数据实体
├── frontend/                   # 前端源码
│   ├── src/
│   │   ├── components/         # 通用组件
│   │   ├── pages/              # 页面组件
│   │   │   ├── AccountsPage    # 账号管理页面
│   │   │   ├── MonitorPage     # 任务监控页面
│   │   │   ├── SettingsPage    # 设置页面
│   │   │   └── ...
│   │   ├── stores/             # 状态管理
│   │   ├── utils/              # 工具函数
│   │   └── App.tsx             # 应用入口组件
│   ├── package.json            # 前端依赖
│   └── vite.config.ts          # Vite 配置
├── build/                      # 构建资源
│   ├── appicon.png             # 应用图标
│   └── windows/                # Windows 构建配置
├── config.yaml                 # 配置文件（运行时生成）
└── README.md                   # 项目说明文档
```

---

## ⚙️ 配置说明

程序运行后会在同目录下生成 `config.yaml` 配置文件，支持以下配置：

```yaml
# 基础配置
setting:
  basicSetting:
    cookie: ""                    # Cookie 登录（部分平台需要）
    
  # 账号列表
  users:
    - platformType: "XUEXITONG"  # 平台类型
      username: "your_username"   # 用户名
      password: "your_password"   # 密码
      
      # AI 答题配置（可选）
      aiSetting:
        aiType: "DEEPSEEK"       # AI 类型
        apiKey: "your_api_key"   # API Key
        apiUrl: ""               # 自定义 API 地址（可选）
        
      # 学习设置
      studySetting:
        videoModel: 1            # 视频模式：0=不刷，1=串行，2=并发课程，3=多任务点
        autoExam: true           # 自动考试
        autoWork: true           # 自动作业
```

### 视频模式说明

| 模式 | 说明 |
|:---|:---|
| `0` | 不刷视频 |
| `1` | 串行模式 - 一个视频看完再看下一个 |
| `2` | 并发课程模式 - 多个课程同时刷 |
| `3` | 多任务点模式 - 单课程内多个任务点同时刷（推荐） |

---

## 👤 作者

**❦Angelic 音乐**

- 📧 QQ：2844189228
- 🔗 GitHub：[@MHL-25](https://github.com/MHL-25)

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

---

## 📄 许可证

本项目基于 MIT 许可证开源 - 详见 [LICENSE](LICENSE) 文件

---

## ⚠️ 免责声明

> **重要提示**：使用本程序前请仔细阅读以下声明

1. 本程序仅供学习、研究与技术交流使用，严禁用于任何商业用途或违法活动。
2. 本程序开源免费，严禁贩卖。若因使用本程序对相关平台或机构造成任何损失，请立即停止使用并删除本程序。
3. 任何个人或组织使用本程序所从事的一切违法行为，均与作者无关，作者不承担任何法律责任。
4. 使用本程序即表示您已阅读、理解并同意遵守上述声明。如不同意，请立即删除本程序。
5. 本程序涉及的部分功能可能违反相关平台的使用条款，使用者需自行承担由此产生的一切风险与后果。

---

## 🙏 致谢

- [yatori-dev/yatori-go-console](https://github.com/yatori-dev/yatori-go-console) - 核心代码提供
- [Wails](https://wails.io/) - 优秀的 Go 桌面应用框架
- [React](https://reactjs.org/) - 前端 UI 库
- [Tailwind CSS](https://tailwindcss.com/) - CSS 框架

---

<div align="center">
  <p>如果觉得有用，请给个 ⭐ Star 支持一下！</p>
</div>
