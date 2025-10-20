package subscribes

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
)

func ParseByBase64(tag string, body string) ([]NodeInfo, error) {
	var nodeStrs []string

	bodyStr := strings.TrimSpace(string(body))
	missingPadding := len(bodyStr) % 4
	if missingPadding != 0 {
		bodyStr += strings.Repeat("=", 4-missingPadding)
	}
	decoded, err := base64.StdEncoding.DecodeString(bodyStr)
	if err == nil {
		nodeStrs = strings.Split(strings.TrimSpace(string(decoded)), "\n")
	} else {
		nodeStrs = strings.Split(string(body), "\n")
	}

	var nodes []NodeInfo
	for _, _node := range nodeStrs {
		node, err := url.QueryUnescape(strings.TrimSpace(_node))
		if node == "" {
			continue
		}
		if err != nil {
			log.Println("节点 decodeURIComponent 失败: " + err.Error())
			continue
		}
		var name string
		if idx := strings.Index(node, "#"); idx != -1 {
			name = node[idx+1:]
			node = node[:idx]
		}
		u, err := url.Parse(node)
		if err != nil || u.Scheme == "" {
			continue
		}
		nodes = append(nodes, TrojanNode{
			Name: fmt.Sprintf("%s#%s", name, tag),
			Url:  u,
		})
	}
	return nodes, nil
}

type TrojanNode struct {
	Name string
	Url  *url.URL
}

func (n TrojanNode) GetPluginOpts() *string {
	return new(string)
}

func (n TrojanNode) GetAllowInsecure() *bool {
	b := true
	return &b
}

func (n TrojanNode) GetName() string {
	return n.Name
}
func (n TrojanNode) GetScheme() string {
	return n.Url.Scheme
}
func (n TrojanNode) GetServer() string {
	return n.Url.Hostname()
}
func (n TrojanNode) GetPort() int {
	return *parsePort(n.Url.Port())
}
func (n TrojanNode) GetCipher() *string {
	return nil
}
func (n TrojanNode) GetPassword() string {
	return n.Url.User.Username()
}
func (n TrojanNode) GetPlugin() *string {
	s := n.Url.Query().Get("plugin")
	return &s
}
func (n TrojanNode) GetUDP() *bool {
	b := true
	return &b
}
func (n TrojanNode) GetUUID() *string {
	p := n.GetPassword()
	return &p
}
func (n TrojanNode) GetAlterID() *int {
	return nil
}
func (n TrojanNode) GetTLS() *bool {
	b := true
	return &b
}
func (n TrojanNode) GetSNI() *string {
	s := n.Url.Query().Get("peer")
	return &s
}
func (n TrojanNode) GetNetwork() *string {
	return nil
}
func (n TrojanNode) GetWSPath() *string {
	return nil
}
func (n TrojanNode) GetWSSHeaders() *map[string]string {
	return nil
}
