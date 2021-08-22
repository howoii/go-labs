package config

const (
	ServerUrl = "http://localhost:%d/%s"

	MethodFormatter = "formatter"
	MethodPublisher = "publisher"
)

var RoleToPort = map[string]int32{
	MethodFormatter: 8080,
	MethodPublisher: 8081,
}
