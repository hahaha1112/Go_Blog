# GoBlog - 基于Go的简易博客系统

这是一个使用Go语言开发的简单博客系统，支持用户注册登录、文章发布、编辑和管理等基本功能。

## 功能特点

- 用户管理：注册、登录、退出
- 文章管理：创建、查看、编辑、删除
- 响应式设计：适配不同设备屏幕大小
- SQLite数据库：轻量级存储解决方案

## 技术栈

- Go 1.20+
- 原生Go标准库 (net/http)
- SQLite3 数据库
- gorilla/sessions 会话管理
- golang.org/x/crypto 密码处理

## 运行环境要求

- Go 1.20 或更高版本
- 支持SQLite的操作系统（Windows、Linux、MacOS）

## 安装和运行

1. 克隆项目到本地：
   ```
   git clone [项目地址]
   cd goblog
   ```

2. 安装依赖：
   ```
   go mod tidy
   ```

3. 运行项目：
   ```
   go run main.go
   ```

4. 打开浏览器，访问 `http://localhost:8080`

## 项目结构

```
goblog/
├── config/         // 配置相关
├── controllers/    // 控制器
├── db/             // 数据库访问
├── middleware/     // 中间件
├── models/         // 数据模型
├── public/         // 静态资源
│   ├── css/        // 样式文件
│   └── js/         // JavaScript文件
├── router/         // 路由配置
├── templates/      // HTML模板
│   ├── posts/      // 文章相关模板
│   └── users/      // 用户相关模板
├── utils/          // 工具函数
├── main.go         // 主入口文件
├── go.mod          // Go模块文件
└── README.md       // 项目说明
```

## 配置说明

系统会自动创建 `config.json` 配置文件，您可以根据需要修改以下配置：

```json
{
  "server": {
    "port": 8080,
    "readTimeout": 60,
    "writeTimeout": 60
  },
  "database": {
    "type": "sqlite3",
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "password",
    "dbname": "goblog"
  }
}
```

## 后续开发计划

- 添加评论功能
- 添加标签和分类
- 增加文章搜索功能
- 支持Markdown编辑器
- 添加管理员后台
- 优化移动端体验

## 贡献指南

欢迎贡献代码或提出建议！请通过GitHub Issue或Pull Request参与项目改进。

## 许可证

本项目采用MIT许可证 - 详见 LICENSE 文件 