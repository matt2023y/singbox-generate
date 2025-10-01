package route

type DNSRuleSet []string

type RouteRuleSetBase struct {
	Type   string `json:"type"`
	Tag    string `json:"tag"`
	Format string `json:"format"`
}

type RouteRuleSetEle interface {
	RuleSetEle()
}
type RouteRuleSetLocalEle struct {
	RouteRuleSetBase
	Path string `json:"path"`
}

func (RouteRuleSetLocalEle) RuleSetEle() {}

type RouteRuleSetRemoteEle struct {
	RouteRuleSetBase
	Url            string `json:"url"`
	DownloadDetour string `json:"download_detour"`
	UpdateInterval string `json:"update_interval"`
}

func (RouteRuleSetRemoteEle) RuleSetEle() {}

type RouteRuleSet []RouteRuleSetEle
