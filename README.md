# dep check

This repo contains a tool to validate imports inside Go monorepos,
and to run tests selectively based on changed files and imports structure.

## Installation

```bash
go install github.com/moveaxlab/dep-check
```

## Configuration

This tool distingueshes between 4 types of packages:

- `external` is used for [vendored](https://go.dev/ref/mod#vendoring) external dependencies
  and other stuff we don't care about. We ignore them, anyone can import them.
- `common` are packages that contain shared code between services.
  Common packages can import from each other, and services can import from common packages.
- `service` are packages containing code for a component or microservice.
  They can import from common packages and from external packages,
  and cannot import code from other services.
- `utility` are packages that contain utility scripts that should only be run locally.
  They can import from wherever they want, and no one can import them.

To configure the tool, create a `dep-check.yaml` file in the root of your Go project,
and fill its contents like this:

```yaml
module_name: github.com/moveaxlab/dep-check  # the module name of your project
folders:
  external:
    - external/*
  common:
    - common/src/*
    - models/*
  utility:
    - scripts
    - common/scripts
  service:
    - service/*
```

The `folders` configuration specifies which folders fall in which category.
All folders are specified from the root of your Go project.
You can use the `*` symbols as a wildcard: this instructs `dep-check` to treat
each folder under the specified path as a separate package.

An example is worth a thousand words.
Given the above configuration and the following directory structure:

```
common/
  scripts/
    script_1/
    script_2/
  src/
    config/
    postgres/
    kafka/
      consumer/
      producer/
external/
  dep_1/
  dep_2/
models/
  auth/
  common/
  order/
  product/
scripts/
  script_3/
  script_4/
service/
  service_1/
  service_2/
  service_3/
other_package/
```

`dep-check` will identify these packages:

- external packages:
  - `external`
  - `other_package`
- utility packages:
  - `common/scripts`
  - `scripts`
- common packages:
  - `common/src/config`
  - `common/src/postgres`
  - `common/src/kafka`
  - `models/auth`
  - `models/common`
  - `models/order`
  - `models/product`
- service packages:
  - `service/service_1`
  - `service/service_2`
  - `service/service_3`

All packages that are not matched by an explicit rule will fall under the external category.
You can add the `'*'` catch all to the service category to treat all other packages as services.

If the root of your repo is not the same as the root of your Go project,
you can add the `root_dir` configuration option to specify the path from the root of the repo
to the root of your Go project.

> You only need this to run tests selectively when using git.

You can now run `dep-check` from the root of your Go project.

## Validating project structure

```bash
dep-check validate
```

## Running tests selectively

Get a list of changed files and feed it to `dep-check`:

```bash
readarray -t PACKAGES < <(git diff --name-only "${TARGET_BRANCH}...${CURRENT_BRANCH}" | dep-check changed-packages)
```

You can also run tests on staged files:

```bash
readarray -t PACKAGES < <(git diff --staged --name-only | dep-check changed-packages)
```

Then, iterate over the changed packages and run your tests:

```bash
for PACKAGE in "${PACKAGES[@]}"
do
  go test  "$PACKAGE"
done
```
