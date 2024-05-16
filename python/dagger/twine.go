package main

import (
	"fmt"
	"time"
)

func (p *Python) WithTwineUpload(
// The version of the build package
// +optional
// +default="5.0.0"
	twineVersion string,
// Directory containing the distribution files to upload to the repository
	dist *Directory,
// The repository (package index) URL to upload the package to
// +optional
// +default="https://test.pypi.org/legacy/"
	repositoryUrl string,
// The username to authenticate with the repository
	username *Secret,
// The password to authenticate with the repository
	password *Secret,
// Force execution of the step even if the cache is present
// +optional
// +default=false
	bustCache bool,
) *Python {
	outDir := "/dist"
	p.Ctr = p.Ctr.Pipeline("python-twine-upload").
		WithExec([]string{"pip", "install", fmt.Sprintf("twine==%s", twineVersion)}).
		WithSecretVariable("TWINE_USERNAME", username).
		WithSecretVariable("TWINE_PASSWORD", password).
		WithDirectory(workdir+outDir, dist)

	if bustCache {
		now := time.Now().Format("YYYYMMDDHHMMSS")
		p.Ctr = p.Ctr.WithEnvVariable("CACHEBUSTER", now)
	}
	p.Ctr = p.Ctr.
		WithExec([]string{
			"python",
			"-m",
			"twine",
			"upload",
			"--non-interactive",
			"--disable-progress-bar",
			"--comment",
			"Uploaded via Dagger",
			"--skip-existing",
			"--repository-url",
			repositoryUrl,
			workdir + outDir + "/*",
		})
	return p
}
