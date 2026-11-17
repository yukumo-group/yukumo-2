package synth

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"yukumo-script/core/port"
)

const indexKey = "index.yaml"

type synthIndex map[string]string // hash → artifact key

func loadIndex(a port.ArtifactStore) (synthIndex, error) {
	if a == nil {
		return synthIndex{}, nil
	}
	data, err := a.Read(indexKey)
	if err != nil {
		if errors.Is(err, port.ErrNotFound) {
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

func saveIndex(a port.ArtifactStore, idx synthIndex) error {
	if a == nil {
		return nil
	}
	data, err := yaml.Marshal(idx)
	if err != nil {
		return err
	}
	return a.Write(indexKey, data)
}

// cleanIndex removes entries whose artifact keys are missing.
func cleanIndex(a port.ArtifactStore) error {
	idx, err := loadIndex(a)
	if err != nil {
		return err
	}
	changed := false
	for hash, key := range idx {
		if key == "" {
			delete(idx, hash)
			changed = true
			continue
		}
		ok, err := a.Exists(key)
		if err != nil {
			return err
		}
		if !ok {
			delete(idx, hash)
			changed = true
		}
	}
	if !changed {
		return nil
	}
	return saveIndex(a, idx)
}

func lookupCache(a port.ArtifactStore, hash string) (key string, ok bool, err error) {
	idx, err := loadIndex(a)
	if err != nil {
		return "", false, err
	}
	key, found := idx[hash]
	if !found || key == "" {
		return "", false, nil
	}
	exists, err := a.Exists(key)
	if err != nil {
		return "", false, err
	}
	if !exists {
		delete(idx, hash)
		if err := saveIndex(a, idx); err != nil {
			return "", false, err
		}
		return "", false, nil
	}
	return key, true, nil
}

// putCache records hash → key; an external output keeps its own key and gets a
// cache copy under cacheKey.
func putCache(a port.ArtifactStore, hash, key string, external bool, wav []byte) error {
	if a == nil {
		return nil
	}
	idx, err := loadIndex(a)
	if err != nil {
		return err
	}
	if external {
		key = cacheKey(hash)
		if err := a.Write(key, wav); err != nil {
			return err
		}
	}
	idx[hash] = key
	return saveIndex(a, idx)
}

func putCacheRemove(a port.ArtifactStore, hash string) error {
	if a == nil {
		return nil
	}
	idx, err := loadIndex(a)
	if err != nil {
		return err
	}
	delete(idx, hash)
	return saveIndex(a, idx)
}

func cacheKey(hash string) string {
	return "cache-" + hash + ".wav"
}

// uniqueKey returns base.ext or base-N.ext within the artifact store.
func uniqueKey(a port.ArtifactStore, base, ext string) (string, error) {
	if len(ext) > 0 && ext[0] != '.' {
		ext = "." + ext
	}
	candidate := base + ext
	exists, err := a.Exists(candidate)
	if err != nil {
		return "", err
	}
	if !exists {
		return candidate, nil
	}
	for i := 1; i < 10000; i++ {
		candidate = fmt.Sprintf("%s-%d%s", base, i, ext)
		exists, err := a.Exists(candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not allocate unique key for %s", base)
}
