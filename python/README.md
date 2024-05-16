## Start a static site built with mkdocs-material

```shell
dagger call with-mkdocs-material --src https://github.com/squidfunk/mkdocs-material#master with-static-site-container as-service up --ports 8080:80
```

## Push the static site image to a registry

```shell
dagger call with-mkdocs-material --src https://github.com/squidfunk/mkdocs-material#master with-static-site-container publish --address ttl.sh/$(uuidgen):1h
```