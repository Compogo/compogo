package flag

import "net"

type IP interface {
	IPP(name, shorthand string, value net.IP, usage string) *net.IP
	IP(name string, value net.IP, usage string) *net.IP
	IPVarP(p *net.IP, name, shorthand string, value net.IP, usage string)
	IPVar(p *net.IP, name string, value net.IP, usage string)
	GetIP(name string) (net.IP, error)
}
