package service

import (
	"fmt"
	"sync"

	"github.com/ossrs/srs-sip/pkg/models"
)

type DeviceInfo struct {
	DeviceID    string   `json:"device_id"`
	SourceAddr  string   `json:"source_addr"`
	NetworkType string   `json:"network_type"`
	ChannelMap  sync.Map `json:"-"`
}

type deviceManager struct {
	devices sync.Map
}

var instance *deviceManager
var once sync.Once

func GetDeviceManager() *deviceManager {
	once.Do(func() {
		instance = &deviceManager{
			devices: sync.Map{},
		}
	})
	return instance
}

func (dm *deviceManager) AddDevice(id string, info *DeviceInfo) {
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
