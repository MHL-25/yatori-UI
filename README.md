<div align="center">
  <h1>🎓 Yatori-UI</h1>
  <p><strong>智能网课助手 · 桌面端</strong></p>
  <p>一款基于 Wails v2 框架构建的跨平台桌面端智能网课辅助工具</p>
  
  <div style="margin: 20px 0;">
    <img src="https://img.shields.io/badge/Version-1.2.19-blue.svg?style=flat-square" alt="Version">
    <img src="https://img.shields.io/badge/Go-1.24.0-00ADD8.svg?style=flat-square&logo=go" alt="Go">
    <img src="https://img.shields.io/badge/Wails-v2.11.0-FF6B6B.svg?style=flat-square" alt="Wails">
    <img src="https://img.shields.io/badge/React-18.2-61DAFB.svg?style=flat-square&logo=react" alt="React">
    <img src="https://img.shields.io/badge/TypeScript-4.6-3178C6.svg?style=flat-square&logo=typescript" alt="TypeScript">
    <img src="https://img.shields.io/badge/License-MIT-green.svg?style=flat-square" alt="License">
  </div>
</div>

---

## 支持项目

如果觉得本项目做的不错，可以赞助一下吗（双手合十），你的赞助将会是我维护的最大动力，当然也可以不赞助.....哦

爱发电地址：**https://afdian.com/a/yatori-UI**

---

## 📖 项目简介

**Yatori-UI** 核心代码由 [yatori-dev/yatori-go-console](https://github.com/yatori-dev/yatori-go-console) 作者：Yatori-Dev 提供

**Yatori-UI** 是一款基于 [Wails v2](https://wails.io/) 框架构建的跨平台桌面端智能网课辅助工具。它将强大的 Go 后端与现代化的 React 前端融为一体，提供流畅的图形化操作体验，支持学习通、英华学堂、仓辉实训、智慧职教、青书学堂、WeLearn 等多个主流在线学习平台的自动学习、自动答题与课程质检。

### 🆕 v1.2.19 更新亮点

- **英华多平台 URL**：同一个英华账号可同时配置多个学校平台地址，不再需要重复添加多个账号。
- **账号密码自动登录**：只需填写账号和密码，程序自动识别验证码并抓取各平台 Cookie；登录失败会自动重新抓取，账号或密码错误不会反复重试。
- **平台独立密码**：各学校平台密码不同时可逐平台填写；密码一致时使用“统一密码按钮”一键套用。
- **多账号稳定性**：完善异常隔离、事件派发与运行时自适应，多账号并发更稳定、界面更流畅。
- **答题与图片分流**：题目含图片时优先使用识图 AI，纯文本题目使用文本 AI；填空题按空位严格解析，避免错填。

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
| 课程质检 | ✅ | 单账号 / 指定账号 / 全部账号检测课程完成度与章节锁状态 |
| 前置与定时锁调度 | ✅ | 受限课程单独串行补前置，定时节点记录开放时间并到点自动续刷 |
| 账号自检时间窗 | ✅ | 按账号配置时段，启动时依据系统时间决定是否自动质检/刷课 |
| 互动测验 | ✅ | 学习通视频弹题自动识别并提交答案 |
| 启动次数赞助提示 | ✅ | 达到伪随机次数时提示，可关闭或不再提示 |
| 多账号并发自检 | ✅ | 启动自动质检最多 5 个账号同时处理 |
| 英华多平台 URL | ✅ | 同一账号可配置多个学校平台，自动登录并抓取 Cookie |
| 平台独立密码 | ✅ | 每个平台可单独设置密码，密码一致时支持统一密码 |
| 图片题识别分流 | ✅ | 含图片题目走识图 AI，纯文本题目走文本 AI |

---

## 🎯 支持平台

> 注：智慧树平台暂未接入，将在后续版本支持。

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
      <td>支持账号密码自动识别验证码登录、同一账号配置多个平台 URL、每个平台独立密码，以及暴力模式（会被检测到）</td>
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
├── app.go                      # Wails 应用入口与自动化调度
├── main.go                     # 主程序入口
├── wails.json                  # Wails 配置文件
├── go.mod                      # Go 依赖管理
├── go.sum                      # Go 依赖校验
├── internal/                   # 内部包
│   ├── activity/               # 各平台刷课逻辑与多平台会话管理
│   ├── checker/                # 课程质检与登录校验
│   ├── config/                 # 配置结构与英华多平台目标解析
│   ├── dao/                    # SQLite 数据访问
│   ├── entity/                 # 数据实体与请求模型
│   ├── monitor/                # 进度与日志事件总线
│   ├── runtimecfg/             # CPU/内存自适应运行时配置
│   ├── service/                # 账号、任务与配置业务服务
│   ├── studyplan/              # 前置任务与定时开放分析
│   ├── yinghuahelper/          # 英华多平台登录与 Cookie 抓取
│   ├── email/                  # 邮件通知
│   └── notify/                 # ShowDoc 等外部通知
├── core-patch/                 # 核心库补丁源码（图片分流、填空解析、稳定性修复）
├── frontend/                   # 前端源码
│   ├── src/
│   │   ├── components/         # 通用组件
│   │   ├── pages/              # 页面组件（账号、监控、设置、质检详情等）
│   │   ├── stores/             # Zustand 状态管理
│   │   ├── utils/              # 工具函数与版本号
│   │   └── App.tsx             # 应用入口组件
│   ├── package.json            # 前端依赖
│   └── vite.config.ts          # Vite 配置
├── build/                      # 构建资源
│   ├── appicon.png             # 应用图标
│   └── windows/                # Windows 构建配置
├── config.yaml                 # 运行时配置文件（首次启动生成）
└── README.md                   # 项目说明文档
```

---

## ⚙️ 配置说明

账号与刷课配置推荐直接在软件的“账号管理”页面完成，数据保存在本地 SQLite 数据库。程序首次启动会在同目录生成 `config.yaml`，用于 AI、通知与自动化等全局设置；也支持在 `users` 中预置账号后导入。

```yaml
setting:
  basicSetting:
    autoCheckSw: 1        # 启动时自动质检
    autoStartSw: 1        # 质检后自动开始刷课
  aiSetting:
    aiType: DEEPSEEK      # 文本 AI 类型
    aiUrl: ""
    model: ""
    API_KEY: ""
  visionAiSetting:
    aiType: TONGYI        # 识图 AI（题目含图片时使用）
    aiUrl: ""
    model: ""
    API_KEY: ""
  apiQueSetting:
    url: http://localhost:8083
users:
  - accountType: XUEXITONG
    account: "your_account"
    password: "your_password"
    url: ""
  - accountType: YINGHUA
    account: "your_account"
    password: "统一密码"
    # 同一账号的多个学校平台，密码留空时使用统一密码
    yingHuaUrls:
      - url: "https://cqcst.suwankj.com/user/login"
      - url: "https://cqcst.tuozhikj.com/user/member/pwd"
        password: "该平台单独密码"
      - url: "https://cqcst.chaoxiankeji.com/user/login"
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
