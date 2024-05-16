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