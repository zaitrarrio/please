package worker

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"golang.org/x/net/context"
	"gopkg.in/op/go-logging.v1"

	"core"
	"test"
	pb "test/proto/worker"
)

var log = logging.MustGetLogger("please_test_worker")

// A Worker implements our remote test worker protocol.
type Worker struct{}

func (w *Worker) Test(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	// Build a sufficient representation of the target that we can call the preexisting code for this.
	state := core.NewBuildState(1, nil, 2, core.DefaultConfiguration())
	target := core.NewBuildTarget(core.BuildLabel{PackageName: req.Rule.PackageName, Name: req.Rule.Name})
	target.TestCommand = req.Command
	state.NeedCoverage = req.Coverage
	state.TestArgs = req.TestName
	target.NoTestOutput = req.NoOutput
	dir := target.TestDir()

	// From here on, if anything goes wrong, remove the temp directory.
	defer w.cleanup(dir)

	// Write the required files
	if err := w.writeFile(target, req.Binary); err != nil {
		return nil, err
	}
	for _, df := range req.Data {
		if err := w.writeFile(target, df); err != nil {
			return nil, err
		}
	}

	out, err := test.RunTest(state, target)
	response := &pb.TestResponse{
		Rule:        req.Rule,
		Success:     true,
		ExitSuccess: err == nil,
		Output:      out,
	}
	// Read test results file & coverage file
	b, err := ioutil
}

// writeFile is a convenience wrapper to create one of the test files & any needed directories.
func (w *Worker) writeFile(dir string, df *pb.DataFile) error {
	// Note that we must use the target's test dir for the normal code to work.
	// We might want a slightly flatter structure here but on the other hand it makes it easy to
	// run these workers from the root of a Please repo without polluting anything else.
	filename := path.Join(dir, df.Filename)
	if err := os.MkdirAll(path.Dir(filename), core.DirPermissions); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, df.Contents, 0755)
}

// cleanup removes the temporary test directory.
// Note that we can't return an error since it's deferred (and strictly doesn't represent an error to
// run the test anyway).
func (w *Worker) cleanup(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Error("Failed to remove temporary test directory: %s", err)
	}
}
