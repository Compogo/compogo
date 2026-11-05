package compogo

import (
	"os"

	cFlag "github.com/Compogo/compogo/flag"
	"github.com/Compogo/types/linker"
	flag "github.com/spf13/pflag"
)

var flagComponent = Component{
	Name:         "compogo.flags",
	Dependencies: nil,
	Steps: linker.NewLinker[Step, StepFunc](
		linker.Link[Step, StepFunc](Init, func(container Container) error {
			return container.Provides(
				func(c *Config) *flag.FlagSet {
					return flag.NewFlagSet(c.Name, flag.ContinueOnError)
				},
				func(flagSet *flag.FlagSet) cFlag.FlagSet {
					return flagSet
				},
			)
		}),
		linker.Link[Step, StepFunc](Configuration, func(container Container) error {
			return container.Invoke(func(flagSet *flag.FlagSet) error {
				return flagSet.Parse(os.Args[1:])
			})
		}),
	),
}
