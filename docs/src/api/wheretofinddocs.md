# Where to find API docs
The API documentation, depends on which wrapper you are using. If you're writing a wrapper yourself, or want some background information on how it works internally, you can find function docs in the [exports/exports.go](https://github.com/eduvpn/eduvpn-common/blob/v2/exports/exports.go) file.

This file is commented using Go comment style. It gives a basic of what the function does, what it returns and what type of arguments you should pass to it. The API documentation for the Python wrapper can be [found here](https://eduvpn.github.io/eduvpn-common/api/python/rtd/index.html).

There is also a Go API that is defined in the [client package](https://github.com/eduvpn/eduvpn-common/tree/v2/client). However, this is not the primary use case for the library.