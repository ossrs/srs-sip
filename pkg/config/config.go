package config

type MainConfig struct {
	Serial       string
	Realm        string
	SipNetType   string
	SipHost      string
	SipPort      uint
	MediaHost    string
	MediaApiPort uint
}
