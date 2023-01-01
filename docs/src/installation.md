# Installation

There are multiple ways to install *Track*:

[[_TOC_]]

## Precompiled binaries

1. Download the [latest binaries](https://github.com/mlange-42/track/releases) for your platform  
   (Binaries are available for Linux, Windows and macOS)
2. Unzip somewhere
3. *Optional:* add the parent directory of the executable to your `PATH` environmental variable

## From GitHub using Go

In case you have [Go](https://go.dev/) installed, you can install *Track* with `go install`:

```
go install github.com/mlange-42/track@latest
```

## Clone and build

To build *Track* locally, e.g. to contribute to the project, you need to clone the repository on your local machine:

```
git clone https://github.com/mlange-42/track
```

`cd` into `track/` and run

```
go build
```

The resulting binaries can be found in the `track` root directory under the name `track` or `track.exe`.
