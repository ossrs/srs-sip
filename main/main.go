package main

import (
	"context"
	"flag"
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
	// 定义配置文件路径参数
	configPath := flag.String("c", "", "配置文件路径")
	flag.Parse()

	if *configPath == "" {
		logger.E(nil, "错误: 通过 -c 参数指定配置文件路径，比如：./srs-sip -c conf/config.yaml")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	conf, err := config.LoadConfig(*configPath)
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

	// 创建文件服务器
	fs := http.FileServer(http.Dir(conf.Http.Dir))

	// 先注册API路由
	apiSvr, err := api.NewHttpApiServer(conf, sipSvr)
	if err != nil {
		logger.Ef("create http service failed. err is %v", err.Error())
		return
	}
	apiSvr.Start(router)

	// 添加静态文件处理 - 使用PathPrefix处理所有非API请求
	router.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 如果是API路径，直接返回404
		if strings.HasPrefix(r.URL.Path, "/srs-sip/v1/") {
			http.NotFound(w, r)
			return
		}

		// 检查请求的文件是否存在
		filePath := path.Join(conf.Http.Dir, r.URL.Path)
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			// 如果文件不存在，返回 index.html
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
		logger.Tf(ctx, "http server listen on %s, home is %v", httpPort, conf.Http.Dir)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Ef(ctx, "listen on %s failed", httpPort)
		}
	}()

	WaitTerminationSignal(cancel)

	sipSvr.Stop()
}
