package config

import (
	"fmt"
	"net"
)

type MainConfig struct {
	Serial  string `ymal:"serial"`
	Realm   string `ymal:"realm"`
	SipHost string `ymal:"sip-host"`
	SipPort int    `ymal:"sip-port"`

	EnableAuth bool   `ymal:"enable-auth"`
	Password   string `ymal:"password"`

	HttpServerPort int `ymal:"http-server-port"`
	APIPort        int `ymal:"api-port"`
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
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagUp == 0 {
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
