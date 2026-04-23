<div align="center"><h1>Yatori-UI</h1></div>

<div align="center"><h2>智能网课助手 · 桌面端</h2></div>

<div align="center">
<img width="125px" src="https://img.shields.io/badge/GO1.24.0-building-r.svg?logo=go"></img>
<img width="100px" src="https://img.shields.io/badge/Wails-v2.11.0-blue.svg"></img>
<img width="90px" src="https://img.shields.io/badge/React-18.2-61DAFB.svg?logo=react"></img>
<img width="80px" src="https://img.shields.io/badge/TypeScript-4.6-3178C6.svg?logo=typescript"></img>
</div>

---

## 📖 项目简介
**Yatori-UI**核心代码由"https://github.com/yatori-dev/yatori-go-console.git"作者：Yatori-Dev提供
**Yatori-UI** 是一款基于 [Wails v2](https://wails.io/) 框架构建的跨平台桌面端智能网课辅助工具。它将强大的 Go 后端与现代化的 React 前端融为一体，提供流畅的图形化操作体验，支持多个主流在线学习平台的自动学习、自动答题等功能。

无需浏览器，无需复杂配置，开箱即用。通过直观的 GUI 界面即可完成账号管理、任务监控、配置调整等全部操作。

## 🎯 功能/特性

| 功能/特性 | 状态 |
| --- | --- |
| 独立桌面程序，不依赖浏览器 | ✅ |
| 图形化界面操作，开箱即用 | ✅ |
| AI 自动识别跳过验证码 | ✅ |
| 多账号同刷 | ✅ |
| 支持状态邮箱通知 | ✅ |
| 支持自动考试 | ✅ |
| 答题支撑 AI 大模型加持 | ✅ |
| 灵活配置文件（YAML） | ✅ |
| 自动继续上次记录时长刷课 | ✅ |
| 系统托盘最小化运行 | ✅ |
| 配置文件热更新（自动监听变更） | ✅ |
| 本地 SQLite 数据持久化 | ✅ |
| 部分平台支持暴力模式（无视前置课程学习限制所有视频同刷） | ✅ |

## 🎯 支持平台

| 平台 | 描述 | 状态 |
| --- | --- | --- |
| 英华学堂 | 支持暴力模式（会被检测到） | 已完成 ✅ |
| 仓辉实训 | 支持暴力模式（套壳英华版本会被检测到） | 已完成 ✅ |
| 创能实训 | 支持暴力模式（会被检测到） | 已完成 ✅ |
| 社会公益课 | 支持暴力模式（会被检测到） | 已完成 ✅ |
| 重庆工业学院 CQIE | 支持暴力模式（支持秒刷） | 已完成 ✅ |
| 学习公社（ENAEA） | 支持暴力模式（倍速刷） | 已完成 ✅ |
| 大学生网络党校（ENAEA） | 支持暴力模式（倍速刷） | 已完成 ✅ |
| 中小学网络党校（ENAEA） | 支持暴力模式（倍速刷） | 已完成 ✅ |
| 学习通 | 支持绕过人脸认证，支持自动写章测、作业、考试，以及多课程、多任务点模式 | 已完成 ✅ |
| 码上研训 | 默认秒刷 | 已完成 ✅ |
| 随行课堂 | 支持秒刷完成度以及学时累计刷取 | 已完成 ✅ |
| 智慧职教（资源库） | 默认秒刷（目前只支持 Cookie 登录方式） | 已完成 ✅ |
| 青书学堂 | 只支持普通模式 | 已完成 ✅ |
| WeLearn | 支持 | 已完成 ✅ |
| 海旗科技 | 支持 | 已完成 ✅ |
| 工学云 | 支持 | 已完成 ✅ |
| 安全微伴 | 开发中 | 开发中 🚧 |

## 🤖 AI 答题支持

| AI 类型 | 名称 |
| --- | --- |
| TONGYI | 通义千问 |
| CHATGLM | 智谱 ChatGLM |
| XINGHUO | 讯飞星火 |
| DOUBAO | 豆包 |
| OPENAI | OpenAI |
| DEEPSEEK | DeepSeek |
| SILICON | 硅基流动 |
| METAAI | 秘塔 AI |
| OTHER | 其他兼容 |

## 🛠️ 技术栈

### 后端（Go）

| 类别 | 技术 | 说明 |
| --- | --- | --- |
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
| --- | --- | --- |
| UI 框架 | React 18.2 | 声明式组件化界面开发 |
| 路由 | react-router-dom 6 | 单页应用路由管理 |
| 状态管理 | Zustand 5 | 轻量级响应式状态管理 |
| CSS 框架 | Tailwind CSS 3.4 | 原子化 CSS 工具类框架 |
| 动画 | Framer Motion 12 | 流畅的 UI 交互动画 |
| 图标 | Lucide React | 精美的开源图标库 |
| 构建工具 | Vite 3 | 极速前端构建与热更新 |
| 语言 | TypeScript 4.6 | 类型安全的 JavaScript 超集 |

### 架构特点

- **分层架构**：后端采用 `Entity → DAO → Service → App` 标准分层，职责清晰
- **事件驱动**：通过 Wails `EventsEmit/EventsOn` 实现前后端实时通信与进度推送
- **并发模型**：每个账号的刷课任务在独立 goroutine 中运行，支持多任务并发
- **无边框窗口**：自定义标题栏，暗色系 + 霓虹风格 UI，科技感十足
- **配置热更新**：自动监听 `config.yaml` 变更，无需重启即可生效

---

## 👤 作者

**❦Angelic 音乐**

📧 联系方式：QQ：2844189228

---

## ⚠️ 免责声明

> 1. 本程序仅供学习、研究与技术交流使用，严禁用于任何商业用途或违法活动。
> 2. 本程序开源免费，严禁贩卖。若因使用本程序对相关平台或机构造成任何损失，请立即停止使用并删除本程序。
> 3. 任何个人或组织使用本程序所从事的一切违法行为，均与作者无关，作者不承担任何法律责任。
> 4. 使用本程序即表示您已阅读、理解并同意遵守上述声明。如不同意，请立即删除本程序。
> 5. 本程序涉及的部分功能可能违反相关平台的使用条款，使用者需自行承担由此产生的一切风险与后果。
