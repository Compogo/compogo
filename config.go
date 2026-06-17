package compogo

const (
	ClusterFieldName       = "cluster"
	NamespaceFieldName     = "namespace"
	ContainerIdFieldName   = "container.id"
	ContainerNameFieldName = "container.name"
	HostnameFieldName      = "hostname"

	ClusterDefault       = "unknown-cluster"
	NamespaceDefault     = "unknown-namespace"
	ContainerIdDefault   = "unknown-container_id"
	ContainerNameDefault = "unknown-container_name"
	HostnameDefault      = "unknown-hostname"
)

type Config struct {
	Name string
	PID  uint64

	// k8s
	Cluster       string
	Namespace     string
	ContainerId   string
	ContainerName string
	Hostname      string
}

func NewConfig() *Config {
	return &Config{}
}
