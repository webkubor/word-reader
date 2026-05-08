package main

import (
	"database/sql"
	"os"
	"path/filepath"

	// 下划线导入 — 只执行包的 init() 函数，不直接使用包里的东西
	// 这里导入 sqlite 驱动，让 sql.Open("sqlite", ...) 能识别这个驱动
	// 类似前端 import 'core-js' 只为 polyfill 副作用
	_ "modernc.org/sqlite"
)

// initDB 初始化数据库连接并建表
// 返回 (*sql.DB, error) 是 Go 的惯用模式：
//   - *sql.DB: 数据库连接池对象（成功时）
//   - error: 错误信息（失败时，成功时为 nil）
//
// 类似前端的 try/catch，但 Go 用返回值处理错误
func initDB() (*sql.DB, error) {
	// os.UserHomeDir() 获取用户主目录
	// macOS: /Users/xxx，Windows: C:\Users\xxx
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// filepath.Join 拼接路径，自动处理不同操作系统的分隔符
	// macOS: /Users/xxx/.word-reader
	// Windows: C:\Users\xxx\.word-reader
	dir := filepath.Join(home, ".word-reader")

	// os.MkdirAll 创建目录（包括所有父目录），0755 是权限位
	// 类似前端 mkdir -p，如果目录已存在不会报错
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// 数据库文件路径，如 /Users/xxx/.word-reader/data.db
	dbPath := filepath.Join(dir, "data.db")

	// sql.Open 只是创建连接池对象，不会真正连接
	// 第一个参数 "sqlite" 是驱动名，对应上面下划线导入的驱动
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// 执行建表迁移
	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

// migrate 数据库迁移 — 建表
// 类似前端的数据库 migration，保证表结构存在
// CREATE TABLE IF NOT EXISTS — 表不存在才创建，已存在则跳过
func migrate(db *sql.DB) error {
	schema := `
	-- words 表：单词缓存，翻译过的单词存在这里，避免重复调 API
	CREATE TABLE IF NOT EXISTS words (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,  -- 自增主键
		word        TEXT NOT NULL UNIQUE,                -- 英文单词，唯一约束
		translation TEXT NOT NULL,                       -- 中文翻译
		phonetic    TEXT DEFAULT '',                     -- 音标，默认空字符串
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP   -- 创建时间，自动填充
	);

	-- word_book 表：单词本，用户收藏的单词
	CREATE TABLE IF NOT EXISTS word_book (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,   -- 自增主键
		word_id     INTEGER NOT NULL REFERENCES words(id), -- 外键，关联 words 表
		note        TEXT DEFAULT '',                      -- 用户备注
		reviewed    INTEGER DEFAULT 0,                    -- 复习次数
		last_review DATETIME,                             -- 上次复习时间
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,   -- 收藏时间
		UNIQUE(word_id)                                   -- 每个单词只能收藏一次
	);

	-- settings 表：用户配置
	CREATE TABLE IF NOT EXISTS settings (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	-- articles 表：用户保存的英文文章
	CREATE TABLE IF NOT EXISTS articles (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,   -- 自增主键
		title      TEXT DEFAULT '',                      -- 文章标题
		content    TEXT NOT NULL,                        -- 文章内容
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP    -- 创建时间
	);
	`
	// db.Exec 执行 SQL 语句（不需要返回结果）
	// 类似前端的 db.execute(sql)
	_, err := db.Exec(schema)
	return err
}
