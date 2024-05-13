package service

import "sync"

type ChannelInfo struct {
	DeviceID     string // 通道ID
	ParentID     string
	Name         string
	Manufacturer string
	Model        string
	Owner        string
	CivilCode    string
	Address      string
	Port         int
	Parental     int
	SafetyWay    int
	RegisterWay  int
	Secrecy      int
	Status       ChannelStatus
}

type ChannelStatus string

type DeviceInfo struct {
	DeviceID    string
	SourceAddr  string
	ChannelList []ChannelInfo
}

var Devices sync.Map

func (s *GB28181Server) AddDevice(id string, sourceAddr string) {
	Devices.Store(id, DeviceInfo{
		DeviceID:   id,
		SourceAddr: sourceAddr,
	})
}

func (s *GB28181Server) RemoveDevice(id string) {
	Devices.Delete(id)
}

func (s *GB28181Server) GetDevice(id string) (DeviceInfo, bool) {
	v, ok := Devices.Load(id)
	if !ok {
		return DeviceInfo{}, false
	}
	return v.(DeviceInfo), true
}

func (s *GB28181Server) UpdateChannels(list ...ChannelInfo) {
	for _, channel := range list {
		if info, ok := s.GetDevice(channel.ParentID); ok {
			info.ChannelList = append(info.ChannelList, channel)
			Devices.Store(channel.DeviceID, info)
		}
	}
}

func (s *GB28181Server) GetDeviceByChannel(channelID string) (DeviceInfo, bool) {
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
