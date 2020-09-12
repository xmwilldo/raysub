package config

type V2rayCoreConfig struct {
	Log       Log         `json:"log"`
	Inbounds  []Inbounds  `json:"inbounds"`
	Outbounds []Outbounds `json:"outbounds"`
	DNS       DNS         `json:"dns"`
	Routing   Routing     `json:"routing"`
	Transport Transport   `json:"transport"`
}

type Log struct {
	Error    string `json:"error"`
	Loglevel string `json:"loglevel"`
	Access   string `json:"access"`
}

type Settings struct {
	UDP     *bool  `json:"udp,omitempty"`
	Auth    string `json:"auth,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

type Inbounds struct {
	Listen   string   `json:"listen"`
	Protocol string   `json:"protocol"`
	Settings Settings `json:"settings"`
	Port     string   `json:"port"`
}

type Mux struct {
	Enabled     bool `json:"enabled"`
	Concurrency int  `json:"concurrency,omitempty"`
}

type Header struct {
	Type string `json:"type,omitempty"`
}

type TCPSettings struct {
	Header Header `json:"header,omitempty"`
}

type TLSSettings struct {
	AllowInsecure bool `json:"allowInsecure,omitempty"`
}

type StreamSettings struct {
	TCPSettings TCPSettings `json:"tcpSettings,omitempty"`
	TLSSettings TLSSettings `json:"tlsSettings,omitempty"`
	Security    string      `json:"security,omitempty"`
	Network     string      `json:"network,omitempty"`
}

type Users struct {
	ID       string `json:"id"`
	AlterID  int    `json:"alterId"`
	Level    int    `json:"level"`
	Security string `json:"security"`
}

type Vnext struct {
	Address string  `json:"address,omitempty"`
	Users   []Users `json:"users,omitempty"`
	Port    int     `json:"port,omitempty"`
}

type Response struct {
	Type string `json:"type,omitempty"`
}

type Settings0 struct {
	Vnext          []Vnext   `json:"vnext,omitempty"`
	DomainStrategy *string   `json:"domainStrategy,omitempty"`
	Redirect       *string   `json:"redirect,omitempty"`
	UserLevel      *int      `json:"userLevel,omitempty"`
	Response       *Response `json:"response,omitempty"`
}

type Outbounds struct {
	Mux            *Mux            `json:"mux,omitempty"`
	Protocol       string          `json:"protocol"`
	StreamSettings *StreamSettings `json:"streamSettings,omitempty"`
	Tag            string          `json:"tag"`
	Settings       *Settings0      `json:"settings"`
}

type DNS struct {
}

type Settings1 struct {
	DomainStrategy string        `json:"domainStrategy"`
	Rules          []interface{} `json:"rules"`
}

type Routing struct {
	Settings Settings1 `json:"settings"`
}

type Transport struct {
}

type ProxyConfig struct {
	V    string `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	ID   string `json:"id"`
	Aid  string `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	Path string `json:"path"`
	TLS  string `json:"tls"`
}
