package main

import (
	"fmt"
	"strings"
)

func (p *Python) WithVersion(
	// The image to use
	// +optional
	// +default="python"
	image string,
	// The version of Python to use
	version string,
	// Additional APT packages to install
	// +optional
	// +default=[]
	packages []string,
	// Enable caching for the container
	// +optional
	// +default=true
	cacheEnabled bool,
	// Commands to run in the container
	// +optional
	// +default=[]
	commands []string,
	// CA File to use for SSL
	// +optional
	// +default=null
	caFile *File,
	// Pip index URL
	// +optional
	// +default="https://pypi.org/simple"
	pipIndexURL string,
	// http_proxy environment variable
	// +optional
	// +default=""
	httpProxy string,
	// https_proxy environment variable
	// +optional
	// +default=""
	httpsProxy string,
	// no_proxy environment variable
	// +optional
	// +default=""
	noProxy string,
) *Python {
	baseImage := fmt.Sprintf("%s:%s", image, version)
	p.setupBaseImage(baseImage, pipIndexURL)
	p.setupCAFile(caFile)
	p.setupCaching(cacheEnabled)
	p.setupProxyEnvVariables(httpProxy, httpsProxy, noProxy)
	p.installPackages(packages)
	p.runCommands(commands)
	p.setupWorkdir()
	p.BaseImageRef = baseImage
	return p
}

func (p *Python) setupBaseImage(baseImage string, pipIndexURL string) {
	p.Ctr = dag.Container().
		From(baseImage).
		Pipeline("python-base").
		WithEnvVariable("DEBIAN_FRONTEND", "noninteractive").
		WithEnvVariable("PYTHONUNBUFFERED", "1").
		WithEnvVariable("PYTHONFAULTHANDLER", "1").
		WithEnvVariable("PIP_INDEX_URL", pipIndexURL).
		WithEnvVariable("PIP_ROOT_USER_ACTION", "ignore").
		WithEnvVariable("PIP_DISABLE_PIP_VERSION_CHECK", "on").
		WithEnvVariable("PIP_DEFAULT_TIMEOUT", "100")
}

func (p *Python) setupCAFile(caFile *File) {
	if caFile != nil {
		p.Ctr = p.Ctr.
			WithEnvVariable("REQUESTS_CA_BUNDLE", caBundlePath).
			WithEnvVariable("GRPC_DEFAULT_SSL_ROOTS_FILE_PATH", caBundlePath).
			WithEnvVariable("CURL_CA_BUNDLE", caBundlePath).
			WithEnvVariable("SSL_CERT_FILE", caBundlePath).
			WithMountedFile(caBundlePath, caFile)
	}
}

func (p *Python) setupProxyEnvVariables(
	httpProxy string,
	httpsProxy string,
	noProxy string,
) {
	if httpProxy != "" {
		p.Ctr = p.Ctr.WithEnvVariable("http_proxy", httpProxy)
	}
	if httpsProxy != "" {
		p.Ctr = p.Ctr.WithEnvVariable("https_proxy", httpsProxy)
	}
	if noProxy != "" {
		p.Ctr = p.Ctr.WithEnvVariable("no_proxy", noProxy)
	}
}

func (p *Python) setupCaching(cacheEnabled bool) {
	if cacheEnabled {
		aptCacheDir := cacheRootDir + "/apt"
		pipCacheDir := cacheRootDir + "/pip"
		poetryCacheDir := cacheRootDir + "/pypoetry"

		aptCache := dag.CacheVolume("apt-cache")
		pipCache := dag.CacheVolume("pip-cache")
		poetryCache := dag.CacheVolume("poetry-cache")

		p.Ctr = p.Ctr.
			WithEnvVariable("PIP_CACHE_DIR", pipCacheDir).
			WithEnvVariable("POETRY_CACHE_DIR", poetryCacheDir).
			WithMountedCache(aptCacheDir, aptCache).
			WithMountedCache(pipCacheDir, pipCache).
			WithMountedCache(poetryCacheDir, poetryCache)
	}
}

func (p *Python) installPackages(packages []string) {
	if len(packages) > 0 {
		p.Ctr = p.Ctr.WithExec(
			ShDashC([]string{
				"apt-get update",
				"apt-get install -yqq --no-install-recommends " + strings.Join(packages, " "),
				"rm -rf /var/lib/apt/lists/*"}))
	}
}

func (p *Python) runCommands(commands []string) {
	if len(commands) > 0 {
		p.Ctr = p.Ctr.WithExec(ShDashC(commands))
	}
}

func (p *Python) setupWorkdir() {
	p.Ctr = p.Ctr.WithWorkdir(workdir)
}

//func (p *Python) WithProductionImage(
//	directory *Directory,
//) *Container {
//	aptCacheDir := cacheRootDir + "/apt"
//	poetryCacheDir := cacheRootDir + "/pypoetry"
//	ctr := dag.Container().
//		From("debian:12-slim").
//		WithEnvVariable("POETRY_HOME", "/usr/local").
//		WithEnvVariable("POETRY_NO_INTERACTION", "1").
//		WithEnvVariable("POETRY_CACHE_DIR", poetryCacheDir).
//		WithMountedCache(aptCacheDir, dag.CacheVolume("apt")).
//		WithMountedCache(poetryCacheDir, dag.CacheVolume("poetry")).
//		WithEnvVariable("PATH", "/root/.local/bin:$PATH", ContainerWithEnvVariableOpts{Expand: true}).
//		WithExec([]string{"sh", "-c", "apt-get update && apt-get install --no-install-suggests --no-install-recommends --yes pipx"}).
//		WithExec([]string{"sh", "-c", "pipx install poetry && pipx inject poetry poetry-plugin-bundle"}).
//		WithWorkdir(workdir).
//		WithDirectory(workdir, directory).
//		WithExec([]string{"poetry", "bundle", "venv", "--python=/usr/bin/python3", "--only=main", "/venv"})
//
//	prodCtr := dag.Container().
//		From("gcr.io/distroless/python3-debian12").
//		WithDirectory("/venv", ctr.Directory("/venv")).
//		WithEntrypoint([]string{"/venv/bin/python"})
//
//	return prodCtr
//}

func (p *Python) WithPackageManager(
	// The package manager to use
	// +optional
	// +default="pip"
	packageManager string,
) *Python {
	switch packageManager {
	case "poetry":
		return p.WithPoetry("", []string{})
	default:
		return p.WithPip("")
	}
}

func (p *Python) WithPip(
	// The version of pip to use
	// +optional
	// +default=""
	version string,
) *Python {
	p.Ctr = p.Ctr.Pipeline("python-pip")
	if version != "" {
		p.Ctr = p.Ctr.WithExec([]string{"pip", "install", "--upgrade", "pip==" + version})
	}
	return p
}

func (p *Python) WithPoetry(
	// The version of poetry to use
	// +optional
	// +default="1.6.1"
	version string,
	// Additional poetry plugins to install
	// +optional
	// +default=[]
	plugins []string,
) *Python {
	p.Ctr = p.Ctr.Pipeline("python-poetry")
	p.Ctr = p.Ctr.
		WithEnvVariable("POETRY_HOME", "/usr/local").
		WithEnvVariable("POETRY_VIRTUALENVS_CREATE", "true").
		WithEnvVariable("POETRY_INSTALLER_MAX_WORKERS", "20").
		WithEnvVariable("POETRY_NO_INTERACTION", "1")
	if version != "" {
		p.Ctr = p.Ctr.WithEnvVariable("POETRY_VERSION", version)
	}
	// Review installing with pipx
	installCmd := "curl -sSL https://install.python-poetry.org | python -"
	p.Ctr = p.Ctr.WithExec(ShDashC([]string{installCmd}))
	if len(plugins) > 0 {
		for _, plugin := range plugins {
			p.Ctr = p.Ctr.WithExec(ShDashC([]string{fmt.Sprintf("pip install %s", plugin)}))
		}
	}
	return p
}
