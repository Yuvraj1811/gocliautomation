package gitops 

import "fmt"

func KustomizationPath(app, env string) string {
	return fmt.Sprintf(
		"apps/%s/overlays/%s/kustomization.yaml",
		app,
		env,
	)
}