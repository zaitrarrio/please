package cache

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"core"
)

func init() {
	// Move this directory from test_data to somewhere local.
	if err := os.Rename("src/cache/test_data/plz-out", "plz-out"); err != nil {
		log.Fatalf("Failed to prepare test directory: %s\n", err)
	}
}

func TestStoreOK(t *testing.T) {
	cache, mockCache, target := setUp("10GB")
	cache.Store(target, nil)
	assert.True(t, mockCache.called)
}

func TestStoreRejected(t *testing.T) {
	cache, mockCache, target := setUp("1B")
	cache.Store(target, nil)
	assert.False(t, mockCache.called)
}

func TestStoreExtraOK(t *testing.T) {
	cache, mockCache, target := setUp("10GB")
	cache.StoreExtra(target, nil, "testfile")
	assert.True(t, mockCache.called)
}

func TestStoreExtraRejected(t *testing.T) {
	cache, mockCache, target := setUp("1B")
	cache.StoreExtra(target, nil, "testfile")
	assert.False(t, mockCache.called)
}

func TestRetrieveAlwaysCalled(t *testing.T) {
	cache, mockCache, target := setUp("1B")
	assert.True(t, cache.Retrieve(target, nil))
	assert.True(t, mockCache.called)
}

func TestRetrieveExtraAlwaysCalled(t *testing.T) {
	cache, mockCache, target := setUp("1B")
	assert.True(t, cache.RetrieveExtra(target, nil, "testfile"))
	assert.True(t, mockCache.called)
}
func TestCleanAlwaysCalled(t *testing.T) {
	cache, mockCache, target := setUp("1B")
	cache.Clean(target)
	assert.True(t, mockCache.called)
}

type mockCache struct {
	called bool
}

func (c *mockCache) Store(target *core.BuildTarget, key []byte) {
	c.called = true
}

func (c *mockCache) StoreExtra(target *core.BuildTarget, key []byte, file string) {
	c.called = true
}

func (c *mockCache) Retrieve(target *core.BuildTarget, key []byte) bool {
	c.called = true
	return true
}

func (c *mockCache) RetrieveExtra(target *core.BuildTarget, key []byte, file string) bool {
	c.called = true
	return true
}

func (c *mockCache) Clean(target *core.BuildTarget) {
	c.called = true
}

func (c *mockCache) Shutdown() {}

func setUp(sizeLimit string) (core.Cache, *mockCache, *core.BuildTarget) {
	target := core.NewBuildTarget(core.NewBuildLabel("pkg/name", "label_name"))
	target.AddOutput("testfile2")
	target.AddOutput("testfile3")
	target.AddOutput("testfile4")
	target.BuildDuration = time.Second
	config := core.DefaultConfiguration()
	config.Cache.MaxSizeFactor.UnmarshalFlag(sizeLimit)
	cache := &mockCache{}
	return newSizeFactorLimit(cache, config), cache, target
}
