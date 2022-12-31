# Installation

There are multiple ways to install Yarner:

[[_TOC_]]

## Binaries

1. Download the [latest binaries](https://github.com/mlange-42/track/releases) for your platform  
   (Binaries are available for Linux, Windows and macOS)
2. Unzip somewhere
3. *Optional:* add the parent directory of the executable to your `PATH` environmental variable

## From GitHub using Go

In case you have [Go](https://go.dev/) installed, you can install with `go install`:

```
go install github.com/mlange-42/track@latest
```

## Clone and build

To build track locally, e.g. to contribute to the project, you will have to clone the repository on your local machine:

```
git clone https://github.com/mlange-42/track
```

`cd` into `track/` and run

```
go build
```

The resulting binary can be found in the `track` directory under the name `track` or `track.exe`.
