# 极简 Go Web Server 项目

## 目标
搭建一个极简、可运行的 Go 语言 Web Server 项目。

## 架构设计
- 单一 main package 入口
- 使用 Go 标准库 `net/http`
- 监听 8080 端口
- 提供根路径 `/` 返回 "Hello, Go!" 响应
- 极简主义：不用第三方库，不拆分多余文件

## 技术选型
- 语言：Go
- 标准库：net/http
- 无外部依赖

## 目录结构
```
.
├── main.go        # 主入口
├── go.mod         # Go module 定义
├── PLAN.md        # 架构文档
└── TODO.md        # 任务跟踪
```
