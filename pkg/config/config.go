package config

type MainConfig struct {
	Serial    string `ymal:"serial"`
	Realm     string `ymal:"realm"`
	SipHost   string `ymal:"sip-host"`
	SipPort   int    `ymal:"sip-port"`
	MediaAddr string `ymal:"media-addr"`
}
