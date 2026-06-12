# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 常用命令

```bash
# 运行主程序（需设置 ZHIPU_API_KEY 环境变量）
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

## 环境变量

- `ZHIPU_API_KEY`：必填，当前 `main.go` 通过 OpenAI 兼容协议指向 DeepSeek (`https://api.deepseek.com/`)。两个 Provider 实现都从此变量读取密钥，名称仅为历史遗留，与实际后端无关。
- 支持 `.env` 文件（由 `godotenv` 自动加载）。

## 项目架构

**go-tiny-claw** 是一个极简 AI Agent 引擎骨架，演示了"两阶段慢思考 + 工具调用"的核心循环。

Module 路径：`github.com/yourname/go-tiny-claw`

### 核心数据类型 (`internal/schema/message.go`)

所有组件之间通过以下类型通信：
- `Message`：上下文历史中的一条消息，角色为 `system` / `user` / `assistant`
- `ToolCall`：模型请求调用某工具时的描述（含 `Arguments json.RawMessage`，延迟解析）
- `ToolResult`：工具实际执行后的返回（含 `IsError` 标志用于错误自愈）
- `ToolDefinition`：向模型暴露的工具元信息（含 JSON Schema）

**关键约定**：工具执行的观察结果（Observation）以 `RoleUser` + `ToolCallID` 的方式追加到上下文，而非独立的 `tool` 角色。Provider 实现必须识别 `msg.ToolCallID != ""` 这一条件，将其翻译为对应 API 的 tool result 格式（Claude SDK 为 `ToolResultBlock`，OpenAI SDK 为 `ToolMessage`）。

### 两阶段 Agent 循环 (`internal/engine/loop.go`)

`AgentEngine.Run()` 实现一个"无限 turn"循环，每轮分两阶段：

1. **Phase 1 — Thinking**（仅 `EnableThinking=true` 时触发）：传入 `availableTools=nil`，剥夺模型的工具可见性，强制其输出纯文本规划；规划结果以 `assistant` 消息追加到上下文
2. **Phase 2 — Action**：恢复完整 `availableTools`，模型基于已有的 Thinking Trace 发起精准工具调用

循环终止条件：模型返回的响应中 `ToolCalls` 为空。

### 两个核心接口

- **`provider.LLMProvider`** (`internal/provider/interface.go`)：单次推理调用，实现者负责将 `[]schema.Message` 和 `[]schema.ToolDefinition` 序列化为具体 API 请求。当 `availableTools` 为 `nil` 时代表 Thinking 阶段，Provider 不得向 API 传递 Tools 字段。
- **`tools.Registry`** (`internal/tools/registry.go`)：工具注册表，`GetAvailableTools` 返回所有工具的 Schema，`Execute` 负责路由并执行工具调用，未知工具返回 `IsError=true` 让模型自愈。

### 添加新工具的模式

实现 `tools.BaseTool` 接口（三个方法：`Name() string`、`Definition() schema.ToolDefinition`、`Execute(ctx, json.RawMessage) (string, error)`），参考 `internal/tools/read_file.go`：
- `workDir` 注入以限制文件操作边界
- `Execute` 内部完成 JSON 反序列化，错误直接返回给 Registry 转发给模型
- `read_file` 工具对输出做了 8000 字节截断保护，防止超大文件撑爆 Context

### 现有 Provider 实现

| 文件 | 构造函数 | SDK | 实际后端（当前配置） |
|------|----------|-----|-----|
| `openai.go` | `NewZhipuOpenAIProvider` | openai-go v3 | DeepSeek (`api.deepseek.com`) |
| `claude.go` | `NewZhipuClaudeProvider` | anthropic-sdk-go | 智谱 (`open.bigmodel.cn`) |

`main.go` 目前使用 `NewZhipuOpenAIProvider("deepseek-chat")`。

### 扩展路径

- 接入新 LLM：实现 `provider.LLMProvider` 接口，在 `cmd/claw/main.go` 中替换即可
- 注册新工具：实现 `tools.BaseTool` 接口，调用 `registry.Register(tool)` 挂载
- 引擎本身无需修改，循环逻辑与具体 provider/tool 完全解耦
