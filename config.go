package compogo

import (
	"os"

	"github.com/Compogo/compogo/configurator"
)

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

func Configuration(config *Config, configurator configurator.Configurator) *Config {
	config.PID = uint64(os.Getpid())

	if config.Cluster == "" || config.Cluster == ClusterDefault {
		configurator.SetDefault(ClusterFieldName, ClusterDefault)
		config.Cluster = configurator.GetString(ClusterFieldName)
	}

	if config.Namespace == "" || config.Namespace == NamespaceDefault {
		configurator.SetDefault(NamespaceFieldName, NamespaceDefault)
		config.Namespace = configurator.GetString(NamespaceFieldName)
	}

	if config.ContainerId == "" || config.ContainerId == ContainerIdDefault {
		configurator.SetDefault(ContainerIdFieldName, ContainerIdDefault)
		config.ContainerId = configurator.GetString(ContainerIdFieldName)
	}

	if config.ContainerName == "" || config.ContainerName == ContainerNameDefault {
		configurator.SetDefault(ContainerNameFieldName, ContainerNameDefault)
		config.ContainerName = configurator.GetString(ContainerNameFieldName)
	}

	if config.Hostname == "" || config.Hostname == HostnameDefault {
		configurator.SetDefault(HostnameFieldName, HostnameDefault)
		config.Hostname = configurator.GetString(HostnameFieldName)
	}

	return config
}
