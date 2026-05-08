package main

// ── 数据模型 ──────────────────────────────────────────────────────────────
// Go 没有 class，用 struct（结构体）定义数据结构
// json:"xxx" 是标签（tag），控制 JSON 序列化时的字段名
// 类似前端的 TypeScript interface，但 Go 的 struct 只定义数据，方法另外写

// Word 单词缓存 — 对应数据库 words 表
// 存储翻译过的单词，避免重复调 API
type Word struct {
	ID          int64  `json:"id"`          // 主键，自增 ID
	Word        string `json:"word"`        // 英文单词
	Translation string `json:"translation"` // 中文翻译
	Phonetic    string `json:"phonetic"`    // 音标，如 /wɜːrd/
	CreatedAt   string `json:"created_at"`  // 创建时间
}

// WordBookItem 单词本条目 — 对应数据库 word_book 表
// 用户收藏的单词，带复习记录
type WordBookItem struct {
	ID          int64  `json:"id"`          // 主键
	Word        string `json:"word"`        // 英文单词
	Translation string `json:"translation"` // 中文翻译
	Phonetic    string `json:"phonetic"`    // 音标
	Note        string `json:"note"`        // 用户备注
	Reviewed    int    `json:"reviewed"`    // 复习次数
	LastReview  string `json:"last_review"` // 上次复习时间
	CreatedAt   string `json:"created_at"`  // 收藏时间
}

// Article 文章 — 对应数据库 articles 表
// 用户保存的英文文章
type Article struct {
	ID        int64  `json:"id"`         // 主键
	Title     string `json:"title"`      // 文章标题
	Content   string `json:"content"`    // 文章内容（Markdown 或纯文本）
	CreatedAt string `json:"created_at"` // 创建时间
}

// TranslateResult 翻译结果 — 返回给前端的数据结构
// 不是数据库表，是 API 响应的 DTO（Data Transfer Object）
type TranslateResult struct {
	Word        string `json:"word"`        // 翻译的单词
	Translation string `json:"translation"` // 翻译结果
	Phonetic    string `json:"phonetic"`    // 音标
	Cached      bool   `json:"cached"`      // 是否来自缓存（前端可据此显示"已缓存"标识）
}
