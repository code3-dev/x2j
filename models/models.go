package models

// V2RayConfig represents the complete V2Ray configuration structure
type V2RayConfig struct {
	Log       LogConfig       `json:"log"`
	Inbounds  []InboundConfig `json:"inbounds"`
	Outbounds []OutboundConfig `json:"outbounds"`
	DNS       DNSConfig       `json:"dns"`
	Routing   RoutingConfig   `json:"routing"`
}

// LogConfig represents the logging configuration
type LogConfig struct {
	Access   string `json:"access"`
	Error    string `json:"error"`
	LogLevel string `json:"loglevel"`
	DNSLog   bool   `json:"dnsLog"`
}

// InboundConfig represents an inbound proxy configuration
type InboundConfig struct {
	Tag      string                 `json:"tag"`
	Port     int                    `json:"port"`
	Protocol string                 `json:"protocol"`
	Listen   string                 `json:"listen"`
	Settings map[string]interface{} `json:"settings"`
	Sniffing *SniffingConfig        `json:"sniffing,omitempty"`
}

// SniffingConfig represents the sniffing configuration
type SniffingConfig struct {
	Enabled       bool     `json:"enabled"`
	DestOverride  []string `json:"destOverride,omitempty"`
	MetadataOnly  *bool    `json:"metadataOnly,omitempty"`
	RouteOnly     *bool    `json:"routeOnly,omitempty"`
}

// OutboundConfig represents an outbound proxy configuration
type OutboundConfig struct {
	Tag             string                 `json:"tag"`
	Protocol        string                 `json:"protocol"`
	Settings        map[string]interface{} `json:"settings"`
	StreamSettings  *StreamSettings        `json:"streamSettings,omitempty"`
	ProxySettings   interface{}            `json:"proxySettings,omitempty"`
	SendThrough     interface{}            `json:"sendThrough,omitempty"`
	Mux             *MuxConfig             `json:"mux,omitempty"`
}

// StreamSettings represents stream settings for connections
type StreamSettings struct {
	Network             string                 `json:"network"`
	Security            string                 `json:"security,omitempty"`
	TCPSettings         interface{}            `json:"tcpSettings,omitempty"`
	KCPSettings         interface{}            `json:"kcpSettings,omitempty"`
	WSSettings          interface{}            `json:"wsSettings,omitempty"`
	HTTPSettings        interface{}            `json:"httpSettings,omitempty"`
	TLSSettings         interface{}            `json:"tlsSettings,omitempty"`
	QUICSettings        interface{}            `json:"quicSettings,omitempty"`
	RealitySettings     interface{}            `json:"realitySettings,omitempty"`
	GRPCSettings        interface{}            `json:"grpcSettings,omitempty"`
	XHTTPSettings       interface{}            `json:"xhttpSettings,omitempty"`
	HTTPUpgradeSettings interface{}            `json:"httpupgradeSettings,omitempty"`
	DSSettings          interface{}            `json:"dsSettings,omitempty"`
	Sockopt             interface{}            `json:"sockopt,omitempty"`
}

// MuxConfig represents mux settings
type MuxConfig struct {
	Enabled     bool `json:"enabled"`
	Concurrency int  `json:"concurrency"`
}

// DNSConfig represents DNS configuration
type DNSConfig struct {
	Servers []string `json:"servers"`
}

// RoutingConfig represents routing configuration
type RoutingConfig struct {
	DomainStrategy string        `json:"domainStrategy"`
	Rules          []interface{} `json:"rules"`
	Balancers      []interface{} `json:"balancers"`
}

// V2RayURL is the base interface for all V2Ray URL parsers
type V2RayURL interface {
	GetConfig() *V2RayConfig
}