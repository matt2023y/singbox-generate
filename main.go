package main

import (
	"encoding/json"
	"fmt"
	"go-singbox/config"
	"go-singbox/dns"
	"go-singbox/logger"
	"go-singbox/outbounds"
	"go-singbox/route"
	"go-singbox/singbox"
	"maps"
	"os"
	"slices"

	"log"
)

func main() {
	proxies, err := outbounds.ParseProxyFromURL()
	if err != nil {
		logger.Error(fmt.Errorf("解析节点失败: %w", err))
		return
	}
	log.Printf("解析到 %d 个节点", len(proxies))
	if len(proxies) == 0 {
		panic("没有收到节点配置")
	}
	proxies = outbounds.ProxyAddOutbounds(proxies)

	outboundRuleSet, outboundRules, dnsRules, err := route.GenRouteAndDNSRules()
	if err != nil {
		logger.Error(err)
	}

	// 生成本地节点的 dns
	peerDomains := make(map[string]struct{})
	for _, p := range proxies {
		s := p.GetServer()
		if s == "" {
			continue
		}
		peerDomains[s] = struct{}{}

		s = p.GetSni()
		if s == "" {
			continue
		}
		peerDomains[s] = struct{}{}
	}

	peersDNS := dns.GetPeerDNS()
	peersDNS.Domain = slices.Sorted(maps.Keys(peerDomains))
	dnsRules = append(dnsRules, &peersDNS)

	dnsConf := dns.GetDNSConf()
	dnsConf.Rules = dnsRules

	peerRoutRule := route.GetPeerRouteRule(slices.Sorted(maps.Keys(peerDomains)))

	singConf := singbox.Singbox{
		Experimental: &singbox.ExperimentalConf,
		Log:          &logger.LogConf,
		DNS:          &dnsConf,
		Inbounds:     &singbox.InboundsConf,
		Outbounds:    &proxies,
		Route: &route.Route{
			Rules:                 append(outboundRules, &peerRoutRule),
			RuleSet:               outboundRuleSet,
			DefaultDomainResolver: config.Conf.DefaultDomainResolver,
			Final:                 "direct", // route
			AutoDetectInterface:   true,
		},
	}
	b, e := json.MarshalIndent(singConf, "", "  ")
	if e != nil {
		logger.Error(e)
	}
	f, _ := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer f.Close()
	f.Write(b)
}
