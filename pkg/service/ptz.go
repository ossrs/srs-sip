package service

import "fmt"

var (
	ptzCmdMap = map[string]uint8{
		"stop":      0,
		"right":     1,
		"left":      2,
		"down":      4,
		"downright": 5,
		"downleft":  6,
		"up":        8,
		"upright":   9,
		"upleft":    10,
		"zoomin":    16,
		"zoomout":   32,
	}

	ptzSpeedMap = map[string]uint8{
		"1":  25,
		"2":  50,
		"3":  75,
		"4":  100,
		"5":  125,
		"6":  150,
		"7":  175,
		"8":  200,
		"9":  225,
		"10": 255,
	}

	defaultSpeed uint8 = 125
)

func getPTZSpeed(speed string) uint8 {
	if v, ok := ptzSpeedMap[speed]; ok {
		return v
	}
	return defaultSpeed
}

func toPTZCmd(cmdName, speed string) (string, error) {
	cmdCode, ok := ptzCmdMap[cmdName]
	if !ok {
		return "", fmt.Errorf("invalid ptz command: %q", cmdName)
	}

	speedValue := getPTZSpeed(speed)

	var horizontalSpeed, verticalSpeed, zSpeed uint8

	switch cmdName {
	case "left", "right":
		horizontalSpeed = speedValue
		verticalSpeed = 0
	case "up", "down":
		verticalSpeed = speedValue
		horizontalSpeed = 0
	case "upleft", "upright", "downleft", "downright":
		verticalSpeed = speedValue
		horizontalSpeed = speedValue
	case "zoomin", "zoomout":
		zSpeed = speedValue << 4 // zoom速度在高4位
	default:
		horizontalSpeed = 0
		verticalSpeed = 0
		zSpeed = 0
	}

	sum := uint16(0xA5) + uint16(0x0F) + uint16(0x01) + uint16(cmdCode) + uint16(horizontalSpeed) + uint16(verticalSpeed) + uint16(zSpeed)
	checksum := uint8(sum % 256)

	return fmt.Sprintf("A50F01%02X%02X%02X%02X%02X",
		cmdCode,
		horizontalSpeed,
		verticalSpeed,
		zSpeed,
		checksum,
	), nil
}
