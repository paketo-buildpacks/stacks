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
- Tag: **`some-registry/build:1.1.0-base`**
- Digest: `some-registry/build@sha256:some-base-sha`
## Build Receipt Diff
```
-ii  ruby-loofah  1.6.10ubuntu0.1  amd64  some description
+ii  ruby-loofah  1.6.12ubuntu0.1  amd64  some description
```
## Run Receipt Diff
```
-ii  ruby-loofah  1.6.10ubuntu0.1  amd64  some description
+ii  ruby-loofah  1.6.12ubuntu0.1  amd64  some description
```
