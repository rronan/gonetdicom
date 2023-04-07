# Gonetdicom

Gonetdicom is an open-source software library written in Go, designed to provide an efficient implementation of tools for DICOM communication. It offers full DICOMweb functionalities, including WADO, STOW, and HTTP PUT support for DICOM instances. Gonetdicom uses the DICOM dataset implementation from [suyashkumar/dicom](https://pkg.go.dev/github.com/suyashkumar/dicom).

## Features

- WADO support (retrieving DICOM instances)
- STOW support (storing DICOM instances)
- HTTP PUT support (sending DICOM instances)

- In-progress: implementation of C-STORE protocol (STORESCP, STORESCU)

## Installation

To install Gonetdicom, run:

```
go get github.com/rronan/gonetdicom
```

## Usage

```go
import (
	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/rronan/gonetdicom/dicomweb"
)
```

Here is a brief overview of the available packages and their functionalities:

### dicomweb

This package contains the implementation of DICOMweb functionalities:

- `dicomweb/get.go`: Functions for retrieving DICOM instances via WADO.
  - `Get()`: Sends a GET request to a specified URL with headers and returns an HTTP response.
  - `ReadMultipart()`: Reads DICOM datasets from a multipart HTTP response.
  - `Wado()`: Retrieves DICOM datasets from a specified URL using the WADO-RS standard.
  - `ReadMultipartToFile()`: Reads DICOM datasets from a multipart HTTP response and saves them to a specified folder.
  - `WadoToFile()`: Retrieves DICOM datasets from a specified URL using the WADO-RS standard and saves them to a specified folder.

- `dicomweb/post.go`: Functions for storing DICOM instances via STOW.
  - `WriteMultipart()`: Writes DICOM datasets to a multipart message.
  - `PostMultipart()`: Sends a POST request to a specified URL with a multipart message and headers, and returns an HTTP response.
  - `Stow()`: Stores DICOM datasets to a specified URL using the STOW-RS standard.
  - `WriteMultipartFromFile()`: Writes DICOM datasets from files to a multipart message.
  - `StowFromFile()`: Stores DICOM datasets from files to a specified URL using the STOW-RS standard.

- `dicomweb/put.go`: Functions for sending DICOM instances via HTTP PUT.
  - `Put()`: Sends a DICOM dataset to a specified URL using HTTP PUT with headers.
  - `PutFromFile()`: Sends a DICOM dataset from a file to a specified URL using HTTP PUT with headers.

### dicomutil

This package provides utility functions for working with DICOM datasets:

- `dicomutil/dicomutil.go`: Various utility functions, including trimming tags, getting UIDs, parsing file UIDs, and converting between DICOM datasets and byte arrays.

### dicomweb

Work in progress, not merge yet.

## Tests

To run the tests:

```
go test -v ./...
```

## Contributing

We welcome and encourage collaboration from the community. If you're interested in contributing to Gonetdicom, don't hesitate to fork the repository and propose a PR.