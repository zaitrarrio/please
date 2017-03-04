package test

import "core"
import "fmt"

// runTestRemotely is a stub used during initial bootstrap when protobufs aren't available.
func runTestRemotely(state *core.BuildState, target *core.BuildTarget) ([]byte, error) {
	return nil, fmt.Errorf("Cannot run test remotely, remote running is not compiled")
}
