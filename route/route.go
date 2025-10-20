package route

import (
	"fmt"
	"go-singbox/config"
	"strings"
)

type BaseRule struct {
	Domain        []string `json:"domain,omitempty" `
	DomainSuffix  []string `json:"domain_suffix,omitempty"`
	DomainKeyword []string `json:"domain_keyword,omitempty"`
	DomainRegex   []string `json:"domain_regex,omitempty"`
	IpCidr        []string `json:"ip_cidr,omitempty"`
	Protocol      []string `json:"protocol,omitempty"`
	RuleSet       []string `json:"rule_set,omitempty"`
	Invert        bool     `json:"invert"`
	ClashMode     *string  `json:"clash_mode,omitempty"`
	Port          []uint16 `json:"port,omitempty"`
	Action        *string  `json:"action,omitempty"`
}

type RouteRule struct {
	BaseRule
	Outbound string `json:"outbound,omitempty"`
}
type DNSRule struct {
	BaseRule
	RuleSetIpCidrMatchSource bool   `json:"rule_set_ip_cidr_match_source"`
	RuleSetIpCidrAcceptEmpty bool   `json:"rule_set_ip_cidr_accept_empty"`
	Action                   string `json:"action,omitempty"`
	Server                   string `json:"server,omitempty"`
	DisableCache             bool   `json:"disable_cache,omitempty"`
}
type Route struct {
	Rules   []*RouteRule `json:"rules"`
	RuleSet RouteRuleSet `json:"rule_set,omitempty"`

	DefaultDomainResolver string `json:"default_domain_resolver"`
	Final                 string `json:"final"`
	AutoDetectInterface   bool   `json:"auto_detect_interface"`
}

func GenRouteAndDNSRules() (RouteRuleSet, []*RouteRule, []*DNSRule, error) {
	conf := config.Conf

	outboundRules := make([]*RouteRule, 0, len(conf.Rules)+3)
	dnsRules := make([]*DNSRule, 0, len(conf.Rules))
	var sniffRouteRule RouteRule
	var hijackRouteRule RouteRule
	var clashGlobalRouteRule RouteRule
	var clashDirectRouteRule RouteRule
	var clashGlobalDNSRule, clashDirectDNSRule DNSRule
	var quicRouteRule RouteRule

	// route 独有的rule
	sniffAction := "sniff"
	sniffRouteRule.Action = &sniffAction
	hijackAction := "hijack-dns"
	hijackRouteRule.Action = &hijackAction
	hijackRouteRule.Protocol = []string{"dns"}
	hijackRouteRule.Port = []uint16{53}

	clashGlobalAction := "global"
	clashGlobalRouteRule.ClashMode = &clashGlobalAction
	clashGlobalRouteRule.Outbound = "global"

	clashDirectAction := "direct"
	clashDirectRouteRule.ClashMode = &clashDirectAction
	clashDirectRouteRule.Outbound = "direct"

	clashGlobalDNSRule.ClashMode = &clashGlobalAction
	clashGlobalDNSRule.Server = config.Conf.OuterDNS

	clashDirectDNSRule.ClashMode = &clashDirectAction
	clashDirectDNSRule.Server = config.Conf.InnerDNS

	quicRouteRule.Protocol = []string{"quic"}
	quicRouteRule.Outbound = "block"

	outboundRules = append(outboundRules, &sniffRouteRule, &hijackRouteRule, &clashGlobalRouteRule, &clashDirectRouteRule, &quicRouteRule)
	dnsRules = append(dnsRules, &clashGlobalDNSRule, &clashDirectDNSRule)

	for _, rule := range conf.Rules {
		lines, err := config.ParseListFile(rule.Name)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to parse rules file [rules.yml]: %w", err)
		}
		outboundRule, dnsRule := groupRule(lines, rule.Group, rule.DNS)

		outboundRules = append(outboundRules, outboundRule)
		dnsRules = append(dnsRules, dnsRule)
	}

	rrs, rr, dr := groupRuleSet(conf.RuleSets)

	outboundRules = append(outboundRules, rr...)
	dnsRules = append(dnsRules, dr...)

	return rrs, outboundRules, dnsRules, nil
}

func groupRule(s []string, outbound string, dnsTag string) (*RouteRule, *DNSRule) {
	var outboundRule RouteRule
	var dnsRule DNSRule

	var rule BaseRule
	for _, l := range s {
		if strings.HasPrefix(l, "#") {
			continue
		}
		parts := strings.Split(strings.TrimSpace(l), ",")
		if len(parts) <= 1 {
			continue
		}
		typ, record := parts[0], parts[1]
		switch strings.ToUpper(typ) {
		case "USER-AGENT", "PROCESS-NAME":
		case "DOMAIN":
			rule.Domain = append(rule.Domain, record)
		case "DOMAIN-SUFFIX":
			rule.DomainSuffix = append(rule.DomainSuffix, record)
		case "DOMAIN-KEYWORD":
			rule.DomainKeyword = append(rule.DomainKeyword, record)
		case "URL-REGEX":
			rule.DomainRegex = append(rule.DomainRegex, record)
		case "IP-CIDR", "IP-CIDR6":
			rule.IpCidr = append(rule.IpCidr, record)
		default:
			fmt.Printf("unknown group rule: %s\n", parts[0])
		}
	}
	outboundRule.BaseRule = rule
	outboundRule.Outbound = outbound

	dnsRule.BaseRule = rule
	if dnsTag == "" {
		dnsRule.Action = "reject"
		dnsRule.Server = dnsTag
		dnsRule.DisableCache = false
	} else {
		dnsRule.Action = "route"
		dnsRule.Server = dnsTag
		dnsRule.DisableCache = false
	}
	return &outboundRule, &dnsRule
}

func groupRuleSet(r []*config.RuleSet) (RouteRuleSet, []*RouteRule, []*DNSRule) {
	var ruleSet RouteRuleSet
	var rr []*RouteRule
	var dr []*DNSRule

	for _, rs := range r {
		base := RouteRuleSetBase{
			Type:   rs.Type,
			Tag:    rs.Name,
			Format: rs.Format,
		}
		if rs.Type == "local" {
			ruleSet = append(ruleSet, RouteRuleSetLocalEle{
				RouteRuleSetBase: base,
				Path:             rs.Path,
			})
		} else {
			ruleSet = append(ruleSet, RouteRuleSetRemoteEle{
				RouteRuleSetBase: base,
				Url:              rs.Url,
				DownloadDetour:   rs.DownloadDetour,
				UpdateInterval:   rs.UpdateInterval,
			})
		}

		baseSet := BaseRule{
			RuleSet: []string{rs.Name},
			Invert:  rs.Invert,
		}
		rr = append(rr, &RouteRule{
			BaseRule: baseSet,
			Outbound: rs.Group,
		})
		dr = append(dr, &DNSRule{
			BaseRule:     baseSet,
			Action:       "route",
			Server:       rs.DNS,
			DisableCache: false,
		})
	}
	return ruleSet, rr, dr
}

func GetPeerRouteRule(domains []string) RouteRule {
	action := "route"

	var peerRule RouteRule
	peerRule.BaseRule = BaseRule{
		Domain: domains,
	}
	peerRule.Action = &action
	peerRule.Outbound = "direct"

	return peerRule
}
