// A generated module for Mkdocs functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import "fmt"

type Mkdocs struct {
	Ctr  *Container
	Site *Directory
}

// Returns a container that echoes whatever string argument is provided
func (m *Mkdocs) WithMkdocsMaterial(
	// The image to use
	// +optional
	// +default="squidfunk/mkdocs-material"
	image string,
	// The version of the image to use
	// +optional
	// +default="latest"
	version string,
	// Your mkdocs.yml file
	// +optional
	// +default=null
	mkdocsYaml *File,
	// Your docs directory
	// +optional
	// +default=null
	docs *Directory,
	// Source directory
	// +optional
	// +default=null
	src *Directory,
) (*Mkdocs, error) {
	m.Ctr = dag.Container().From(fmt.Sprintf("%s:%s", image, version))
	if src != nil {
		m.Ctr = m.Ctr.WithDirectory("/docs", src)
	} else if docs != nil && mkdocsYaml != nil {
		m.Ctr = m.Ctr.
			WithDirectory("/docs/", docs).
			WithFile("/docs/mkdocs.yml", mkdocsYaml)
	} else {
		return nil, fmt.Errorf("either --src or both --docs and --mkdocs-aml must be provided")
	}
	m.Ctr = m.Ctr.WithExec([]string{"build"})
	m.Site = m.Ctr.Directory("/docs/site")
	return m, nil
}

func (m *Mkdocs) WithStaticSiteContainer(
	// The image to use
	// +optional
	// +default="nginx"
	image string,
	// The image version
	// +optional
	// +default="alpine"
	version string,
) *Container {
	return dag.Container().
		Pipeline("mkdocs-nginx-build").
		From(fmt.Sprintf("%s:%s", image, version)).
		WithDirectory("/usr/share/nginx/html", m.Site)
}
