package component

import (
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
)

// StepFunc defines a function that executes at a specific lifecycle step.
// It receives the DI container which should already contain all initialized dependencies.
type StepFunc func(container container.Container) error

// BindFlags defines a function for registering command-line flags.
// Called during the flag binding phase before any step execution.
type BindFlags func(flagSet flag.FlagSet, container container.Container) error

// Components is a collection of Component pointers for easier grouping.
type Components []*Component

// Component represents a modular piece of application logic with a defined lifecycle.
// Each component can declare dependencies and hook into various execution steps.
type Component struct {
	// Name component in core logs
	Name string

	// Dependencies lists other components that must be initialized before this one
	Dependencies Components

	// Init performs initial setup and registers services in the container
	Init StepFunc
	// BindFlags registers component-specific command-line flags
	BindFlags BindFlags

	// PreRun executes before the main Run step
	PreRun StepFunc
	// Run contains the main component logic
	Run StepFunc
	// PostRun executes after the main Run step
	PostRun StepFunc

	// PreWait executes before entering wait state
	PreWait StepFunc
	// Wait typically blocks for signals or async operations
	Wait StepFunc
	// PostWait executes after wait state
	PostWait StepFunc

	// PreStop executes before shutdown
	PreStop StepFunc
	// Stop performs cleanup and graceful shutdown
	Stop StepFunc
	// PostStop executes after shutdown is complete
	PostStop StepFunc
}
