package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ossrs/go-oryx-lib/logger"
	"github.com/ossrs/srs-sip/pkg/api"
	"github.com/ossrs/srs-sip/pkg/config"
	"github.com/ossrs/srs-sip/pkg/service"
)

func WaitTerminationSignal(cancel context.CancelFunc) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)
	<-sigc
	cancel()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	conf, err := config.LoadConfig("conf/config.yaml")
	if err != nil {
		logger.E(nil, "load config failed: %v", err)
		return
	}

	sipSvr, err := service.NewService(ctx, conf)
	if err != nil {
		logger.Ef("create service failed. err is %v", err.Error())
		return
	}

	if err := sipSvr.Start(); err != nil {
		logger.Ef("start sip service failed. err is %v", err.Error())
		return
	}

	// 创建主路由
	router := mux.NewRouter().StrictSlash(true)

	// CORS配置
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// 设置API路由 - 需要在静态文件路由之前设置
	apiSvr, err := api.NewHttpApiServer(conf, sipSvr)
	if err != nil {
		logger.Ef("create http service failed. err is %v", err.Error())
		return
	}
	apiSvr.Start(router)

	// 使用配置中指定的目录，如果不存在则尝试备选目录
	targetDir := conf.Http.Dir
	if _, err := os.Stat(path.Join(targetDir, "index.html")); err != nil {
		backupDirs := []string{"./html", "../web/NextGB/dist"}
		for _, dir := range backupDirs {
			if _, err := os.Stat(path.Join(dir, "index.html")); err == nil {
				targetDir = dir
				break
			}
		}
	}
	if targetDir == "" {
		logger.Ef(ctx, "index.html not found")
		return
	}

	// 创建文件服务器
	fs := http.FileServer(http.Dir(targetDir))

	// 添加静态文件处理 - 使用NotFoundHandler来处理未匹配的路由
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Tf(context.Background(), "Handling request: %s", r.URL.Path)

		// 检查请求的文件是否存在
		filePath := path.Join(targetDir, r.URL.Path)
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			// 如果文件不存在，返回 index.html
			r.URL.Path = "/"
		}
		fs.ServeHTTP(w, r)
	})

	// 启动合并后的HTTP服务
	go func() {
		httpPort := strconv.Itoa(conf.Http.Port)
		handler := handlers.CORS(headers, methods, origins)(router)
		server := &http.Server{
			Addr:              ":" + httpPort,
			Handler:           handler,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		}
		logger.Tf(ctx, "http server listen on %s, home is %v", httpPort, targetDir)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Ef(ctx, "listen on %s failed", httpPort)
		}
	}()

	WaitTerminationSignal(cancel)

	sipSvr.Stop()
}
