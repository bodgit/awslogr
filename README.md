[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/bodgit/awslogr/badge)](https://securityscorecards.dev/viewer/?uri=github.com/bodgit/awslogr)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10586/badge)](https://www.bestpractices.dev/projects/10586)
[![GitHub release](https://img.shields.io/github/v/release/bodgit/awslogr)](https://github.com/bodgit/awslogr/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/bodgit/awslogr/build.yml?branch=main)](https://github.com/bodgit/awslogr/actions?query=workflow%3ABuild)
[![Coverage Status](https://coveralls.io/repos/github/bodgit/awslogr/badge.svg?branch=main)](https://coveralls.io/github/bodgit/awslogr?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/bodgit/awslogr)](https://goreportcard.com/report/github.com/bodgit/awslogr)
[![GoDoc](https://godoc.org/github.com/bodgit/awslogr?status.svg)](https://godoc.org/github.com/bodgit/awslogr)
![Go version](https://img.shields.io/badge/Go-1.24-brightgreen.svg)
![Go version](https://img.shields.io/badge/Go-1.23-brightgreen.svg)

# awslogr

This simple package provides a means to use [github.com/go-logr/logr](https://pkg.go.dev/github.com/go-logr/logr) with the [github.com/aws/aws-sdk-go-v2](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2) AWS SDK.

This can be useful when adding OpenTelemetry tracing and logging, using the [otellogr](https://go.opentelemetry.io/contrib/bridges/otellogr) and [otelaws](http://go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws) packages.

An OpenTelemetry example:

```golang
l := logr.New(otellogr.NewLogSink("my/pkg/name"))

nl, err := awslogr.New(l)
if err != nil {
	log.Fatal(err)
}

cfg, err := config.LoadDefaultConfig(context.Background(),
	config.WithClientLogMode(aws.LogRetries|aws.logRequest),
	config.WithLogger(nl))
if err != nil {
	log.Fatal(err)
}

otelaws.AppendMiddlewares(&cfg.APIOptions)

// Use cfg to now create an AWS service client
```
