package cmd

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"timeLineGin/internal/route"
	"timeLineGin/pkg/config"
	"timeLineGin/pkg/logger"
	"timeLineGin/pkg/mysql"
	"timeLineGin/pkg/redis"
)

// 初始化 gin, zap 日志
func init() {
	// 初始化配置
	config.Initialize()
	// 初始化日志
	logger.Initialize(config.GetInstance().Logger)
	// 初始化数据库
	mysql.Initialize(config.GetInstance().DB)
	// 初始化 redis
	redis.Initialize(config.GetInstance().Redis)
}

type App struct{}

func NewApp() *App {
	return &App{}
}

// Run 启动服务
func (a *App) Run() {
	// 实例化 gin
	gin.SetMode(config.GetInstance().System.Env)
	app := gin.New()

	app.Use(logger.GinLogger(logger.GetInstance()))
	app.Use(gin.Recovery())

	// 实例化 http 服务
	srv := &http.Server{
		Handler: app,
		Addr:    config.GetInstance().System.Addr,
	}

	// 注册路由
	route.RegisterRoute(app)
	// 启动服务
	logger.GetInstance().Info("start server", zap.String("addr", srv.Addr))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.GetInstance().Error("server has benn stopped", zap.Error(err))
		}
	}()

	logger.GetInstance().Info("server start successful")

	// 监听外部停止信号
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	// 阻塞服务运行
	<-ch

	// 延时关闭
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.GetInstance().System.ShutdownWaitTime),
	)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.GetInstance().Error("服务关闭失败", zap.Error(err))
	}
	<-ctx.Done()
	logger.GetInstance().Info("服务关闭成功")
}
