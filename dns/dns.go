package dns

import (
	"go-singbox/config"
	"go-singbox/route"
)

type DNSServer struct {
	Tag        string `json:"tag"`
	Type       string `json:"type"`
	Server     string `json:"server"`
	ServerPort int    `json:"server_port,omitempty"`
}
type DNS struct {
	Servers       []DNSServer      `json:"servers"`
	Rules         []*route.DNSRule `json:"rules"`
	DisableCache  bool             `json:"disable_cache"`
	Strategy      string           `json:"strategy"`
	CacheCapacity int              `json:"cache_capacity,omitempty"`
	Final         string           `json:"final"`
}

func GetDNSConf() DNS {
	var servers []DNSServer
	for _, serverConfig := range config.Conf.DNSServers {
		servers = append(servers, DNSServer{
			Tag:        serverConfig.Tag,
			Type:       serverConfig.Type,
			Server:     serverConfig.Server,
			ServerPort: serverConfig.ServerPort,
		})
	}
	return DNS{
		DisableCache:  false,
		Strategy:      "ipv4_only",
		CacheCapacity: 5120,
		Final:         config.Conf.DefaultDomainResolver, // "114-dns",
		Servers:       servers,
	}
}

func GetPeerDNS() route.DNSRule {
	return route.DNSRule{
		Action:       "route",
		Server:       "cc-dns",
		DisableCache: false,
	}
}

// TLS 853、HTTPS 443、UDP 53
