package main

import "fmt"

func (d *Dotnet) WithVersion(
	// The image to use
	// +optional
	// +default="mcr.microsoft.com/dotnet/sdk"
	image string,
	// The version of the .NET SDK to use
	// +optional
	// +default="8.0"
	version string,
	// Cache enabled
	// +optional
	// +default=true
	cacheEnabled bool,
) *Dotnet {
	baseImage := fmt.Sprintf("%s:%s", image, version)
	d.setupBaseImage(baseImage)
	d.setupCache(cacheEnabled)
	return d
}

func (d *Dotnet) setupBaseImage(baseImage string) {
	p.Ctr = dag.Pipeline("dotnet-base").
		Container().
		From(fmt.Sprintf("mcr.microsoft.com/dotnet/sdk:%s", baseImage)).
		WithEnvVariable("DOTNET_CLI_TELEMETRY_OPTOUT", "1").
		WithEnvVariable("DOTNET_NOLOGO", "1").
		WithEnvVariable("DOTNET_MULTILEVEL_LOOKUP", "0")
}

func (d *Dotnet) setupCache(cacheEnabled bool) {
	if cacheEnabled {
		httpCachePath := "/root/.local/share/NuGet/v3-cache"
		nugetScratchPath := "/root/.local/share/NuGet/Scratch"
		nugetPluginCachePath := "/root/.local/share/NuGet/plugins-cache"

		d.Ctr = d.Ctr.WithMountedCache(globalPackagesCachePath, dag.CacheVolume("nuget-packages-cache")).
			WithMountedCache(httpCachePath, dag.CacheVolume("nuget-http-cache")).
			WithMountedCache(nugetScratchPath, dag.CacheVolume("nuget-scratch")).
			WithMountedCache(nugetPluginCachePath, dag.CacheVolume("nuget-plugin-cache"))
	}
}

func (d *Dotnet) Restore(
	// The project path to restore
	// +required
	project string,
) *Dotnet {
	d.Ctr = d.Ctr.WithExec([]string{"dotnet", "restore", "--packages", globalPackagesCachePath, project})
	return d
}
