# word-reader — 项目宪法

## Stack
- **Runtime**: Wails v2 + Go 1.25
- **Frontend**: React + TypeScript + Vanilla CSS (CSS Variables)
- **DB**: SQLite via `modernc.org/sqlite`（纯 Go，无 CGO）
- **CI**: GitHub Actions — Ubuntu 需用 `libwebkit2gtk-4.1-dev`（非 4.0）

## UI / UX 基调
- 主色：**琥珀暖金** (`#F59E0B` / Amber-400)
- 背景：**极深暖黑** (`#1C1A17`)
- 氛围：深夜书屋，沉浸阅读，禁用冰冷科技蓝
- 单词 hover：琥珀色背景 + 平滑下划线
- 弹窗：磨砂玻璃 `backdrop-filter: blur`
- 行间距：1.85

## 翻译
- 只认 **DeepL Free API**
- **存储**: 存储于数据库 `settings` 表（`deepl_api_key`），不再依赖环境变量
- **发音 (TTS)**: 必须在前端通过 Web Speech API (`utils/tts.ts`) 实现，禁止后端执行系统命令

## 开发环境检查
```bash
cs check          # 项目级环境自检（wails / go / node 是否就绪）
```

## 红线
- **禁止**修改 `wailsjs/` 下的自动生成文件，只改 `app.go` 后重新生成（如环境缺失 Wails CLI，需提醒用户修复）
- **禁止**引入第三方 UI 组件库（Tailwind、MUI 等），保持 Vanilla CSS 原生感
- **禁止**用 CGO，保持纯 Go 编译

## 协作分工 (Gemini User Preference)
- **Gemini**: 专注于 **前端 UI/UX 实现**、CSS 系统、React 交互逻辑。
- **Other Agents**: 处理后端 Go 逻辑、DB 迁移及系统级运维。
