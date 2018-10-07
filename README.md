# Chrome-proxy

Chrome-proxy is a tcp proxy to expose the
[Chrome headless remote debug protocol](https://chromedevtools.github.io/devtools-protocol/).

## Usage

The command line accepts some options:
- `--help` shows the option available
- `--bind` is the chrome's address to request, default `127.0.0.1:9222`
- `--key` is a secret need to request the proxy via `Api-Key` http header, default `secret`
- `--listen` is the proxy's server address, default `127.0.0.1:8080`

## Build

```
$ go build
```

## Quick start

The repository includes a docker container with a chrome headless.
You can use it to try easily the proxy.

These commands will build the container and start a chrome headless. It will
expose the `9222` port into the host.

```
$ make docker-build-chrome
$ make docker-run-chrome
```

Then you can run the proxy
```
$ ./chrome-proxy
```

And test via curl
```
$ curl -H 'Api-Key: secret'  http://127.0.0.1:8080/json/version
{
   "Browser": "HeadlessChrome/69.0.3497.100",
   "Protocol-Version": "1.3",
   "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.100 Safari/537.36",
   "V8-Version": "6.9.427.23",
   "WebKit-Version": "537.36 (@8920e690dd011895672947112477d10d5c8afb09)",
   "webSocketDebuggerUrl": "ws://127.0.0.1:8080/devtools/browser/bed69809-6abd-4ecc-b8f8-af9bbcd29157"
}
```

## Unit test

```
$ go test
```
