package subscribes

//func ParseByBase64(body string) ([]string, error) {
//	var nodeStrs []string
//
//	bodyStr := strings.TrimSpace(string(body))
//	missingPadding := len(bodyStr) % 4
//	if missingPadding != 0 {
//		bodyStr += strings.Repeat("=", 4-missingPadding)
//	}
//	decoded, err := base64.StdEncoding.DecodeString(bodyStr)
//	if err == nil {
//		nodeStrs = strings.Split(strings.TrimSpace(string(decoded)), "\n")
//	} else {
//		nodeStrs = strings.Split(string(body), "\n")
//	}
//
//	var nodes []TrojanNode
//	for _, _node := range nodeStrs {
//		node, err := url.QueryUnescape(strings.TrimSpace(_node))
//		if node == "" {
//			continue
//		}
//		if err != nil {
//			log.Println("节点 decodeURIComponent 失败: " + err.Error())
//			continue
//		}
//		var name string
//		if idx := strings.Index(node, "#"); idx != -1 {
//			name = node[idx+1:]
//			node = node[:idx]
//		}
//		u, err := url.Parse(node)
//		if err != nil || u.Scheme == "" {
//			continue
//		}
//		nodes = append(nodes, TrojanNode{
//			Name: name,
//			Url:  u,
//		})
//	}
//	return nodeStrs, nil
//}
//
//type TrojanNode struct {
//	Name string
//	Url  *url.URL
//}

//
//func (n TrojanNode) GetName() string {
//	return n.Name
//}
//func (n TrojanNode) GetScheme() string {
//	return n.Url.Scheme
//}
//func (n TrojanNode) GetServer() string {
//	return n.Url.Host
//}
//func (n TrojanNode) GetPort() int {
//	return *parsePort(n.Url.Port())
//}
//func (n TrojanNode) GetCipher() *string {
//	return nil
//}
//func (n TrojanNode) GetPassword() *string {
//	s := n.Url.User.Username()
//	return &s
//}
//func (n TrojanNode) GetPlugin() *string {
//	s := n.Url.Query().Get("plugin")
//	return &s
//}
//func (n TrojanNode) GetPluginOpts() *map[string]interface{} {
//	return nil
//}
//func (n TrojanNode) GetUDP() *bool {
//	b := true
//	return &b
//}
//func (n TrojanNode) GetUUID() *string {
//	return n.GetPassword()
//}
//func (n TrojanNode) GetAlterID() *int {
//	return nil
//}
//func (n TrojanNode) GetTLS() *bool {
//	b := true
//	return &b
//}
//func (n TrojanNode) GetSNI() *string {
//	return &n.Url.Query().Get("")
//}
//func (n TrojanNode) GetNetwork() *string {
//	return nil
//}
//func (n TrojanNode) GetWSPath() *string {
//	return nil
//}
//func (n TrojanNode) GetWSSHeaders() *map[string]string {
//	return nil
//}
