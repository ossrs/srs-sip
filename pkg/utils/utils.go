package utils

import (
	"crypto/rand"
	"math/big"
)

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

// GetSessionName 根据播放类型返回会话名称
func GetSessionName(playType int) string {
	switch playType {
	case 1:
		return "Playback"
	case 2:
		return "Download"
	case 3:
		return "Talk"
	default:
		return "Play"
	}
}
