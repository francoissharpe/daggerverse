package main

import "fmt"

func (p *Python) WithPypaBuild(
// The version of the build package
// +optional
// +default="1.2.1"
	buildVersion string,
// Directory containing the source code
	src *Directory,
// Version of the package to be built
//version string,

) *Directory {
	outDir := "/output"
	p.Ctr = p.Ctr.Pipeline("python-pypa-build").
		WithExec([]string{"pip", "install", fmt.Sprintf("build==%s", buildVersion)}).
		WithDirectory(workdir, src).
		WithExec([]string{"python", "-m", "build", "--outdir", outDir, workdir})
	return p.Ctr.Directory(outDir)
}
