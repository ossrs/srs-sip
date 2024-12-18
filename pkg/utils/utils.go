package utils

import (
	"context"
	"crypto/rand"
	"flag"
	"math/big"
	"os"

	"github.com/ossrs/srs-sip/pkg/config"
)

func Parse(ctx context.Context) interface{} {
	fl := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	var conf config.MainConfig
	fl.StringVar(&conf.Serial, "serial", "34020000002000000001", "The serial number")
	fl.StringVar(&conf.Realm, "realm", "3402000000", "The realm")
	fl.StringVar(&conf.SipHost, "sip-host", "0.0.0.0", "The SIP host")
	fl.IntVar(&conf.SipPort, "sip-port", 5060, "The SIP port")
	fl.StringVar(&conf.Password, "password", "123456", "The password")
	fl.StringVar(&conf.MediaAddr, "media-addr", "127.0.0.1:1985", "The api address of media server. like: 127.0.0.1:1985")
	fl.IntVar(&conf.HttpServerPort, "http-server-port", 8888, "The port of http server")
	fl.IntVar(&conf.APIPort, "api-port", 2020, "The port of http api server")

	fl.Usage = func() {
		fl.PrintDefaults()
	}

	if err := fl.Parse(os.Args[1:]); err == flag.ErrHelp {
		os.Exit(0)
	}

	showHelp := conf.MediaAddr == ""
	if showHelp {
		fl.Usage()
		os.Exit(-1)
	}

	return &conf
}

func GenRandomNumber(n int) string {
	var result string
	for i := 0; i < n; i++ {
		randomDigit, _ := rand.Int(rand.Reader, big.NewInt(10))
		result += randomDigit.String()
	}
	return result
}

func CreateSSRC(isLive bool) string {
	ssrc := make([]byte, 10)
	if isLive {
		ssrc[0] = '0'
	} else {
		ssrc[0] = '1'
	}
	copy(ssrc[1:], GenRandomNumber(9))
	return string(ssrc)
}

// @see GB/T28181—2016 附录D 统一编码规则
func IsVideoChannel(channelID string) bool {
	deviceType := channelID[10:13]
	return deviceType == "131" || deviceType == "132"
}
