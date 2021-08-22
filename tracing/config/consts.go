package config

const (
	ServerUrl = "http://localhost:%d/%s"

	MethodFormat  = "format"
	MethodPublish = "publish"
)

var RoleToPort = map[string]int32{
	MethodFormat:  8080,
	MethodPublish: 8081,
}
