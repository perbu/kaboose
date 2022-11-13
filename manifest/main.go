package manifest

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

type k8sYaml interface {
	Kustomization | Manifest | LocalConfig
}

// Load is a generic function to load a yaml file.
func Load[T k8sYaml](file string, target *T) error {
	fh, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %w", file, err)
	}
	defer fh.Close()
	err = loadFromReader(fh, target)
	if err != nil {
		return fmt.Errorf("failed to load file '%s': %w", file, err)
	}
	// fix the paths if needed:

	return nil
}

func loadFromReader[T k8sYaml](r io.Reader, target *T) error {
	// parse the yaml into the Manifest struct
	content, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read manifest file: %w", err)
	}
	err = yaml.Unmarshal(content, target)
	if err != nil {
		return fmt.Errorf("failed to parse manifest file: %w", err)
	}
	return nil
}

func expandTilde(path string) string {
	if path[0] == '~' {
		return filepath.Join(os.Getenv("HOME"), path[1:])
	}
	return path
}
