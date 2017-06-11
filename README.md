
##### Current version: UNRELEASED / FIRST VERSION IN DEV

[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/millisecond/olb/master/LICENSE)

OLB is a fast, modern, load balancer (HTTP(S) / TCP / UDP) for deploying applications on AWS.

OLB is developed and maintained by [Casey Haakenson](https://twitter.com/millisecond).

It supports ([Full feature list](https://github.com/millisecond/olb/wiki/Features))

* [TLS termination with dynamic certificate stores](https://github.com/millisecond/olb/wiki/Features#certificate-stores)
* [Raw TCP proxy](https://github.com/millisecond/olb/wiki/Features#tcp-proxy-support)
* [TCP+SNI proxy for full end-to-end TLS](https://github.com/millisecond/olb/wiki/Features#tcpsni-proxy-support) without decryption
* [HTTPS upstream support](https://github.com/millisecond/olb/wiki/Features#https-upstream-support)
* [Websockets](https://github.com/millisecond/olb/wiki/Features#websocket-support) and
  [SSE](https://github.com/millisecond/olb/wiki/Features#sse---server-sent-events)
* [Dynamic reloading without restart](https://github.com/millisecond/olb/wiki/Features#dynamic-reloading)
* [Traffic shaping](https://github.com/millisecond/olb/wiki/Features#traffic-shaping) for "blue/green" deployments,
* [Circonus](https://github.com/millisecond/olb/wiki/Features#metrics-support),
  [Graphite](https://github.com/millisecond/olb/wiki/Features#metrics-support) and
  [StatsD/DataDog](https://github.com/millisecond/olb/wiki/Features#metrics-support) metrics
* [WebUI](https://github.com/millisecond/olb/wiki/Features#web-ui)

The full documentation is on the [Wiki](https://github.com/millisecond/olb/wiki).

## Getting started

1. Install from source, [binary](https://github.com/millisecond/olb/releases),
   [Docker](https://hub.docker.com/r/millisecond/olb/) or [Homebrew](http://brew.sh).
    ```
	# go 1.8 or higher is required
    go get github.com/millisecond/olb                     (>= go1.8)

    brew install olb                                  (OSX/macOS stable)
    brew install --devel olb                          (OSX/macOS devel)

    docker pull millisecond/olb                           (Docker)

    https://github.com/millisecond/olb/releases           (pre-built binaries)
    ```

2. Create a Target Group

2a. Route 53

3. Create an Auto Scaling Group

* User properties

* IAM Role

* Public/private routes

4. 

5. 

6. 

7. Done

## Maintainers

* Casey Haakenson [@millisecond](https://twitter.com/millisecond)

## License

See [LICENSE](https://github.com/millisecond/olb/blob/master/LICENSE) for details.
