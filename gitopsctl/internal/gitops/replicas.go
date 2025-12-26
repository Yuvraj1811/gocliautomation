package gitops

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type replicaPatch struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
	} `yaml:"spec"`
}

func UpdateReplicas(app, env string, replicas int, dryRun bool) (string, error) {
	overlayDir := filepath.Join("apps", app, "overlays", env)
	patchPath := filepath.Join(overlayDir, "replicas.yaml")

	if _, err := os.Stat(overlayDir); err != nil {
		return "", fmt.Errorf("overlay not found: %s", overlayDir)
	}

	var p replicaPatch
	p.APIVersion = "apps/v1"
	p.Kind = "Deployment"
	p.Metadata.Name = app
	p.Spec.Replicas = replicas

	out, _ := yaml.Marshal(&p)

	if !dryRun {
		_ = os.WriteFile(patchPath, out, 0644)
		_ = ensurePatchInKustomization(overlayDir, "replicas.yaml")
	}

	return patchPath, nil
}

func ensurePatchInKustomization(dir, patch string) error {
	path := filepath.Join(dir, "kustomization.yaml")

	type kustomization struct {
		Patches []string `yaml:"patchesStrategicMerge,omitempty"`
	}

	var k kustomization

	if data, err := os.ReadFile(path); err == nil {
		_ = yaml.Unmarshal(data, &k)
	}

	for _, p := range k.Patches {
		if p == patch {
			goto WRITE
		}
	}

	k.Patches = append(k.Patches, patch)

WRITE:
	out, _ := yaml.Marshal(&k)
	return os.WriteFile(path, out, 0644)
}
