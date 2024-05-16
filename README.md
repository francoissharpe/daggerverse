# Daggerverse

## Modules
| Module | Description |
|--------|-------------|
| [python](./python) | Tooling for building and publishing Python projects |
## Getting started

### Install the dagger cli

```shell
curl -L https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh
```

### Run a container with Python 3.11 and poetry
```shell
dagger -m github.com/francoissharpe/daggerverse/python@main call \
  with-version --version 3.11 \
  with-package-manager --package-manager poetry \
  container terminal --cmd /bin/bash
```

### Build a python project

```shell
dagger call with-version --version 3.11 \
  with-pypa-build --src "https://github.com/francoissharpe/pekish.git#master" \
  directory --path ./ entries
```

## References
[Installing the Dagger CLI](https://docs.dagger.io/quickstart/729237/cli/)