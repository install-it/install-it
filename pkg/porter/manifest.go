package porter

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"time"
)

type Manifest struct {
	FormatVersion int       `json:"format_version"`
	ExportedAt    time.Time `json:"exported_at"`
}

const CurrentFormatVersion = 1

func newManifest() Manifest {
	return Manifest{
		FormatVersion: 1,
		ExportedAt:    time.Now(),
	}
}

func readManifest(zr *zip.ReadCloser) (Manifest, error) {
	for _, f := range zr.File {
		if f.Name == "manifest.json" {
			rc, err := f.Open()
			if err != nil {
				return Manifest{}, fmt.Errorf("porter: cannot open manifest.json: %w", err)
			}
			defer rc.Close()

			var m Manifest
			if err := json.NewDecoder(rc).Decode(&m); err != nil {
				return Manifest{}, fmt.Errorf("porter: invalid manifest.json: %w", err)
			}
			return m, nil
		}
	}
	return Manifest{}, fmt.Errorf("porter: manifest.json not found in archive")
}

func writeManifest(zw *zip.Writer, m Manifest) error {
	entry, err := zw.Create("manifest.json")
	if err != nil {
		return fmt.Errorf("porter: cannot create manifest.json entry: %w", err)
	}
	if err := json.NewEncoder(entry).Encode(m); err != nil {
		return fmt.Errorf("porter: cannot encode manifest.json: %w", err)
	}
	return nil
}
