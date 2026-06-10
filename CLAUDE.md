# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 常用命令

```bash
# 运行主程序（含 mock provider 的演示）
go run ./cmd/claw/

# 编译
go build ./...

# 运行测试
go test ./...

# 运行单个包的测试
go test ./internal/engine/

# 格式化代码
go fmt ./...

# 静态检查
go vet ./...
```

## 项目架构

**go-tiny-claw** 是一个极简 AI Agent 引擎骨架，演示了"两阶段慢思考 + 工具调用"的核心循环。

### 核心数据类型 (`internal/schema/message.go`)

所有组件之间通过以下类型通信：
- `Message`：上下文历史中的一条消息，角色为 `system` / `user` / `assistant`
- `ToolCall`：模型请求调用某工具时的描述（含 `Arguments json.RawMessage`，延迟解析）
- `ToolResult`：工具实际执行后的返回（含 `IsError` 标志用于错误自愈）
- `ToolDefinition`：向模型暴露的工具元信息（含 JSON Schema）

工具执行的观察结果以 `RoleUser` + `ToolCallID` 的方式追加到上下文历史，而非独立的 `tool` 角色。

### 两阶段 Agent 循环 (`internal/engine/loop.go`)

`AgentEngine.Run()` 实现一个"无限 turn"循环，每轮分两阶段：

1. **Phase 1 — Thinking**（仅 `EnableThinking=true` 时触发）：传入 `availableTools=nil`，剥夺模型的工具可见性，强制其输出纯文本规划；规划结果以 `assistant` 消息追加到上下文
2. **Phase 2 — Action**：恢复完整 `availableTools`，模型基于已有的 Thinking Trace 发起精准工具调用

循环终止条件：模型返回的响应中 `ToolCalls` 为空。

### 两个核心接口

- **`provider.LLMProvider`**：封装对 LLM 的单次推理调用，实现者负责将 `[]schema.Message` 和 `[]schema.ToolDefinition` 序列化为具体 API 请求
- **`tools.Registry`**：工具注册表，负责返回当前可用工具的 Schema（`GetAvailableTools`）以及分发执行（`Execute`）

### 扩展路径

- 接入真实 LLM：实现 `provider.LLMProvider` 接口替换 `cmd/claw/main.go` 中的 `mockProvider`
- 注册真实工具（如 `bash`、文件读写）：实现 `tools.Registry` 接口替换 `mockRegistry`
- 引擎本身无需修改，循环逻辑与具体 provider/tool 完全解耦
