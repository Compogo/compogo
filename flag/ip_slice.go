package flag

import "net"

type IPSlice interface {
	IPSliceP(name, shorthand string, value []net.IP, usage string) *[]net.IP
	IPSlice(name string, value []net.IP, usage string) *[]net.IP
	IPSliceVarP(p *[]net.IP, name, shorthand string, value []net.IP, usage string)
	IPSliceVar(p *[]net.IP, name string, value []net.IP, usage string)
	GetIPSlice(name string) ([]net.IP, error)
}
