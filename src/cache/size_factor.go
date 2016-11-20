package cache

import (
	"os"
	"path"
	"path/filepath"

	"core"
)

// A sizeFactorLimit is a wrapper around a real cache which limits which artifacts to
// cache based on the ratio of artifact size to the time they took to build.
// The scheme here is to avoid caching artifacts that are large and relatively easy
// to rebuild.
type sizeFactorLimit struct {
	realCache core.Cache
	limit     int
}

func newSizeFactorLimit(realCache core.Cache, config *core.Configuration) core.Cache {
	return &sizeFactorLimit{
		realCache: realCache,
		limit:     int(config.Cache.MaxSizeFactor),
	}
}

func (c *sizeFactorLimit) Store(target *core.BuildTarget, key []byte) {
	if c.shouldStore(target) {
		c.realCache.Store(target, key)
	}
}

func (c *sizeFactorLimit) StoreExtra(target *core.BuildTarget, key []byte, file string) {
	// TODO(pebers): Maybe we should remember the answer we came up with in Store()
	//               rather than recalculating here?
	if c.shouldStore(target) {
		c.realCache.StoreExtra(target, key, file)
	}
}

func (c *sizeFactorLimit) Retrieve(target *core.BuildTarget, key []byte) bool {
	return c.realCache.Retrieve(target, key)
}

func (c *sizeFactorLimit) RetrieveExtra(target *core.BuildTarget, key []byte, file string) bool {
	return c.realCache.RetrieveExtra(target, key, file)
}

func (c *sizeFactorLimit) Clean(target *core.BuildTarget) {
	c.realCache.Clean(target)
}

func (c *sizeFactorLimit) Shutdown() {}

func (c *sizeFactorLimit) sizeFactor(target *core.BuildTarget) int {
	duration := int64(target.BuildDuration.Seconds())
	if duration == 0 {
		// Avoid divide by zero later on, if target took basically no time to build
		// then flooring it at 1 second seems pretty reasonable.
		duration = 1
	}
	var size int64
	for out := range cacheArtifacts(target) {
		filepath.Walk(path.Join(core.RepoRoot, target.OutDir(), out), func(name string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			size += info.Size()
			return nil
		})
	}
	return int(size / duration)
}

func (c *sizeFactorLimit) shouldStore(target *core.BuildTarget) bool {
	factor := c.sizeFactor(target)
	if factor < c.limit {
		log.Debug("Will store %s in cache, size factor %d under limit", target.Label, factor)
		return true
	}
	log.Info("Not storing %s in cache, size factor %d exceeds limit", target.Label, factor)
	return false
}
