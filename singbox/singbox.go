package singbox

import (
	"go-singbox/config"
	"go-singbox/dns"
	"go-singbox/logger"
	"go-singbox/outbounds"
	"go-singbox/route"
)

type InboundInterface interface {
	Inbound()
}
type InboundBase struct {
	Type                     string `json:"type"`
	Tag                      string `json:"tag"`
	Sniff                    *bool  `json:"sniff,omitempty"`
	SniffOverrideDestination *bool  `json:"sniff_override_destination,omitempty"`
}

type HttpInbound struct {
	InboundBase

	Listen     string `json:"listen"`
	ListenPort int    `json:"listen_port"`
}

func (in HttpInbound) Inbound() {}

type TunInbound struct {
	InboundBase
	InterfaceName *string  `json:"interface_name,omitempty"`
	Address       []string `json:"address"`
	AutoRoute     bool     `json:"auto_route"`
	AutoRedirect  bool     `json:"auto_redirect"` // 使用 NFTables
	StrictRoute   bool     `json:"strict_route"`
}

func (in TunInbound) Inbound() {}

type TProxyInbound struct {
	InboundBase

	Listen      string `json:"listen"`
	ListenPort  int    `json:"listen_port"`
	TcpFastOpen bool   `json:"tcp_fast_open"`
	UdpFragment bool   `json:"udp_fragment"`
}

func (in TProxyInbound) Inbound() {}

type DirectInbound struct {
	InboundBase
	Network      string `json:"network"`
	OverridePort uint16 `json:"override_port"`
}

func (in DirectInbound) Inbound() {}

var True = true
var InboundsConf []InboundInterface

func init() {
	if config.Conf.Inbounds.Tun {
		InboundsConf = append(InboundsConf, TunInbound{
			InboundBase: InboundBase{
				Type:                     "tun",
				Tag:                      "tun-in",
				Sniff:                    &True,
				SniffOverrideDestination: &True,
			},
			InterfaceName: nil,
			Address: []string{
				"172.18.0.1/30",
			},
			AutoRoute:    true,
			StrictRoute:  true,
			AutoRedirect: true,
		})
	}
	if config.Conf.Inbounds.Mixed != 0 {
		InboundsConf = append(InboundsConf, HttpInbound{
			InboundBase: InboundBase{
				Type:                     "mixed",
				Tag:                      "mixed-in",
				Sniff:                    &True,
				SniffOverrideDestination: &True,
			},
			Listen:     "::",
			ListenPort: config.Conf.Inbounds.Mixed,
		})
	}

	if config.Conf.Inbounds.TProxy != 0 {
		InboundsConf = append(InboundsConf, TProxyInbound{
			InboundBase: InboundBase{
				Type:                     "tproxy",
				Tag:                      "tproxy-in",
				Sniff:                    &True,
				SniffOverrideDestination: &True,
			},
			Listen:      "::",
			ListenPort:  config.Conf.Inbounds.TProxy,
			TcpFastOpen: true,
			UdpFragment: true,
		})
	}

	InboundsConf = append(InboundsConf, DirectInbound{
		InboundBase: InboundBase{
			Type: "direct",
			Tag:  "direct-in",
		},
		Network:      "udp",
		OverridePort: 53,
	})

}

type ClashApi struct {
	ExternalController               string   `json:"external_controller"`
	AccessControlAllowOrigin         []string `json:"access_control_allow_origin"`
	AccessControlAllowPrivateNetwork bool     `json:"access_control_allow_private_network"`
}
type CacheFile struct {
	Enabled   bool `json:"enabled"`
	StoreRdrc bool `json:"store_rdrc"`
}
type Experimental struct {
	ClashApi  ClashApi  `json:"clash_api"`
	CacheFile CacheFile `json:"cache_file"`
}

var ExperimentalConf = Experimental{
	ClashApi: ClashApi{
		ExternalController: "0.0.0.0:9090",
		AccessControlAllowOrigin: []string{
			"http://127.0.0.1",
			"*",
		},
		AccessControlAllowPrivateNetwork: true,
	},
	CacheFile: CacheFile{
		Enabled:   true,
		StoreRdrc: true,
	},
}

type Singbox struct {
	Log          *logger.Log         `json:"log"`
	Experimental *Experimental       `json:"experimental"`
	Inbounds     *[]InboundInterface `json:"inbounds"`
	DNS          *dns.DNS            `json:"dns"`
	Outbounds    *[]outbounds.Proxy  `json:"outbounds"`
	Route        *route.Route        `json:"route"`
}
