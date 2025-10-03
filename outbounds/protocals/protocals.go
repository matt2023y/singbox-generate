package protocals

import "go-singbox/config"

type TLS struct {
	Enabled    bool          `json:"enabled,omitempty"`
	DisableSni bool          `json:"disable_sni,omitempty"`
	ServerName *string       `json:"server_name,omitempty"`
	Insecure   *bool         `json:"insecure,omitempty"`
	Alpn       []interface{} `json:"alpn,omitempty"`
}
type Trojan struct {
	Tag  string `yaml:"tag" json:"tag"`
	Type string `yaml:"type" json:"type"`

	Server     string `yaml:"server" json:"server"`
	ServerPort int    `yaml:"server_port" json:"server_port"`
	Password   string `yaml:"password" json:"password"`
	TLS        TLS    `yaml:"tls" json:"tls"`
}

func (t Trojan) GetTag() string {
	return t.Tag
}
func (t Trojan) GetServer() string {
	return t.Server
}
func (t Trojan) GetSni() string {
	return *t.TLS.ServerName
}

type Vless struct {
	Tag        string `yaml:"tag" json:"tag"`
	Type       string `yaml:"type" json:"type"`
	Server     string `yaml:"server" json:"server"`
	ServerPort int    `yaml:"server_port" json:"server_port"`
	Uuid       string `yaml:"uuid" json:"uuid"`
	Encryption string `yaml:"encryption,omitempty" json:"encryption,omitempty"`
	Flow       string `yaml:"flow,omitempty" json:"flow,omitempty"`
}

func (t Vless) GetTag() string {
	return t.Tag
}

func (t Vless) GetServer() string {
	return t.Server
}
func (t Vless) GetSni() string {
	return ""
}

type Vmess struct {
	Tag  string `yaml:"tag" json:"tag"`
	Type string `yaml:"type" json:"type"`

	Server     string `yaml:"server" json:"server"`
	ServerPort int    `yaml:"server_port" json:"server_port"`

	Uuid     string `yaml:"uuid" json:"uuid"`
	AlterId  int    `yaml:"alterId,omitempty" json:"alterId,omitempty"`
	Security string `yaml:"security,omitempty" json:"security,omitempty"`
}

func (t Vmess) GetTag() string {
	return t.Tag
}
func (t Vmess) GetServer() string {
	return t.Server
}
func (t Vmess) GetSni() string {
	return ""
}

type Hysteria2 struct {
	Tag  string `yaml:"tag" json:"tag"`
	Type string `yaml:"type" json:"type"`

	Server     string `yaml:"server" json:"server"`
	ServerPort int    `yaml:"server_port" json:"server_port"`

	Password string `yaml:"password" json:"password"`
	Obfs     string `yaml:"obfs,omitempty" json:"obfs,omitempty"`
	Tls      bool   `yaml:"tls,omitempty" json:"tls,omitempty"`
}

func (h Hysteria2) GetTag() string {
	return h.Tag
}
func (h Hysteria2) GetServer() string {
	return h.Server
}
func (h Hysteria2) GetSni() string {
	return ""
}

type SS struct {
	Type string `yaml:"type" json:"type"`
	Tag  string `yaml:"tag" json:"tag"`

	Server     string `yaml:"server" json:"server"`
	ServerPort int    `yaml:"server_port" json:"server_port"`

	Method   *string `yaml:"method" json:"method"`
	Password string  `yaml:"password" json:"password"`

	Plugin     *string `yaml:"plugin,omitempty" json:"plugin,omitempty"`
	PluginOpts *string `yaml:"plugin_opts,omitempty" json:"plugin_opts,omitempty"`
	Network    *string `yaml:"network,omitempty" json:"network,omitempty"`
}

func (s SS) GetTag() string {
	return s.Tag
}
func (s SS) GetServer() string {
	return s.Server
}
func (s SS) GetSni() string {
	return ""
}

type Base struct {
	Type      string   `yaml:"type" json:"type"`
	Tag       string   `yaml:"tag" json:"tag"`
	Outbounds []string `yaml:"outbounds,omitempty" json:"outbounds,omitempty"`
	Default   string   `yaml:"default,omitempty" json:"default,omitempty"`
}

func (b Base) GetTag() string {
	return b.Tag
}
func (b Base) GetServer() string {
	return ""
}
func (b Base) GetSni() string {
	return ""
}

type UrlTest struct {
	Base
	config.UrlTest
	Outbounds []string `yaml:"outbounds,omitempty" json:"outbounds,omitempty"`
}
