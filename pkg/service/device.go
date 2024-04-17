package service

import (
	"sync"

	"github.com/ossrs/srs-sip/pkg/utils"
)

// <Item>
// 	<DeviceID>34020000001320000002</DeviceID>
// 	<Name>209</Name>
// 	<Manufacturer>UNIVIEW</Manufacturer>
// 	<Model>HIC6622-IR@X33-VF</Model>
// 	<Owner>IPC-B2202.7.11.230222</Owner>
// 	<CivilCode>CivilCode</CivilCode>
// 	<Address>Address</Address>
// 	<Parental>1</Parental>
// 	<ParentID>75015310072008100002</ParentID>
// 	<SafetyWay>0</SafetyWay>
// 	<RegisterWay>1</RegisterWay>
// 	<Secrecy>0</Secrecy>
// 	<Status>ON</Status>
// 	<Longitude>0.0000000</Longitude>
// 	<Latitude>0.0000000</Latitude>
// 	<Info>
// 		<PTZType>1</PTZType>
// 		<Resolution>6/4/2</Resolution>
// 		<DownloadSpeed>0</DownloadSpeed>
// 	</Info>
// </Item>

type ChannelInfo struct {
	DeviceID     string        `json:"device_id"`
	ParentID     string        `json:"parent_id"`
	Name         string        `json:"name"`
	Manufacturer string        `json:"manufacturer"`
	Model        string        `json:"model"`
	Owner        string        `json:"owner"`
	CivilCode    string        `json:"civil_code"`
	Address      string        `json:"address"`
	Port         int           `json:"port"`
	Parental     int           `json:"parental"`
	SafetyWay    int           `json:"safety_way"`
	RegisterWay  int           `json:"register_way"`
	Secrecy      int           `json:"secrecy"`
	IPAddress    string        `json:"ip_address"`
	Status       ChannelStatus `json:"status"`
	Longitude    float64       `json:"longitude"`
	Latitude     float64       `json:"latitude"`
	Info         struct {
		PTZType       int    `json:"ptz_type"`
		Resolution    string `json:"resolution"`
		DownloadSpeed string `json:"download_speed"` // 1/2/4/8
	} `json:"info"`

	// custom fields
	Ssrc string `json:"ssrc"`
}

type ChannelStatus string

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

func (dm *deviceManager) UpdateChannels(deviceID string, list ...ChannelInfo) {
	device, ok := dm.GetDevice(deviceID)
	if !ok {
		return
	}

	for _, channel := range list {
		device.ChannelMap.Store(channel.DeviceID, channel)
	}
	dm.devices.Store(deviceID, device)
}

func (dm *deviceManager) ApiGetChannelByDeviceId(deviceID string) []ChannelInfo {
	device, ok := dm.GetDevice(deviceID)
	if !ok {
		return nil
	}

	channels := make([]ChannelInfo, 0)
	device.ChannelMap.Range(func(key, value interface{}) bool {
		channels = append(channels, value.(ChannelInfo))
		return true
	})
	return channels
}

func (dm *deviceManager) GetAllVideoChannels() []ChannelInfo {
	channels := make([]ChannelInfo, 0)
	dm.devices.Range(func(key, value interface{}) bool {
		device := value.(*DeviceInfo)
		device.ChannelMap.Range(func(key, value interface{}) bool {
			if utils.IsVideoChannel(value.(ChannelInfo).DeviceID) {
				channels = append(channels, value.(ChannelInfo))
				return true
			}
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
