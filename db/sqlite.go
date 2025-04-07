package db

import (
	"database/sql"
	"goblog/models"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// SQLiteStore SQLite存储实现
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore 创建新的SQLite存储
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	log.Printf("尝试打开数据库: %s", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("打开数据库失败: %v", err)
		return nil, err
	}

	log.Println("尝试Ping数据库连接...")
	if err := db.Ping(); err != nil {
		log.Printf("Ping数据库失败: %v", err)
		return nil, err
	}
	log.Println("数据库连接成功")

	store := &SQLiteStore{db: db}

	log.Println("开始初始化数据库结构...")
	if err := store.Initialize(); err != nil {
		log.Printf("初始化数据库结构失败: %v", err)
		return nil, err
	}
	log.Println("数据库结构初始化成功")

	return store, nil
}

// Initialize 初始化数据库表
func (s *SQLiteStore) Initialize() error {
	log.Println("开始初始化数据库表...")

	// 创建用户表
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`)
	if err != nil {
		log.Printf("创建用户表失败: %v", err)
		return err
	}
	log.Println("用户表创建成功或已存在")

	// 创建文章表
	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (id)
	)`)
	if err != nil {
		log.Printf("创建文章表失败: %v", err)
		return err
	}
	log.Println("文章表创建成功或已存在")

	log.Println("数据库初始化完成")
	return nil
}

// Close 关闭数据库连接
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// FindAllPosts 查找所有文章
func (s *SQLiteStore) FindAllPosts() ([]*models.Post, error) {
	log.Println("正在查询所有文章...")

	// 首先检查posts表是否存在
	var tableName string
	err := s.db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='posts'`).Scan(&tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("posts表不存在，返回空列表")
			return []*models.Post{}, nil
		}
		log.Printf("检查posts表是否存在时出错: %v", err)
		return nil, err
	}

	// 检查是否有任何文章
	var count int
	err = s.db.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&count)
	if err != nil {
		log.Printf("计算文章数量时出错: %v", err)
		return nil, err
	}

	if count == 0 {
		log.Println("没有任何文章，返回空列表")
		return []*models.Post{}, nil
	}

	// 继续原来的查询
	rows, err := s.db.Query(`
		SELECT p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at,
			   u.id, u.username, u.email, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		log.Printf("查询文章失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		var user models.User
		var postCreatedAt, postUpdatedAt, userCreatedAt, userUpdatedAt string

		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.UserID, &postCreatedAt, &postUpdatedAt,
			&user.ID, &user.Username, &user.Email, &userCreatedAt, &userUpdatedAt,
		)
		if err != nil {
			log.Printf("扫描文章行失败: %v", err)
			return nil, err
		}

		post.CreatedAt, _ = time.Parse(time.RFC3339, postCreatedAt)
		post.UpdatedAt, _ = time.Parse(time.RFC3339, postUpdatedAt)
		user.CreatedAt, _ = time.Parse(time.RFC3339, userCreatedAt)
		user.UpdatedAt, _ = time.Parse(time.RFC3339, userUpdatedAt)
		post.User = &user

		posts = append(posts, &post)
	}

	// 检查遍历过程中是否有错误
	if err = rows.Err(); err != nil {
		log.Printf("遍历结果集时出错: %v", err)
		return nil, err
	}

	log.Printf("查询到 %d 篇文章", len(posts))
	return posts, nil
}

// FindPostByID 根据ID查找文章
func (s *SQLiteStore) FindPostByID(id int) (*models.Post, error) {
	var post models.Post
	var user models.User
	var postCreatedAt, postUpdatedAt, userCreatedAt, userUpdatedAt string

	err := s.db.QueryRow(`
		SELECT p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at,
			   u.id, u.username, u.email, u.created_at, u.updated_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, id).Scan(
		&post.ID, &post.Title, &post.Content, &post.UserID, &postCreatedAt, &postUpdatedAt,
		&user.ID, &user.Username, &user.Email, &userCreatedAt, &userUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	post.CreatedAt, _ = time.Parse(time.RFC3339, postCreatedAt)
	post.UpdatedAt, _ = time.Parse(time.RFC3339, postUpdatedAt)
	user.CreatedAt, _ = time.Parse(time.RFC3339, userCreatedAt)
	user.UpdatedAt, _ = time.Parse(time.RFC3339, userUpdatedAt)
	post.User = &user

	return &post, nil
}

// CreatePost 创建文章
func (s *SQLiteStore) CreatePost(post *models.Post) error {
	now := time.Now().Format(time.RFC3339)
	result, err := s.db.Exec(`
		INSERT INTO posts (title, content, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, post.Title, post.Content, post.UserID, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.ID = int(id)
	post.CreatedAt, _ = time.Parse(time.RFC3339, now)
	post.UpdatedAt = post.CreatedAt

	return nil
}

// UpdatePost 更新文章
func (s *SQLiteStore) UpdatePost(post *models.Post) error {
	now := time.Now().Format(time.RFC3339)
	_, err := s.db.Exec(`
		UPDATE posts
		SET title = ?, content = ?, updated_at = ?
		WHERE id = ?
	`, post.Title, post.Content, now, post.ID)
	if err != nil {
		return err
	}

	post.UpdatedAt, _ = time.Parse(time.RFC3339, now)
	return nil
}

// DeletePost 删除文章
func (s *SQLiteStore) DeletePost(id int) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}

// FindUserByID 根据ID查找用户
func (s *SQLiteStore) FindUserByID(id int) (*models.User, error) {
	var user models.User
	var createdAt, updatedAt string

	err := s.db.QueryRow(`
		SELECT id, username, email, password, created_at, updated_at
		FROM users WHERE id = ?
	`, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &user, nil
}

// FindUserByUsername 根据用户名查找用户
func (s *SQLiteStore) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	var createdAt, updatedAt string

	err := s.db.QueryRow(`
		SELECT id, username, email, password, created_at, updated_at
		FROM users WHERE username = ?
	`, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &user, nil
}

// FindUserByEmail 根据邮箱查找用户
func (s *SQLiteStore) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	var createdAt, updatedAt string

	err := s.db.QueryRow(`
		SELECT id, username, email, password, created_at, updated_at
		FROM users WHERE email = ?
	`, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &user, nil
}

// CreateUser 创建用户
func (s *SQLiteStore) CreateUser(user *models.User) error {
	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	result, err := s.db.Exec(`
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, user.Username, user.Email, string(hashedPassword), now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	user.CreatedAt, _ = time.Parse(time.RFC3339, now)
	user.UpdatedAt = user.CreatedAt

	return nil
}

// UpdateUser 更新用户
func (s *SQLiteStore) UpdateUser(user *models.User) error {
	now := time.Now().Format(time.RFC3339)
	_, err := s.db.Exec(`
		UPDATE users
		SET username = ?, email = ?, updated_at = ?
		WHERE id = ?
	`, user.Username, user.Email, now, user.ID)
	if err != nil {
		return err
	}

	user.UpdatedAt, _ = time.Parse(time.RFC3339, now)
	return nil
}

// DeleteUser 删除用户
func (s *SQLiteStore) DeleteUser(id int) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// Authenticate 认证用户
func (s *SQLiteStore) Authenticate(username, password string) (*models.User, error) {
	log.Printf("尝试验证用户: %s", username)
	user, err := s.FindUserByUsername(username)
	if err != nil {
		log.Printf("查找用户错误: %v", err)
		return nil, err
	}

	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("密码比较错误: %v", err)
		return nil, err
	}

	log.Printf("用户验证成功: %s", user.Username)
	return user, nil
}
