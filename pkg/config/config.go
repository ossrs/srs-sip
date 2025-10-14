package config

import (
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

// 通用配置
type CommonConfig struct {
	LogLevel string `yaml:"log-level"`
	LogFile  string `yaml:"log-file"`
}

// GB28181配置
type GB28181AuthConfig struct {
	Enable   bool   `yaml:"enable"`
	Password string `yaml:"password"`
}

type GB28181Config struct {
	Serial string            `yaml:"serial"`
	Realm  string            `yaml:"realm"`
	Host   string            `yaml:"host"`
	Port   int               `yaml:"port"`
	Auth   GB28181AuthConfig `yaml:"auth"`
}

// HTTP服务配置
type HttpConfig struct {
	Port int    `yaml:"listen"`
	Dir  string `yaml:"dir"`
}

// 主配置结构
type MainConfig struct {
	Common  CommonConfig  `yaml:"common"`
	GB28181 GB28181Config `yaml:"gb28181"`
	Http    HttpConfig    `yaml:"http"`
}

// 获取默认配置
func DefaultConfig() *MainConfig {
	return &MainConfig{
		Common: CommonConfig{
			LogLevel: "info",
			LogFile:  "app.log",
		},
		GB28181: GB28181Config{
			Serial: "34020000002000000001",
			Realm:  "3402000000",
			Host:   "0.0.0.0",
			Port:   5060,
			Auth: GB28181AuthConfig{
				Enable:   false,
				Password: "123456",
			},
		},
		Http: HttpConfig{
			Port: 8025,
			Dir:  "./html",
		},
	}
}

func LoadConfig(filename string) (*MainConfig, error) {
	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %v", err)
	}

	var config MainConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse config file failed: %v", err)
	}

	return &config, nil
}

func GetLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", nil
	}
	type Iface struct {
		Name string
		Addr net.IP
	}
	var candidates []Iface
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagUp == 0 {
			continue
		}
		if ifc.Flags&(net.FlagPointToPoint|net.FlagLoopback) != 0 {
			continue
		}
		addrs, err := ifc.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				candidates = append(candidates, Iface{
					Name: ifc.Name, Addr: ip4,
				})
				//logger.Tf("considering interface", "iface", ifc.Name, "ip", ip4)
			}
		}
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("No local IP found")
	}
	return candidates[0].Addr.String(), nil
}
