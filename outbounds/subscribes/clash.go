package subscribes

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// 根配置
type ClashConfig struct {
	DNS     *DNSConfig  `yaml:"dns,omitempty"`
	Proxies []ClashNode `yaml:"proxies,omitempty"`
	//Rules   []string    `yaml:"rules,omitempty"`
}

type DNSConfig struct {
	Enable       bool             `yaml:"enable,omitempty"`
	IPv6         bool             `yaml:"ipv6,omitempty"`
	Listen       string           `yaml:"listen,omitempty"`
	FakeIPFilter []string         `yaml:"fake-ip-filter,omitempty"`
	Nameserver   []string         `yaml:"nameserver,omitempty"`
	Timeout      *DurationSeconds `yaml:"timeout,omitempty"`
}

type DurationSeconds struct {
	time.Duration `yaml:"-"`
}

func (d *DurationSeconds) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v float64
	if err := unmarshal(&v); err == nil {
		d.Duration = time.Duration(v) * time.Second
		return nil
	}
	var s string
	if err := unmarshal(&s); err == nil {
		dur, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		d.Duration = dur
		return nil
	}
	return nil
}

// 代理（支持多种类型，部分字段仅针对 ss/vmess/trojan 等）
type ClashNode struct {
	Name       string             `yaml:"name"`
	Type       string             `yaml:"type,omitempty"`
	Server     string             `yaml:"server,omitempty"`
	Port       int                `yaml:"port,omitempty"`
	Cipher     *string            `yaml:"cipher,omitempty"`
	Password   string             `yaml:"password,omitempty"`
	Plugin     *string            `yaml:"plugin,omitempty"`
	PluginOpts *map[string]string `yaml:"plugin-opts,omitempty"`
	UDP        *bool              `yaml:"udp,omitempty"`
	// vmess/trojan 特有字段（可选）
	UUID       *string            `yaml:"uuid,omitempty"`
	AlterID    *int               `yaml:"alterId,omitempty"`
	TLS        *bool              `yaml:"tls,omitempty"`
	SNI        *string            `yaml:"sni,omitempty"`
	Network    *string            `yaml:"network,omitempty"`
	WSPath     *string            `yaml:"ws-path,omitempty"`
	WSSHeaders *map[string]string `yaml:"wss-headers,omitempty"`
}

func (n ClashNode) GetName() string {
	return n.Name
}
func (n ClashNode) GetScheme() string {
	return n.Type
}
func (n ClashNode) GetServer() string {
	return n.Server
}
func (n ClashNode) GetPort() int {
	return n.Port
}
func (n ClashNode) GetCipher() *string {
	return n.Cipher
}
func (n ClashNode) GetPassword() string {
	return n.Password
}
func (n ClashNode) GetPlugin() *string {
	if n.Plugin == nil {
		return nil
	}
	p := "obfs-local"
	if *n.Plugin == "obfs" {
		p = "obfs-local"
	}
	return &p
}
func (n ClashNode) GetPluginOpts() *string {
	if n.PluginOpts == nil {
		return nil
	}
	p := make([]string, 0, len(*n.PluginOpts))
	for key, value := range *n.PluginOpts {
		if key == "mode" {
			key = "bofs"
		}
		if key == "host" {
			key = "obfs-host"
		}
		p = append(p, key+"="+value)
	}
	s := strings.Join(p, ";")
	return &s
}
func (n ClashNode) GetUDP() *bool {
	return n.UDP
}
func (n ClashNode) GetUUID() *string {
	return n.UUID
}
func (n ClashNode) GetAlterID() *int {
	return n.AlterID
}
func (n ClashNode) GetTLS() *bool {
	return n.TLS
}
func (n ClashNode) GetSNI() *string {
	return n.SNI
}
func (n ClashNode) GetNetwork() *string {
	net := "tcp"
	//if n.UDP == nil || *n.UDP == true {
	//	net = "udp"
	//}
	return &net
}
func (n ClashNode) GetWSPath() *string {
	return n.WSPath
}
func (n ClashNode) GetWSSHeaders() *map[string]string {
	return n.WSSHeaders
}

func (n ClashNode) GetAllowInsecure() *bool {
	b := true
	return &b
}

func ReturnNodesFromYaml(tag string, body string) ([]NodeInfo, error) {
	var conf ClashConfig
	if err := yaml.Unmarshal([]byte(body), &conf); err != nil {
		return nil, err
	}
	var _conf []NodeInfo
	for _, node := range conf.Proxies {
		node.Name = fmt.Sprintf("%s#%s", node.Name, tag)
		_conf = append(_conf, node)
	}

	return _conf, nil
}
