// +build proto
package test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"core"
	pb "test/proto/worker"
)

// runTestRemotely runs a single test against a remote worker and returns the output & any error.
// It assembles the output files in the target's temp directory.
// TODO(pebers): probably we should minimise the required file writing here and return the
//               information without writing it?
func runTestRemotely(state *core.BuildState, target *core.BuildTarget) ([]byte, error) {
	client, err := manager.GetClient(state.Config)
	if err != nil {
		return nil, err
	}
	timeout := target.TestTimeout
	if timeout == 0 {
		timeout = time.Duration(state.Config.Test.Timeout)
	}
	request := pb.TestRequest{
		Rule:     &pb.BuildLabel{PackageName: target.Label.PackageName, Name: target.Label.Name},
		Command:  target.GetTestCommand(),
		Coverage: state.NeedCoverage,
		TestName: state.TestArgs,
		Timeout:  int32(timeout.Seconds()),
		Labels:   target.Labels,
		NoOutput: target.NoTestOutput,
	}
	// Attach the test binary to the request
	b, err := ioutil.ReadFile(path.Join(target.OutDir(), target.Outputs()[0]))
	if err != nil {
		return nil, err
	}
	request.Binary = &pb.DataFile{Filename: target.Outputs()[0], Contents: b}
	// Attach its runtime files
	for _, datum := range target.Data {
		fullPaths := datum.FullPaths(state.Graph)
		for i, path := range datum.Paths(state.Graph) {
			b, err := ioutil.ReadFile(fullPaths[i])
			if err != nil {
				return nil, err
			}
			request.Data = append(request.Data, &pb.DataFile{Filename: path, Contents: b})
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := client.Test(ctx, &request)
	if err != nil {
		// N.B. we only get an error here if something failed structurally about the RPC - it is
		//      not an error if we communicate failure in the response.
		return nil, err
	}
	if !response.Success {
		return nil, fmt.Errorf("Failed to run test: %s", strings.Join(response.Messages, "\n"))
	}
	if !target.NoTestOutput {
		if err := ioutil.WriteFile(target.TestResultsFile(), response.Results, 0644); err != nil {
			return nil, err
		}
	}
	if len(response.Coverage) > 0 {
		if err := ioutil.WriteFile(target.TestCoverageFile(), response.Coverage, 0644); err != nil {
			return nil, err
		}
	}
	if !response.ExitSuccess {
		return response.Output, fmt.Errorf("process exited unsuccessfully")
	}
	return response.Output, nil
}

type clientManager struct {
	sync.RWMutex
	clients map[string]pb.TestWorkerClient
}

var manager = clientManager{clients: map[string]pb.TestWorkerClient{}}

// GetClient returns a client of one of our remote test workers.
// TODO(pebers): we lazy-initialise these, doing it eagerly at startup would probably be a
//               small optimisation so they're dialed and ready when we want to use them.
func (m *clientManager) GetClient(config *core.Configuration) (pb.TestWorkerClient, error) {
	address := config.Test.RemoteWorker[rand.Intn(len(config.Test.RemoteWorker))].String()
	m.RLock()
	client, present := m.clients[address]
	m.RUnlock()
	if present {
		return client, nil
	}
	// Need to initialise a new one.
	// Technically this is a little racy but there is little harm creating an extra client.
	m.Lock()
	defer m.Unlock()
	// TODO(pebers): Support secure connections here.
	conn, err := grpc.Dial(address, grpc.WithTimeout(5*time.Second), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client = pb.NewTestWorkerClient(conn)
	m.clients[address] = client
	return client, nil
}
