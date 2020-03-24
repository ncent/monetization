# Monetization

This module will be responsible for all the Monetization related events on the platform, the events are managed by the EventBridge service on AWS

## Getting Started

### Prerequisites

```
GO
Node
npm
Serverless Framework
ginkgo
```

### Installing

Install Serverless Framework

```
npm install -g serverless
```

Install node dependencies

```
npm install
```

## Building

```
make build
```

## Running the tests

```
go test ./...
```

or

```
ginkgo ./...
```

## Deployment

Development environment
```
SLS_DEBUG=* serverless deploy --verbose --stage development
```

Production environment
```
SLS_DEBUG=* serverless deploy --verbose --stage production
```

## Built With

* [GO](https://golang.org) - The Language
* [Serverless Framework](https://serverless.com) - Deployment Framework

## Authors

* **Eduardo Nunes Peireira** - *Initial work* - [eduardonunesp](https://gitlab.com/eduardonunesp)