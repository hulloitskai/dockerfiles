# dockerfiles

_Dockerfiles for personal utility images._

[![Drone: Status][drone-img]][drone]

> Mo' Dockerfiles, mo' problems.

## Catalogue

- [`frp`](./frp) – A fast reverse proxy to help you expose a local server behind
  a NAT or firewall to the internet.
  [Sourced from `fatedier/frp`](https://github.com/fatedier/frp).

- [`kaniko-drone`](./kaniko-drone) – A [Drone](https://drone.io) plugin that
  runs Docker image builds on
  [Kaniko](https://github.com/GoogleContainerTools/kaniko).

- [`golinter`](./golinter) – A linter for Go source code that combines
  [`goimports`](./golang.org/x/tools/cmd/goimports),
  [`revive`](https://github.com/mgechev/revive), and
  [`go vet`](https://golang.org/cmd/vet/).

- [`grpc-go`](./grpc-go) – A builder Docker image containing Go and GRPC /
  protobuf tooling.

[drone]: https://ci.stevenxie.me/stevenxie/dockerfiles
[drone-img]: https://ci.stevenxie.me/api/badges/stevenxie/dockerfiles/status.svg
