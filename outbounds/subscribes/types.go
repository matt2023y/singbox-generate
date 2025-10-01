package subscribes

import "strconv"

type NodeInfo interface {
	GetName() string   // 节点名称
	GetScheme() string // 协议类型

	GetServer() string
	GetPort() int

	GetCipher() *string // 加密方式

	GetPlugin() *string
	GetPluginOpts() *string

	GetUDP() *bool

	GetAllowInsecure() *bool

	GetPassword() string
	GetUUID() *string

	GetAlterID() *int
	GetTLS() *bool
	GetSNI() *string
	GetNetwork() *string
	GetWSPath() *string
	GetWSSHeaders() *map[string]string
}

func parsePort(s string) *int {
	if s == "" {
		return nil
	}
	port, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &port
}

func firstParam(s1, s2 string) string {
	if s1 == s2 {
		return s1
	}
	if s1 == "" {
		return s2
	}
	if s2 == "" {
		return s1
	}
	return s1
}

func toBool(s string) bool {
	if s == "false" {
		return false
	}
	if s == "0" {
		return false
	}
	return true
}
