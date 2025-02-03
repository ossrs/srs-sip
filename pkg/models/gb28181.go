package models

import "encoding/xml"

type Record struct {
	DeviceID  string `xml:"DeviceID" json:"device_id"`
	Name      string `xml:"Name" json:"name"`
	FilePath  string `xml:"FilePath" json:"file_path"`
	Address   string `xml:"Address" json:"address"`
	StartTime string `xml:"StartTime" json:"start_time"`
	EndTime   string `xml:"EndTime" json:"end_time"`
	Secrecy   int    `xml:"Secrecy" json:"secrecy"`
	Type      string `xml:"Type" json:"type"`
}

// Example XML structure for channel info:
//
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
		DownloadSpeed string `json:"download_speed"` // Speed levels: 1/2/4/8
	} `json:"info"`

	// Custom fields
	Ssrc string `json:"ssrc"`
}

type ChannelStatus string

type XmlMessageInfo struct {
	XMLName      xml.Name
	CmdType      string
	SN           int
	DeviceID     string
	DeviceName   string
	Manufacturer string
	Model        string
	Channel      string
	DeviceList   []ChannelInfo `xml:"DeviceList>Item"`
	RecordList   []*Record     `xml:"RecordList>Item"`
	SumNum       int
}
