package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ── App 结构体 ──────────────────────────────────────────────────────────────
// Go 没有 class，用 struct + 方法实现面向对象
// struct 定义数据，func (a *App) xxx() 定义方法
// 类似前端的 class App { constructor() {} method() {} }

// App 应用核心结构体，绑定到前端后，前端可以调用它的所有公开方法
// 大写字母开头的方法 = 公开（public），小写 = 私有（private）
type App struct {
	ctx context.Context // 上下文，用于控制生命周期（取消、超时等）
	db  *sql.DB         // 数据库连接池，整个应用共享
}

// NewApp 构造函数 — 创建 App 实例
// Go 没有构造函数语法，用 NewXxx 函数代替
// 类似前端的 new App()
func NewApp() *App {
	return &App{}
}

// startup 应用启动回调，由 Wails 框架在窗口创建后调用
// 类似前端的 onMounted 生命周期钩子
// ctx 由框架传入，用于后续的上下文控制
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化数据库，失败只打印日志，不退出
	// 生产环境应该给用户提示
	db, err := initDB()
	if err != nil {
		fmt.Printf("数据库初始化失败: %v\n", err)
		return
	}
	a.db = db
}

// ── 设置 ──────────────────────────────────────────────────────────────────

// GetSetting 获取配置
func (a *App) GetSetting(key string) (string, error) {
	var value string
	err := a.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// SetSetting 保存配置
func (a *App) SetSetting(key, value string) error {
	_, err := a.db.Exec("INSERT INTO settings(key, value) VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value = ?", key, value, value)
	return err
}

// ── 翻译 ──────────────────────────────────────────────────────────────────

// Translate 翻译单词，优先读缓存，缓存未命中则调 DeepL API
// 前端调用：const result = await window.go.main.App.Translate("hello")
func (a *App) Translate(word string) (*TranslateResult, error) {
	// 第一步：查数据库缓存
	// db.QueryRow 查询单行，Scan 把结果扫描到变量中
	// 类似前端的 db.query("SELECT ...").get()
	var result TranslateResult
	err := a.db.QueryRow(
		"SELECT word, translation, phonetic FROM words WHERE word = ?", word,
	).Scan(&result.Word, &result.Translation, &result.Phonetic)

	// err == nil 表示查到了
	if err == nil {
		result.Cached = true
		return &result, nil
	}

	// 第二步：缓存未命中，调 DeepL Free API
	apiKey, _ := a.GetSetting("deepl_api_key")
	if apiKey == "" {
		return nil, fmt.Errorf("请先在设置中配置 DeepL API Key")
	}

	// url.Values 构建 HTTP 表单数据
	// 类似前端的 new URLSearchParams()
	form := url.Values{}
	form.Set("auth_key", apiKey)
	form.Set("text", word)
	form.Set("target_lang", "ZH") // 翻译成中文

	// http.PostForm 发送 POST 请求（Content-Type: application/x-www-form-urlencoded）
	// 类似前端的 fetch(url, { method: "POST", body: form })
	resp, err := http.PostForm("https://api-free.deepl.com/v2/translate", form)
	if err != nil {
		return nil, err
	}
	// defer — 函数返回时执行，类似前端的 finally
	// 确保响应体被关闭，防止资源泄漏
	defer resp.Body.Close()

	// 读取响应体
	raw, _ := io.ReadAll(resp.Body)

	// json.Unmarshal 解析 JSON，类似前端的 JSON.parse()
	// 这里用匿名结构体（临时定义，不复用），映射 DeepL 的响应格式
	var deeplResp struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	if err := json.Unmarshal(raw, &deeplResp); err != nil || len(deeplResp.Translations) == 0 {
		return nil, fmt.Errorf("DeepL 响应解析失败: %s", raw)
	}

	// 取第一条翻译结果
	translation := deeplResp.Translations[0].Text

	// 写入数据库缓存，下次查同样的单词就不用再调 API
	// INSERT OR IGNORE — 如果 word 已存在（UNIQUE 约束），则跳过不报错
	a.db.Exec("INSERT OR IGNORE INTO words(word, translation) VALUES(?, ?)", word, translation)

	return &TranslateResult{Word: word, Translation: translation}, nil
}

// ── 单词本 ────────────────────────────────────────────────────────────────

// SaveWord 收藏单词到单词本
// 前端调用：window.go.main.App.SaveWord("hello", "你好", "/həˈloʊ/")
func (a *App) SaveWord(word, translation, phonetic string) error {
	// ON CONFLICT(word) DO NOTHING — 如果单词已存在，跳过插入
	// 这是 SQLite 的 upsert 语法，类似前端的 INSERT ... ON DUPLICATE KEY UPDATE
	_, err := a.db.Exec(
		"INSERT INTO words(word, translation, phonetic) VALUES(?, ?, ?) ON CONFLICT(word) DO NOTHING",
		word, translation, phonetic,
	)
	if err != nil {
		return err
	}

	// 查出单词的 ID，用于关联 word_book 表
	var wordID int64
	if err := a.db.QueryRow("SELECT id FROM words WHERE word = ?", word).Scan(&wordID); err != nil {
		return err
	}

	// 插入到单词本，INSERT OR IGNORE 防止重复收藏
	_, err = a.db.Exec("INSERT OR IGNORE INTO word_book(word_id) VALUES(?)", wordID)
	return err
}

// GetWordBook 获取单词本列表
// 前端调用：const list = await window.go.main.App.GetWordBook()
func (a *App) GetWordBook() ([]WordBookItem, error) {
	// db.Query 查询多行，返回 *sql.Rows 游标
	// JOIN 关联 words 和 word_book 两张表
	// COALESCE — 如果值为 NULL 则返回默认值，类似前端的 value || ''
	rows, err := a.db.Query(`
		SELECT w.word, w.translation, w.phonetic,
		       COALESCE(wb.note, ''), wb.reviewed,
		       COALESCE(wb.last_review, ''), wb.created_at
		FROM word_book wb
		JOIN words w ON w.id = wb.word_id
		ORDER BY wb.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	// defer rows.Close() — 函数返回时关闭游标，防止连接泄漏
	// 类似前端的 finally { cursor.close() }
	defer rows.Close()

	// 遍历结果集
	var items []WordBookItem
	for rows.Next() {
		var item WordBookItem
		// Scan 把当前行的列值扫描到 struct 字段
		// 字段顺序必须和 SELECT 的列顺序一致
		rows.Scan(&item.Word, &item.Translation, &item.Phonetic,
			&item.Note, &item.Reviewed, &item.LastReview, &item.CreatedAt)
		items = append(items, item)
	}
	return items, nil
}

// UpdateNote 更新单词备注
// 前端调用：window.go.main.App.UpdateNote("hello", "常用打招呼用语")
func (a *App) UpdateNote(word, note string) error {
	// 子查询 (SELECT id FROM words WHERE word = ?) 先查单词 ID
	// 再用 ID 定位 word_book 中的记录
	_, err := a.db.Exec(`
		UPDATE word_book SET note = ?
		WHERE word_id = (SELECT id FROM words WHERE word = ?)
	`, note, word)
	return err
}

// MarkReviewed 标记已复习，复习次数 +1
// 前端调用：window.go.main.App.MarkReviewed("hello")
func (a *App) MarkReviewed(word string) error {
	// reviewed = reviewed + 1 — 原子自增，不用先查再改
	// CURRENT_TIMESTAMP — 当前时间，SQLite 内置函数
	_, err := a.db.Exec(`
		UPDATE word_book
		SET reviewed = reviewed + 1, last_review = CURRENT_TIMESTAMP
		WHERE word_id = (SELECT id FROM words WHERE word = ?)
	`, word)
	return err
}

// RemoveWord 从单词本移除（不删除 words 表的缓存）
// 前端调用：window.go.main.App.RemoveWord("hello")
func (a *App) RemoveWord(word string) error {
	_, err := a.db.Exec(`
		DELETE FROM word_book
		WHERE word_id = (SELECT id FROM words WHERE word = ?)
	`, word)
	return err
}

// ── 文章 ──────────────────────────────────────────────────────────────────

// SaveArticle 保存文章
// 前端调用：window.go.main.App.SaveArticle("My Article", "Once upon a time...")
func (a *App) SaveArticle(title, content string) error {
	_, err := a.db.Exec(
		"INSERT INTO articles(title, content) VALUES(?, ?)", title, content,
	)
	return err
}

// GetArticles 获取文章列表
// 前端调用：const list = await window.go.main.App.GetArticles()
func (a *App) GetArticles() ([]Article, error) {
	rows, err := a.db.Query(
		"SELECT id, title, content, created_at FROM articles ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var art Article
		rows.Scan(&art.ID, &art.Title, &art.Content, &art.CreatedAt)
		articles = append(articles, art)
	}
	return articles, nil
}

// DeleteArticle 删除文章
// 前端调用：window.go.main.App.DeleteArticle(1)
func (a *App) DeleteArticle(id int64) error {
	_, err := a.db.Exec("DELETE FROM articles WHERE id = ?", id)
	return err
}
