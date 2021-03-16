## Release Images

### Runtime Base Images

#### For CNB Builds:
- Tag: **`some-registry/run:1.0-tiny-cnb`**
- Digest: `some-registry/run@sha256:some-cnb-sha`

#### Source Image (e.g., for Dockerfile builds):
- Tag: **`some-registry/run:1.0-tiny`**
- Digest: `some-registry/run@sha256:some-base-sha`

### Build-time Base Images

#### For CNB Builds:
- Tag: **`some-registry/build:1.0-tiny-cnb`**
- Digest: `some-registry/build@sha256:some-cnb-sha`

#### Source Image (e.g., for Dockerfile builds):
- Tag: **`some-registry/build:1.0-tiny`**
- Digest: `some-registry/build@sha256:some-base-sha`
## Patched USNs
[USN-4498-1](https://ubuntu.com/security/notices/USN-4498-1):  Loofah vulnerability
* [CVE-2019-15587](https://people.canonical.com/~ubuntu-security/cve/CVE-2019-15587): In the Loofah gem for Ruby through v2.3.0 unsanitized JavaScript may occur in sanitized output when a crafted SVG element is republished.

## Build Receipt Diff
```
-ii  ruby-loofah  1.6.10ubuntu0.1  amd64  some-description
+ii  ruby-loofah  1.6.12ubuntu0.1  amd64  some-description
```
