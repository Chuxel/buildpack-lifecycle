package layer

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/buildpacks/imgutil"
	"github.com/pkg/errors"

	"github.com/buildpacks/lifecycle/api"
	"github.com/buildpacks/lifecycle/buildpack"
	io2 "github.com/buildpacks/lifecycle/internal/io"
	"github.com/buildpacks/lifecycle/launch"
	"github.com/buildpacks/lifecycle/layers"
)

//go:generate mockgen -package testmock -destination testmock/sbom_restorer.go github.com/buildpacks/lifecycle/internal/layer SBOMRestorer
type SBOMRestorer interface {
	RestoreFromPrevious(image imgutil.Image, layerDigest string) error
	RestoreFromCache(cache Cache, layerDigest string) error
	RestoreToBuildpackLayers(detectedBps []buildpack.GroupBuildpack) error
}

type Cache interface {
	RetrieveLayer(sha string) (io.ReadCloser, error)
}

type SBOMRestorerOpts struct {
	LayersDir string
	Logger    Logger
	Nop       bool
}

func NewSBOMRestorer(opts SBOMRestorerOpts, platformAPI *api.Version) SBOMRestorer {
	if opts.Nop || platformAPI.LessThan("0.8") {
		return &NopSBOMRestorer{}
	}
	return &DefaultSBOMRestorer{
		layersDir: opts.LayersDir,
		logger:    opts.Logger,
	}
}

type DefaultSBOMRestorer struct {
	layersDir string
	logger    Logger
}

func (r *DefaultSBOMRestorer) RestoreFromPrevious(image imgutil.Image, layerDigest string) error {
	// Sanity check to prevent panic.
	if image == nil {
		return errors.Errorf("restoring layer: previous image not found for %q", layerDigest)
	}

	if !image.Found() || layerDigest == "" {
		return nil
	}
	r.logger.Infof("Restoring data for SBOM from previous image")

	r.logger.Debugf("Retrieving previous image SBOM layer for %q", layerDigest)
	rc, err := image.GetLayer(layerDigest)
	if err != nil {
		return err
	}
	defer rc.Close()

	return layers.Extract(rc, "")
}

func (r *DefaultSBOMRestorer) RestoreFromCache(cache Cache, layerDigest string) error {
	// Sanity check to prevent panic.
	if cache == nil {
		return errors.New("restoring layer: cache not provided")
	}
	r.logger.Debugf("Retrieving SBOM layer data for %q", layerDigest)

	rc, err := cache.RetrieveLayer(layerDigest)
	if err != nil {
		return err
	}
	defer rc.Close()

	return layers.Extract(rc, "")
}

func (r *DefaultSBOMRestorer) RestoreToBuildpackLayers(detectedBps []buildpack.GroupBuildpack) error {
	var (
		cacheDir  = filepath.Join(r.layersDir, "sbom", "cache")
		launchDir = filepath.Join(r.layersDir, "sbom", "launch")
	)
	defer os.RemoveAll(filepath.Join(r.layersDir, "sbom"))

	if err := filepath.Walk(cacheDir, r.restoreSBOMFunc(detectedBps, "cache")); err != nil {
		return err
	}

	return filepath.Walk(launchDir, r.restoreSBOMFunc(detectedBps, "launch"))
}

func (r *DefaultSBOMRestorer) restoreSBOMFunc(detectedBps []buildpack.GroupBuildpack, bomType string) func(path string, info fs.FileInfo, err error) error {
	var bomRegex *regexp.Regexp

	if runtime.GOOS == "windows" {
		bomRegex = regexp.MustCompile(fmt.Sprintf(`%s\\(.+)\\(.+)\\(sbom.+json)`, bomType))
	} else {
		bomRegex = regexp.MustCompile(fmt.Sprintf(`%s/(.+)/(.+)/(sbom.+json)`, bomType))
	}

	return func(path string, info fs.FileInfo, err error) error {
		if info == nil || !info.Mode().IsRegular() {
			return nil
		}

		matches := bomRegex.FindStringSubmatch(path)
		if len(matches) != 4 {
			return nil
		}

		var (
			bpID      = matches[1]
			layerName = matches[2]
			fileName  = matches[3]
			destDir   = filepath.Join(r.layersDir, bpID)
		)

		// don't try to restore sbom files when the bp layers directory doesn't exist
		// this can happen when there are sbom files for launch but the cache is empty
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			return nil
		}

		if !r.contains(detectedBps, bpID) {
			return nil
		}

		return io2.Copy(path, filepath.Join(destDir, fmt.Sprintf("%s.%s", layerName, fileName)))
	}
}

func (r *DefaultSBOMRestorer) contains(detectedBps []buildpack.GroupBuildpack, id string) bool {
	for _, bp := range detectedBps {
		if launch.EscapeID(bp.ID) == id {
			return true
		}
	}
	return false
}

type NopSBOMRestorer struct{}

func (r *NopSBOMRestorer) RestoreFromPrevious(_ imgutil.Image, _ string) error {
	return nil
}

func (r *NopSBOMRestorer) RestoreFromCache(_ Cache, _ string) error {
	return nil
}

func (r *NopSBOMRestorer) RestoreToBuildpackLayers(_ []buildpack.GroupBuildpack) error {
	return nil
}
