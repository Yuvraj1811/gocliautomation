package gitops

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	
)

type Kustomization struct {
	APIVersion string `yaml:"apiVersion"`
	Kind string `yaml:"kind"`
	Resources []string `yaml:"resources,omitempty"`
	Images []Image `yaml:"images,omitempty"`

}

type Image struct {
	Name string `yaml:"name"`
	NewTag string `yaml:"newTag,omitempty"`
}

func LoadKustomization(path string) (*Kustomization, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var k Kustomization
	if err := yaml.Unmarshal(data, &k); err != nil {
		return nil, err
	}

	return &k, nil
}

func UpdateImageTag(k *Kustomization, app string, newTag string) error {
	for i, img := range k.Images {
		if img.Name == app {
			k.Images[i].NewTag = newTag
			return nil
		}
	}

	return fmt.Errorf("image %s not found in kustomization", app)
}

func SaveKustomization(path string, k *Kustomization) error {
	data, err := yaml.Marshal(k)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func RenderYAML(k *Kustomization) ([]byte, error) {
	return yaml.Marshal(k)
}