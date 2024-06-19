# lac - Lightweight API Client

`lac` is a cli only API client. It can be used for API exploration and testing.

![screenshot of lac](docs/images/carbon.png)

## Installation

### Install with Go

```sh
go install github.com/eynopv/lac@latest
```

## Usage

```sh
lac [flags] request.json
```

### Example

**Create request object**

```js
// get.json

{
  "path": "https://httpbin.org/get"
}
```

**Send the request**

```sh
lac get.json
```

### Guide

Read [user guide](docs/UserGuide.md)

## License

This project is licensed under the BSD 3-Clause License.
