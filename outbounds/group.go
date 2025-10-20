package outbounds

import (
	"fmt"
	"go-singbox/config"
	"go-singbox/outbounds/protocals"
	"strings"
)

var groups = config.Conf.PeerGroups
var _groups map[string]string

func init() {
	_groups = make(map[string]string)
	for k, v := range groups {
		for _, s := range v {
			_groups[s] = k
		}
	}
}

var Apps = config.Conf.AppGroups

func ProxyAddOutbounds(proxy []Proxy) []Proxy {
	tagCounts := make(map[string]int, len(proxy))

	allTags := make([]string, 0, len(proxy))
	groupOutbounds := make(map[string][]string)
	// 将 proxy 按照关键字进行分组
	for _, p := range proxy {
		tag := p.GetTag()
		tagCounts[tag] = 1
		allTags = append(allTags, tag)
		ct := 0
		for m, v := range _groups {
			if strings.Contains(tag, m) {
				tagCounts[v] = 1
				groupOutbounds[v] = append(groupOutbounds[v], tag)
				ct++
			}
		}
		if ct == 0 {
			groupOutbounds["Other"] = append(groupOutbounds["Other"], tag)
			tagCounts["Other"] = 1
		}
	}

	// 添加 HK JP 分组
	for zone, tags := range groupOutbounds {
		proxy = append(proxy,
			protocals.UrlTest{
				Base: protocals.Base{
					Type: "urltest",
					//Tag:  fmt.Sprintf("auto-%s", zone),
					Tag: zone,
				},
				UrlTest:   config.UrlTestConf,
				Outbounds: tags,
			},
		)
	}

	// auto group
	for app, tags := range Apps {
		_tags := removeEle(tags, "Select")
		for _, tag := range tags {
			if tagCounts[tag] == 0 {
				_tags = removeEle(_tags, tag)
			}
		}
		if len(_tags) == 0 {
			_tags = append(_tags, "Select")
		}
		proxy = append(proxy, protocals.UrlTest{
			Base: protocals.Base{
				Type: "urltest",
				Tag:  fmt.Sprintf("Auto%s", app),
			},
			UrlTest:   config.UrlTestConf,
			Outbounds: _tags,

			//Type:      "selector",
			//Tag:       fmt.Sprintf("Auto%s", app),
			//Outbounds: tags,
			//Default:   tags[0],
		})
	}
	// group
	for app, tags := range Apps {
		for _, tag := range tags { // 清理没有的group
			if tagCounts[tag] == 0 {
				tags = removeEle(tags, tag)
			}
		}

		var _tags []string
		_tags = append(_tags, fmt.Sprintf("Auto%s", app))
		_tags = append(_tags, tags...)
		proxy = append(proxy, protocals.Base{
			Type:      "selector",
			Tag:       app,
			Outbounds: _tags,
			Default:   _tags[0],
		})
	}

	proxy = append(proxy,
		protocals.UrlTest{
			Base: protocals.Base{
				Type: "urltest",
				Tag:  "AutoSelect",
			},
			UrlTest:   config.UrlTestConf,
			Outbounds: removePartialEle(allTags, "个人专用"),
		},
		protocals.Base{
			Type:      "selector",
			Tag:       "Select",
			Outbounds: allTags,
			Default:   allTags[0],
		},
		protocals.Base{
			Type: "direct",
			Tag:  "direct",
		},
		protocals.Base{
			Type: "block",
			Tag:  "block",
		},
	)
	return proxy
}

func removeEle(s []string, e string) []string {
	var ss []string
	for _, v := range s {
		if v == e {
			continue
		}
		ss = append(ss, v)
	}
	return ss
}

func removePartialEle(s []string, e string) []string {
	var ss []string
	for _, v := range s {
		if strings.Contains(v, e) {
			continue
		}
		ss = append(ss, v)
	}
	return ss
}
