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

	// 2. 初始化真实的大脑 (指向智谱 GLM-4.5，使用上一讲的 OpenAI 适配器)
	llmProvider := provider.NewZhipuOpenAIProvider("deepseek-chat")
	registry := tools.NewRegistry()

	// 挂载工具全家桶
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewBashTool(workDir))

	// 【新增挂载】
	registry.Register(tools.NewEditFileTool(workDir))

	// 实例化引擎，开启 EnableThinking = true (开启慢思考，促使模型一次性统筹规划)
	eng := engine.NewAgentEngine(llmProvider, registry, workDir, true)

	// 下发一个需要收集多源信息的任务
	prompt := `
    我当前目录下有 a.txt, b.txt, c.txt 三个文件。
    为了节省时间，请你同时一次性读取这三个文件，并将它们的内容综合起来，告诉我它们分别记录了什么领域的信息。
    `

	err := eng.Run(context.Background(), prompt)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}
}
