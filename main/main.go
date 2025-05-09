package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ossrs/srs-sip/pkg/api"
	"github.com/ossrs/srs-sip/pkg/config"
	"github.com/ossrs/srs-sip/pkg/service"
	"github.com/ossrs/srs-sip/pkg/utils"
)

func WaitTerminationSignal(cancel context.CancelFunc) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)
	<-sigc
	cancel()
}

func main() {
	configPath := flag.String("c", "", "配置文件路径")
	flag.Parse()

	if *configPath == "" {
		slog.Error("error: specify the config file path, like: ./srs-sip -c conf/config.yaml")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		slog.Error("load config failed", "error", err)
		return
	}

	if err := utils.SetupLogger(conf.Common.LogLevel, conf.Common.LogFile); err != nil {
		slog.Error("setup logger failed", "error", err)
		return
	}

	slog.Info("*****************************************************")
	slog.Info("          ☆☆☆ 欢迎使用 SRS-SIP 服务 ☆☆☆")
	slog.Info("*****************************************************")
	slog.Info("srs-sip service starting", "config", *configPath, "log_file", conf.Common.LogFile)

	sipSvr, err := service.NewService(ctx, conf)
	if err != nil {
		slog.Error("create service failed", "error", err.Error())
		return
	}

	if err := sipSvr.Start(); err != nil {
		slog.Error("start sip service failed", "error", err.Error())
		return
	}

	// 创建主路由
	router := mux.NewRouter().StrictSlash(true)

	// CORS配置
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// 创建文件服务器
	fs := http.FileServer(http.Dir(conf.Http.Dir))

	// 先注册API路由
	apiSvr, err := api.NewHttpApiServer(conf, sipSvr)
	if err != nil {
		slog.Error("create http service failed", "error", err.Error())
		return
	}
	apiSvr.Start(router)

	// 添加静态文件处理 - 使用PathPrefix处理所有非API请求
	router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 如果是API路径，直接返回404
		if strings.HasPrefix(r.URL.Path, "/srs-sip/v1/") {
			slog.Info("api path not found", "path", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		// 检查请求的文件是否存在
		filePath := path.Join(conf.Http.Dir, r.URL.Path)
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			// 如果文件不存在，返回 index.html
			slog.Info("file not found, redirect to index", "path", r.URL.Path)
			r.URL.Path = "/"
		}
		fs.ServeHTTP(w, r)
	}))

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
		slog.Info("http server listen", "port", httpPort, "home", conf.Http.Dir)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen failed", "port", httpPort, "error", err)
		}
	}()

	WaitTerminationSignal(cancel)

	sipSvr.Stop()
	slog.Info("srs-sip service stopped")
}
