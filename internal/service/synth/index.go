package synth

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const indexFileName = "index.yaml"

type synthIndex map[string]string // hash → basename under output dir

func indexPath(outputDir string) string {
	return filepath.Join(outputDir, indexFileName)
}

func loadIndex(outputDir string) (synthIndex, error) {
	path := indexPath(outputDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return synthIndex{}, nil
		}
		return nil, err
	}
	var idx synthIndex
	if err := yaml.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	if idx == nil {
		idx = synthIndex{}
	}
	return idx, nil
}

func saveIndex(outputDir string, idx synthIndex) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(idx)
	if err != nil {
		return err
	}
	return os.WriteFile(indexPath(outputDir), data, 0o644)
}

// CleanIndex removes entries whose WAV files are missing under outputDir.
func CleanIndex(outputDir string) error {
	outputDir = filepath.Clean(outputDir)
	idx, err := loadIndex(outputDir)
	if err != nil {
		return err
	}
	changed := false
	for hash, name := range idx {
		if name == "" || !fileExists(filepath.Join(outputDir, name)) {
			delete(idx, hash)
			changed = true
		}
	}
	if !changed {
		return nil
	}
	return saveIndex(outputDir, idx)
}

func lookupCache(outputDir, hash string) (absPath string, ok bool, err error) {
	idx, err := loadIndex(outputDir)
	if err != nil {
		return "", false, err
	}
	name, found := idx[hash]
	if !found || name == "" {
		return "", false, nil
	}
	abs := filepath.Join(outputDir, name)
	if !fileExists(abs) {
		delete(idx, hash)
		if err := saveIndex(outputDir, idx); err != nil {
			return "", false, err
		}
		return "", false, nil
	}
	return abs, true, nil
}

func putCache(outputDir, hash, basename string) error {
	idx, err := loadIndex(outputDir)
	if err != nil {
		return err
	}
	idx[hash] = basename
	return saveIndex(outputDir, idx)
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	tmp := dst + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, dst); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("copy cache: %w", err)
	}
	return nil
}
