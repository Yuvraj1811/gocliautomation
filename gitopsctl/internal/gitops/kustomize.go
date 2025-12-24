package gitops

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// UpdateImageTagInFile updates only the newTag of the image in the kustomization.yaml
func UpdateImageTagInFile(path, image string, dryRun bool) ([]byte, error) {
	
	parts := strings.Split(image, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid image format, expected name:tag")
	}
	imageName := parts[0]
	imageTag := parts[1]

	fmt.Println("DEBUG: imageName =", imageName)
    fmt.Println("DEBUG: imageTag  =", imageTag)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	imagesRaw, ok := obj["images"]
	if !ok {
		return nil, fmt.Errorf("no images section found in %s", path)
	}

	images := imagesRaw.([]interface{})
	found := false

	for _, img := range images {
		m := img.(map[string]interface{})
		if m["name"] == imageName {
			m["newTag"] = imageTag
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("image %s not found in %s", imageName, path)
	}

	out, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}

	if !dryRun {
		if err := os.WriteFile(path, out, 0644); err != nil {
			return nil, err
		}
	}

	return out, nil
}
