package main

import "timeLineGin/cmd"

import (
	// 注册 logic 实例
	_ "timeLineGin/internal/logic"
)

func main() {
	cmd.NewApp().Run()
}
