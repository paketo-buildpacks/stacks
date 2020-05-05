# Metadata for Paketo Build/Run Stack Images

The following build/run stack images are supported:  
**Base**: 
	- `gcr.io/paketo-buildpacks/[build/run]:base-cnb` 
	- `paketobuildpacks/[build/run]:base-cnb`
**Full**: 
	- `gcr.io/paketo-buildpacks/[build/run]:full-cnb-cf`  
	- `paketobuildpacks/[build/run]:full-cnb-cf`
**Tiny**: 
	- `gcr.io/paketo-buildpacks/[build/run]:tiny-cnb`
	- `paketobuildpacks/[build/run]:tiny-cnb`

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

Contains:
- Build: ubuntu:bionic + openssl + CA certs + compilers + shell utilities
- Run: distroless-like bionic + glibc + openssl + CA certs
