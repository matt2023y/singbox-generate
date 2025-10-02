package config

import (
	"errors"
	"fmt"
	"go-singbox/outbounds/subscribes"
	"io"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type RuleSet struct {
	Name           string `yaml:"-"`
	Type           string `yaml:"-"`
	Format         string `yaml:"-"`
	Url            string `yaml:"-"`
	Path           string `yaml:"path"`
	Group          string `yaml:"group"`
	DNS            string `yaml:"dns"`
	DownloadDetour string `yaml:"download_detour"`
	UpdateInterval string `yaml:"update_interval"`
	Invert         bool   `yaml:"invert"`
}

type UrlTest struct {
	Url       *string `yaml:"url,omitempty" json:"url,omitempty"`
	Interval  *string `yaml:"interval,omitempty" json:"interval,omitempty"`
	Tolerance *int    `yaml:"tolerance,omitempty" json:"tolerance,omitempty"`
}
type Config struct {
	ClashURL  map[string]string               `yaml:"clash_url"`
	SelfClash map[string]subscribes.ClashNode `yaml:"self_clash"`

	DefaultDomainResolver string `yaml:"default_domain_resolver"`

	InnerDNS string `yaml:"inner_dns"`
	OuterDNS string `yaml:"outer_dns"`

	Inbounds struct {
		Tun    bool `yaml:"tun"`
		Mixed  int  `yaml:"mixed"`
		TProxy int  `yaml:"tproxy"`
	} `yaml:"inbounds"` // 入站类型

	PeerGroups map[string][]string `yaml:"peer_groups"` // peer 分组
	AppGroups  map[string][]string `yaml:"app_groups"`  // app 分组

	DNSServers []struct {
		Tag        string `yaml:"tag"`
		Type       string `yaml:"type"`
		Server     string `yaml:"server"`
		ServerPort int    `yaml:"port"`
		Final      bool   `yaml:"final"`
	} `yaml:"dns_servers"`

	Rules []struct { // 规则list 列表
		Name  string `yaml:"name"`
		Group string `yaml:"group"`
		DNS   string `yaml:"dns"`
	} `yaml:"rules"`

	RuleSets []*RuleSet `yaml:"rule_sets"`

	UrlTest *UrlTest `yaml:"url_test,omitempty"`
}

const BaseDir = "./config"

var Conf Config
var UrlTestConf UrlTest

func init() {
	file, err := os.ReadFile(path.Join(BaseDir, "config.yaml"))
	if err != nil {
		panic(fmt.Errorf("failed to read config.yaml: %w", err))
	}
	if err := yaml.Unmarshal(file, &Conf); err != nil {
		panic(fmt.Errorf("failed to parse config.yaml: %w", err))
	}

	UrlTestConf = *Conf.UrlTest
	if UrlTestConf.Interval == nil {
		i := "10m"
		UrlTestConf.Interval = &i
	}
	if UrlTestConf.Tolerance == nil {
		t := 100
		UrlTestConf.Tolerance = &t
	}
	if UrlTestConf.Url == nil {
		u := "https://www.gstatic.com/generate_204"
		UrlTestConf.Url = &u
	}

	for _, ruleSet := range Conf.RuleSets {
		rulePath := ruleSet.Path
		if strings.HasSuffix(rulePath, ".srs") {
			ruleSet.Format = "binary"
		} else {
			ruleSet.Format = "source"
		}

		if strings.HasPrefix(rulePath, "http://") || strings.HasPrefix(rulePath, "https://") {
			ruleSet.Type = "remote"
			ruleSet.Url = rulePath
			ruleSet.Path = ""
		} else {
			ruleSet.Type = "local"
		}

		_, fileName := path.Split(rulePath)
		ruleSet.Name = fileName

		ruleSet.Name = strings.TrimSpace(ruleSet.Name)
	}
}

func handleFileError(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("rules file doesn't exist: %w", err)
	}
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("rules file doesn't have permission: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to open rules file: %w", err)
	}
	return nil
}
func closeNoErr(f *os.File) {
	_ = f.Close()
}

func ParseListFile(p string) ([]string, error) {
	f, err := os.Open(fmt.Sprintf(path.Join(BaseDir, "list/%s"), p))
	if err = handleFileError(err); err != nil {
		return nil, err
	}
	defer closeNoErr(f)
	buf, _ := io.ReadAll(f)
	lines := strings.Split(
		strings.TrimSpace(
			string(buf),
		),
		"\n",
	)
	return lines, nil
}
