# word-reader — 项目宪法 (Single Source of Truth)

## 📋 项目概览
本项目是一个基于 Wails 的沉浸式英语阅读与背单词工具。本文件是项目级指令的唯一事实来源，所有 Agent 必须严格遵守。

## 🛠️ Stack (技术栈)
- **Runtime**: Wails v2 + Go 1.25
- **Frontend**: React + TypeScript + Vanilla CSS
- **DB**: SQLite (`modernc.org/sqlite`) - **纯 Go 实现，禁用 CGO**
- **CI**: GitHub Actions (Ubuntu 需使用 `libwebkit2gtk-4.1-dev`)

## 🎨 UI / UX 基调 (沉浸式阅读)
- **主题色**: 琥珀暖金 (`#F59E0B` / Amber-400)
- **背景**: 极深暖黑 (`#1C1A17`)
- **氛围**: 深夜书屋感，禁用冰冷科技蓝
- **细节**: 弹窗模糊 (`backdrop-filter`)，行间距 1.85，琥珀色 Hover 效果

## ⚙️ 核心逻辑规范
- **翻译 (Translation)**: 只认 DeepL Free API。
- **存储 (Storage)**: API Key 存储在 DB 的 `settings` 表中，**严禁**依赖环境变量。
- **发音 (TTS)**: **必须**在前端通过 Web Speech API 实现，禁止后端执行系统命令。

## 🤝 协作分工 (Role Alignment)
- **Gemini**: **首席前端专家**。负责 UI 实现、CSS 架构、React 交互逻辑。
- **Other Agents**: 处理后端 Go 逻辑、数据库迁移及系统架构。

## 🚩 红线 (Hard Redlines)
- **禁止**手动修改 `wailsjs/` 下的自动生成文件。
- **禁止**引入第三方 UI 组件库（如 Tailwind, MUI, AntD），保持 Vanilla CSS 原生性。
- **禁止**引入 CGO 依赖，确保跨平台编译兼容性。

## 🔍 环境检查
在执行开发前，请运行 `cs check` 确保本地 Wails 和 Go 环境已就绪。
