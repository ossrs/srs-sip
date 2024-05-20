package service

import (
	"sync"

	"github.com/ossrs/go-oryx-lib/logger"
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
	Status       ChannelStatus `json:"status"`
	Longitude    float64       `json:"longitude"`
	Latitude     float64       `json:"latitude"`
	Info         struct {
		PTZType       int    `json:"ptz_type"`
		Resolution    string `json:"resolution"`
		DownloadSpeed int    `json:"download_speed"`
	} `json:"info"`
}

type ChannelStatus string

type DeviceInfo struct {
	DeviceID    string        `json:"device_id"`
	SourceAddr  string        `json:"source_addr"`
	NetworkType string        `json:"network_type"`
	ChannelList []ChannelInfo `json:"-"`
}

var Devices sync.Map

func AddDevice(id string, info DeviceInfo) {
	Devices.Store(id, info)
}

func RemoveDevice(id string) {
	Devices.Delete(id)
}

func GetDevice(id string) (DeviceInfo, bool) {
	v, ok := Devices.Load(id)
	if !ok {
		return DeviceInfo{}, false
	}
	return v.(DeviceInfo), true
}

func UpdateChannels(deviceID string, list ...ChannelInfo) {
	for _, channel := range list {
		if info, ok := GetDevice(deviceID); ok {
			info.ChannelList = append(info.ChannelList, channel)
			Devices.Store(deviceID, info)
			logger.Tf("Update channel %s", channel.DeviceID)
		}
	}
}

func GetDeviceByChannel(channelID string) (DeviceInfo, bool) {
	var result DeviceInfo
	Devices.Range(func(key, value interface{}) bool {
		info := value.(DeviceInfo)
		for _, channel := range info.ChannelList {
			if channel.DeviceID == channelID {
				result = info
				return false
			}
		}
		return true
	})
	return result, result.DeviceID != ""
}
