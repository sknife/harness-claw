package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yourname/go-tiny-claw/internal/engine"
	"github.com/yourname/go-tiny-claw/internal/provider"
	"github.com/yourname/go-tiny-claw/internal/tools"
)

func main() {
	_ = godotenv.Load()

	// 确保已设置 ZHIPU_API_KEY
	if os.Getenv("ZHIPU_API_KEY") == "" {
		log.Fatal("请先导出 ZHIPU_API_KEY 环境变量")
	}

	// 1. 获取工作区物理边界
	workDir, _ := os.Getwd()
	workDir += "\\workspace"

	// 2. 初始化真实的大脑 (指向智谱 GLM-4.5，使用上一讲的 OpenAI 适配器)
	llmProvider := provider.NewZhipuOpenAIProvider("deepseek-chat")

	registry := tools.NewRegistry()
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewBashTool(workDir))
	registry.Register(tools.NewEditFileTool(workDir))

	// 实例化引擎，开启慢思考
	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)
	// 【注入新实现的终端输出器】
	reporter := engine.NewTerminalReporter()

	prompt := `
    我需要在当前目录下新建一个 ping.go，提供一个简单的 http ping 接口。
    写完之后，帮我把代码用 git 提交一下。
    `

	err := eng.Run(context.Background(), prompt, reporter)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}
}
