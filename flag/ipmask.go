package flag

import "net"

type IPMask interface {
	IPMaskP(name, shorthand string, value net.IPMask, usage string) *net.IPMask
	IPMask(name string, value net.IPMask, usage string) *net.IPMask
	IPMaskVarP(p *net.IPMask, name, shorthand string, value net.IPMask, usage string)
	IPMaskVar(p *net.IPMask, name string, value net.IPMask, usage string)
	GetIPv4Mask(name string) (net.IPMask, error)
}
