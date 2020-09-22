# Metadata for Paketo Build/Run Stack Images

The following build/run stack images are supported:\
**Base**:\
Build Images
- `index.docker.io/paketobuildpacks/build:base-cnb`  

Run Images
- `index.docker.io/paketobuildpacks/run:base-cnb` 

**Full**:\
Build Images
- `index.docker.io/paketobuildpacks/build:full-cnb`

Run Images
- `index.docker.io/paketobuildpacks/run:full-cnb`

**Tiny**:\
Build Images
- `index.docker.io/paketobuildpacks/build:tiny-cnb`

Run Images
- `index.docker.io/paketobuildpacks/run:tiny-cnb`

## Use Cases


### Base (aka "bionic")
Ideal for:
- Java apps and .NET Core apps
- Go apps that require some C libraries
- Node.js/Python/Ruby/etc. apps **without** many native extensions

Contains:
- Build: ubuntu:bionic + openssl + CA certs + compilers + shell utilities
- Run: ubuntu:bionic + openssl + CA certs

### Full
Ideal for:
- PHP/Node.js/Python/Ruby/etc. apps **with** many native extensions

Contains:
- Build: ubuntu:bionic + many common C libraries and utilities
- Run: ubuntu:bionic + many common libraries and utilities

### Tiny
Ideal for:
- Most Go apps
- Java [GraalVM Native Images](https://www.graalvm.org/docs/reference-manual/native-image/)

Contains:
- Build: ubuntu:bionic + openssl + CA certs + compilers + shell utilities
- Run: distroless-like bionic + glibc + openssl + CA certs
