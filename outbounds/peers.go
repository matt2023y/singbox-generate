package outbounds

import (
	"fmt"
	"go-singbox/config"
	"go-singbox/outbounds/protocals"
	"go-singbox/outbounds/subscribes"
	"log"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Proxy interface {
	GetTag() string
	GetServer() string
	GetSni() string
}

func ParseProxyFromURL() ([]Proxy, error) {
	links := config.Conf.ClashURL
	var nodes []subscribes.NodeInfo

	for tag, link := range links {
		if strings.HasPrefix(link, "https://") {
			client := resty.New()
			resp, err := client.R().SetHeader("User-Agent", UA).Get(link)
			if err != nil {
				return nil, err
			}
			body := resp.Body()
			if f, e := os.OpenFile(fmt.Sprintf("config/proxy-%s.yaml", tag), os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, os.ModePerm); e == nil {
				defer f.Close()
				_, _ = f.Write(body)
			}

			_nodes, err := subscribes.ReturnNodesFromYaml(tag, string(body))
			if err != nil {
				fmt.Println("解析节点出错", err)
				return nil, err
			}
			nodes = append(nodes, _nodes...)
			//nodes, err = types.ParseByBase64(string(body)) // Base64 -- trojan 格式
			//if err != nil {
			//}

		} else {
			//nodes = strings.Split(link, "\n") // 自建
		}
	}

	for tag, c := range config.Conf.SelfClash {
		c.Name = fmt.Sprintf("%s#%s", c.Name, tag)
		nodes = append(nodes, c)
	}

	proxies := make([]Proxy, 0)
	for _, u := range nodes {
		var p Proxy
		switch u.GetScheme() {
		case "trojan":
			tls := protocals.TLS{
				Enabled:    true,
				ServerName: u.GetSNI(),
				Insecure:   u.GetAllowInsecure(), // toBool(u.Query().Get("allowInsecure")),
			}
			p = protocals.Trojan{
				Type:       "trojan",
				Server:     u.GetServer(),
				ServerPort: u.GetPort(),
				Password:   u.GetPassword(),
				Tag:        u.GetName(),
				TLS:        tls,
			}
		case "vless", "vmess":
			p = protocals.Vless{
				Type:       u.GetScheme(),
				Server:     u.GetServer(),
				ServerPort: u.GetPort(),
				Uuid:       u.GetPassword(),
				Tag:        u.GetName(),
			}
		case "hysteria2":
			p = protocals.Hysteria2{
				Type:       "hysteria2",
				Server:     u.GetServer(),
				ServerPort: u.GetPort(),
				Password:   u.GetPassword(),
				Tag:        u.GetName(),
			}
		case "ss":
			{
				p = protocals.SS{
					Type: "shadowsocks",
					Tag:  u.GetName(),

					Server:     u.GetServer(),
					ServerPort: u.GetPort(),

					Method:   u.GetCipher(),
					Password: u.GetPassword(),

					Plugin:     u.GetPlugin(),
					PluginOpts: u.GetPluginOpts(),

					Network: u.GetNetwork(),
				}
			}
		default:
			log.Println("不支持的协议: " + u.GetScheme())
			continue
		}
		proxies = append(proxies, p)
	}
	return proxies, nil
}
