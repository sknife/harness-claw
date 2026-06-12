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

	// 挂载极简工具集
	registry.Register(tools.NewReadFileTool(workDir))
	registry.Register(tools.NewWriteFileTool(workDir))
	registry.Register(tools.NewBashTool(workDir))

	// 实例化核心引擎，关闭慢思考阶段，享受 YOLO 急速模式
	eng := engine.NewAgentEngine(llmProvider, registry, workDir, false)

	// 发起一个需要连贯物理动作的任务
	prompt := `
    请帮我执行以下操作：
    1. 用 bash 查看一下我当前电脑的 Go 版本。
    2. 帮我写一个简单的 helloworld.go 文件，输出 "Hello, go-tiny-claw!"。
    3. 用 bash 编译并运行这个 go 文件，确认它能正常工作。
    `

	err := eng.Run(context.Background(), prompt)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}
}
