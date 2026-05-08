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
- 只认 **DeepL Free API**（`DEEPL_API_KEY` 环境变量）
- 不引入其他翻译服务

## 关键路径
```
main.go          → Wails App 入口
app.go           → 前后端绑定方法（Go → JS 通过 wailsjs/ 自动生成）
db.go            → SQLite 操作
frontend/src/    → React 前端
build/           → 打包产物配置
```

## 开发环境检查
```bash
cs check          # 项目级环境自检（wails / go / node 是否就绪）
```

## 红线
- **禁止**修改 `wailsjs/` 下的自动生成文件，只改 `app.go` 后重新生成
- **禁止**引入第三方 UI 组件库（Tailwind、MUI 等）
- **禁止**用 CGO，保持纯 Go 编译

## Agent 分工
- Claude：前端 UI、CSS、React 组件
- Codex：后端 Go 逻辑、DB 操作
