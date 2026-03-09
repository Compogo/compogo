package compogo

import (
	"os"
	"time"

	"github.com/Compogo/compogo/configurator"
)

const (
	ClusterFieldName       = "cluster"
	NamespaceFieldName     = "namespace"
	ContainerIdFieldName   = "container.id"
	ContainerNameFieldName = "container.name"
	HostnameFieldName      = "hostname"

	InitDurationFieldName = "compogo.duration.init"

	PreRunDurationFieldName  = "compogo.duration.run.pre"
	RunDurationFieldName     = "compogo.duration.run.run"
	PostRunDurationFieldName = "compogo.duration.run.post"

	PreWaitDurationFieldName  = "compogo.duration.wait.pre"
	PostWaitDurationFieldName = "compogo.duration.wait.post"

	PreStopDurationFieldName  = "compogo.duration.stop.pre"
	StopDurationFieldName     = "compogo.duration.stop.stop"
	PostStopDurationFieldName = "compogo.duration.stop.post"

	ClusterDefault       = "unknown-cluster"
	NamespaceDefault     = "unknown-namespace"
	ContainerIdDefault   = "unknown-container_id"
	ContainerNameDefault = "unknown-container_name"
	HostnameDefault      = "unknown-hostname"

	InitDurationDefault = 100 * time.Millisecond

	PreRunDurationDefault  = 100 * time.Millisecond
	RunDurationDefault     = 100 * time.Millisecond
	PostRunDurationDefault = 100 * time.Millisecond

	PreWaitDurationDefault  = 100 * time.Millisecond
	PostWaitDurationDefault = 100 * time.Millisecond

	PreStopDurationDefault  = 100 * time.Millisecond
	StopDurationDefault     = 100 * time.Millisecond
	PostStopDurationDefault = 100 * time.Millisecond
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

	// Duration
	InitDuration time.Duration

	PreRunDuration  time.Duration
	RunDuration     time.Duration
	PostRunDuration time.Duration

	PreWaitDuration  time.Duration
	PostWaitDuration time.Duration

	PreStopDuration  time.Duration
	StopDuration     time.Duration
	PostStopDuration time.Duration
}

func NewConfig() *Config {
	return &Config{
		InitDuration:     InitDurationDefault,
		PreRunDuration:   PreRunDurationDefault,
		RunDuration:      RunDurationDefault,
		PostRunDuration:  PostRunDurationDefault,
		PreWaitDuration:  PreWaitDurationDefault,
		PostWaitDuration: PostWaitDurationDefault,
		PreStopDuration:  PreStopDurationDefault,
		StopDuration:     StopDurationDefault,
		PostStopDuration: PostStopDurationDefault,
	}
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

	if config.InitDuration == 0 || config.InitDuration == InitDurationDefault {
		configurator.SetDefault(InitDurationFieldName, InitDurationDefault)
		config.InitDuration = configurator.GetDuration(InitDurationFieldName)
	}

	if config.PreRunDuration == 0 || config.PreRunDuration == PreRunDurationDefault {
		configurator.SetDefault(PreRunDurationFieldName, PreRunDurationDefault)
		config.PreRunDuration = configurator.GetDuration(PreRunDurationFieldName)
	}

	if config.RunDuration == 0 || config.RunDuration == RunDurationDefault {
		configurator.SetDefault(RunDurationFieldName, RunDurationDefault)
		config.RunDuration = configurator.GetDuration(RunDurationFieldName)
	}

	if config.PostRunDuration == 0 || config.PostRunDuration == PostRunDurationDefault {
		configurator.SetDefault(PostRunDurationFieldName, PostRunDurationDefault)
		config.PostRunDuration = configurator.GetDuration(PostRunDurationFieldName)
	}

	if config.PreWaitDuration == 0 || config.PreWaitDuration == PreWaitDurationDefault {
		configurator.SetDefault(PreWaitDurationFieldName, PreWaitDurationDefault)
		config.PreWaitDuration = configurator.GetDuration(PreWaitDurationFieldName)
	}

	if config.PostWaitDuration == 0 || config.PostWaitDuration == PostWaitDurationDefault {
		configurator.SetDefault(PostWaitDurationFieldName, PostWaitDurationDefault)
		config.PostWaitDuration = configurator.GetDuration(PostWaitDurationFieldName)
	}

	if config.PreStopDuration == 0 || config.PreStopDuration == PreStopDurationDefault {
		configurator.SetDefault(PreStopDurationFieldName, PreStopDurationDefault)
		config.PreStopDuration = configurator.GetDuration(PreStopDurationFieldName)
	}

	if config.StopDuration == 0 || config.StopDuration == StopDurationDefault {
		configurator.SetDefault(StopDurationFieldName, StopDurationDefault)
		config.StopDuration = configurator.GetDuration(StopDurationFieldName)
	}

	if config.PostStopDuration == 0 || config.PostStopDuration == PostStopDurationDefault {
		configurator.SetDefault(PostStopDurationFieldName, PostStopDurationDefault)
		config.PostStopDuration = configurator.GetDuration(PostStopDurationFieldName)
	}

	return config
}
