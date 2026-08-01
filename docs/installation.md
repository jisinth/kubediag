# Installation

## Prerequisites

- Go 1.22+
- A kubeconfig with access to the cluster you want to diagnose

## Build from source

```
git clone https://github.com/jisinth/kubediag.git
cd kubediag
go mod tidy
go build -o kubediag .
```

This produces a `kubediag` binary in the repo root. Move it onto your `PATH`, e.g.:

```
mv kubediag /usr/local/bin/kubediag
```

## Verify

```
kubediag version
```

## Releases

Prebuilt binaries via GitHub Actions/GitHub Releases are planned — see
[docs/roadmap.md](roadmap.md). Until then, build from source as above.
