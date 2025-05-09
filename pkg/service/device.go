package service

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ossrs/srs-sip/pkg/models"
)

type DeviceInfo struct {
	DeviceID          string    `json:"device_id"`
	SourceAddr        string    `json:"source_addr"`
	NetworkType       string    `json:"network_type"`
	ChannelMap        sync.Map  `json:"-"`
	Online            bool      `json:"online"`
	HeartBeatInterval int       `json:"heart_beat_interval"`
	HeartBeatCount    int       `json:"heart_beat_count"`
	lastHeartbeat     time.Time `json:"-"`
}

const (
	DefaultHeartbeatInterval = 60 * time.Second // 心跳检查间隔时间
	DefaultHeartbeatCount    = 3                // 心跳检查次数
)

type deviceManager struct {
	devices          sync.Map
	heartbeatChecker *time.Ticker  // 心跳检查定时器
	stopChan         chan struct{} // 停止信号通道
}

var instance *deviceManager
var once sync.Once

func GetDeviceManager() *deviceManager {
	once.Do(func() {
		instance = &deviceManager{
			devices:  sync.Map{},
			stopChan: make(chan struct{}),
		}
		// 启动心跳检查
		instance.startHeartbeatChecker()
	})
	return instance
}

// 启动心跳检查器
func (dm *deviceManager) startHeartbeatChecker() {
	dm.heartbeatChecker = time.NewTicker(3 * time.Second) // 每3秒检查一次
	go func() {
		for {
			select {
			case <-dm.heartbeatChecker.C:
				dm.checkHeartbeats()
			case <-dm.stopChan:
				dm.heartbeatChecker.Stop()
				return
			}
		}
	}()
}

// 停止心跳检查器
func (dm *deviceManager) stopHeartbeatChecker() {
	close(dm.stopChan)
}

// 检查所有设备的心跳状态
func (dm *deviceManager) checkHeartbeats() {
	now := time.Now()
	dm.devices.Range(func(key, value interface{}) bool {
		device := value.(*DeviceInfo)

		if device.HeartBeatInterval == 0 {
			device.HeartBeatInterval = int(DefaultHeartbeatInterval)
		}
		if device.HeartBeatCount == 0 {
			device.HeartBeatCount = DefaultHeartbeatCount
		}

		// 如果最后心跳时间超过超时时间，则将设备所有通道状态设置为离线
		if now.Sub(device.lastHeartbeat) > time.Duration(device.HeartBeatInterval*device.HeartBeatCount)*time.Second {
			isOffline := false
			device.ChannelMap.Range(func(key, value interface{}) bool {
				channel := value.(models.ChannelInfo)
				if channel.Status != models.ChannelStatus("OFF") {
					isOffline = true
					channel.Status = models.ChannelStatus("OFF")
					device.ChannelMap.Store(key, channel)
				}
				return true
			})

			if isOffline {
				device.SourceAddr = ""
				device.Online = false
				dm.devices.Store(key, device)
				slog.Warn("Device is offline due to heartbeat timeout",
					"device_id", device.DeviceID,
					"heartbeat_interval", device.HeartBeatInterval,
					"heartbeat_count", device.HeartBeatCount)
			}
		}
		return true
	})
}

func (dm *deviceManager) UpdateDeviceHeartbeat(id string) {
	if device, ok := dm.GetDevice(id); ok {
		device.lastHeartbeat = time.Now()

		// 检查是否需要将通道状态设置为在线
		isUpdated := false
		device.ChannelMap.Range(func(key, value interface{}) bool {
			channel := value.(models.ChannelInfo)
			if channel.Status != models.ChannelStatus("ON") {
				isUpdated = true
				channel.Status = models.ChannelStatus("ON")
				device.ChannelMap.Store(key, channel)
			}
			return true
		})

		if isUpdated {
			device.Online = true
			dm.devices.Store(id, device)
		}
	}
}

func (dm *deviceManager) AddDevice(id string, info *DeviceInfo) {
	// 设置初始心跳时间
	info.lastHeartbeat = time.Now()
	info.Online = true

	channel := models.ChannelInfo{
		DeviceID: id,
		ParentID: id,
		Name:     id,
		Status:   models.ChannelStatus("ON"),
	}
	info.ChannelMap.Store(channel.DeviceID, channel)
	dm.devices.Store(id, info)
}

func (dm *deviceManager) RemoveDevice(id string) {
	dm.devices.Delete(id)
}

func (dm *deviceManager) GetDevices() []*DeviceInfo {
	list := make([]*DeviceInfo, 0)
	dm.devices.Range(func(key, value interface{}) bool {
		list = append(list, value.(*DeviceInfo))
		return true
	})
	return list
}

func (dm *deviceManager) GetDevice(id string) (*DeviceInfo, bool) {
	v, ok := dm.devices.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*DeviceInfo), true
}

// UpdateDevice 更新设备信息
func (dm *deviceManager) UpdateDevice(id string, device *DeviceInfo) {
	device.lastHeartbeat = time.Now()
	device.Online = true
	dm.devices.Store(id, device)
}

// UpdateDeviceConfig 更新设备配置信息
func (dm *deviceManager) UpdateDeviceConfig(deviceID string, basicParam *models.BasicParam) {
	device, ok := dm.GetDevice(deviceID)
	if !ok {
		return
	}

	if basicParam != nil {
		// 更新设备心跳相关配置
		if basicParam.HeartBeatInterval > 0 {
			device.HeartBeatInterval = basicParam.HeartBeatInterval
		}
		if basicParam.HeartBeatCount > 0 {
			device.HeartBeatCount = basicParam.HeartBeatCount
		}

		dm.devices.Store(deviceID, device)
	}
}

// ChannelParser defines interface for different manufacturer's channel parsing
type ChannelParser interface {
	ParseChannels(list ...models.ChannelInfo) ([]models.ChannelInfo, error)
}

// channelParserRegistry manages registration and lookup of manufacturer-specific parsers
type channelParserRegistry struct {
	parsers map[string]ChannelParser
	mu      sync.RWMutex
}

var (
	parserRegistry = &channelParserRegistry{
		parsers: make(map[string]ChannelParser),
	}
)

// RegisterParser registers a parser for a specific manufacturer
func (r *channelParserRegistry) RegisterParser(manufacturer string, parser ChannelParser) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parsers[manufacturer] = parser
}

// GetParser retrieves parser for a specific manufacturer
func (r *channelParserRegistry) GetParser(manufacturer string) (ChannelParser, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	parser, ok := r.parsers[manufacturer]
	return parser, ok
}

// UpdateChannels updates device channel information
func (dm *deviceManager) UpdateChannels(deviceID string, list ...models.ChannelInfo) error {
	device, ok := dm.GetDevice(deviceID)
	if !ok {
		return fmt.Errorf("device not found: %s", deviceID)
	}

	// clear ChannelMap
	device.ChannelMap.Range(func(key, value interface{}) bool {
		device.ChannelMap.Delete(key)
		return true
	})

	parser, ok := parserRegistry.GetParser(list[0].Manufacturer)
	if !ok {
		return fmt.Errorf("no parser found for manufacturer: %s", list[0].Manufacturer)
	}

	channels, err := parser.ParseChannels(list...)
	if err != nil {
		return fmt.Errorf("failed to parse channels: %v", err)
	}

	for _, channel := range channels {
		device.ChannelMap.Store(channel.DeviceID, channel)
	}
	dm.devices.Store(deviceID, device)
	return nil
}

func (dm *deviceManager) ApiGetChannelByDeviceId(deviceID string) []models.ChannelInfo {
	device, ok := dm.GetDevice(deviceID)
	if !ok {
		return nil
	}

	channels := make([]models.ChannelInfo, 0)
	device.ChannelMap.Range(func(key, value interface{}) bool {
		channels = append(channels, value.(models.ChannelInfo))
		return true
	})
	return channels
}

func (dm *deviceManager) GetAllVideoChannels() []models.ChannelInfo {
	channels := make([]models.ChannelInfo, 0)
	dm.devices.Range(func(key, value interface{}) bool {
		device := value.(*DeviceInfo)
		device.ChannelMap.Range(func(key, value interface{}) bool {
			channels = append(channels, value.(models.ChannelInfo))
			return true
		})
		return true
	})
	return channels
}

func (dm *deviceManager) GetDeviceInfoByChannel(channelID string) (*DeviceInfo, bool) {
	var device *DeviceInfo
	found := false
	dm.devices.Range(func(key, value interface{}) bool {
		d := value.(*DeviceInfo)
		_, ok := d.ChannelMap.Load(channelID)
		if ok {
			device = d
			found = true
			return false
		}
		return true
	})
	return device, found
}

// Hikvision channel parser implementation
type HikvisionParser struct{}

func (p *HikvisionParser) ParseChannels(list ...models.ChannelInfo) ([]models.ChannelInfo, error) {
	return list, nil
}

// Dahua channel parser implementation
type DahuaParser struct{}

func (p *DahuaParser) ParseChannels(list ...models.ChannelInfo) ([]models.ChannelInfo, error) {
	return list, nil
}

// Uniview channel parser implementation
type UniviewParser struct{}

func (p *UniviewParser) ParseChannels(list ...models.ChannelInfo) ([]models.ChannelInfo, error) {
	videoChannels := make([]models.ChannelInfo, 0)
	for _, channel := range list {
		// 只有Parental为1的通道，才是视频通道
		if channel.Parental == 1 {
			videoChannels = append(videoChannels, channel)
		}
	}
	return videoChannels, nil
}

func init() {
	parserRegistry.RegisterParser("Hikvision", &HikvisionParser{})
	parserRegistry.RegisterParser("DAHUA", &DahuaParser{})
	parserRegistry.RegisterParser("UNIVIEW", &UniviewParser{})
}
