# deps check

This repo contains a tool to validate imports inside Go monorepos,
and to run tests selectively based on changed files and imports structure.

## Installation

```bash
go install github.com/moveaxlab/deps-check
```

## Configuration

This tool distingueshes between 4 layers:

- `external` are for vendored external dependencies. We ignore them, anyone can import them.
- `common` are packages that contain shared code between services.
  Common packages can import from each other, and services can import from common packages.
- `service` are packages containing code for a component or microservice.
  They can import from common packages and from external packages,
  and cannot import code from other services.
- `utility` are packages that contain utility scripts that should only be run locally.
  They can import from wherever they want, and no one can import them.

## Validating project structure

```bash
deps-check validate
```

## Running tests selectively

Get a list of changed files and feed it to `deps-check`:

```bash
readarray -t PACKAGES < <(git diff --name-only "${TARGET_BRANCH}...${CURRENT_BRANCH}" | deps-check changed-packages)
```

You can also run tests on staged files:

```bash
readarray -t PACKAGES < <(git diff --staged --name-only | deps-check changed-packages)
```

Then, iterate over the changed packages and run your tests:

```bash
for PACKAGE in "${PACKAGES[@]}"
do
  go test  "$PACKAGE"
done
```
