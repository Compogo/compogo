package flag

import "net"

type IPNet interface {
	IPNetP(name, shorthand string, value net.IPNet, usage string) *net.IPNet
	IPNet(name string, value net.IPNet, usage string) *net.IPNet
	IPNetVarP(p *net.IPNet, name, shorthand string, value net.IPNet, usage string)
	IPNetVar(p *net.IPNet, name string, value net.IPNet, usage string)
	GetIPNet(name string) (net.IPNet, error)
}
