# pprof-upload

[![Build Status](https://travis-ci.com/rakyll/pprof-upload.svg?token=Quf3mWAszVwsMXDMWxkm&branch=master)](https://travis-ci.com/rakyll/pprof-upload)

Uploads pprof files to [Stackdriver Profiler](https://cloud.google.com/profiler/).

## Requirements

* Enable [Stackdriver Profiler API](https://console.cloud.google.com/apis/library/cloudprofiler.googleapis.com).
* If running outside of Google Compute Engine, install [gcloud](https://cloud.google.com/sdk/gcloud/) and run `gcloud auth application-default login`.

## Installation

Linux 64-bit:

```
$ curl http://storage.googleapis.com/jbd-releases/pprof-upload-linuxamd64 > pprof-upload && chmod +x pprof-upload
```

macOS 64-bit:

```
$ curl http://storage.googleapis.com/jbd-releases/pprof-upload-darwinamd64 > pprof-upload && chmod +x pprof-upload
```

Windows 64-bit:

* Download http://storage.googleapis.com/jbd-releases/pprof-upload-windowsamd64
* Put the download folder to your PATH.

## Usage

Capture pprof profiles, for example by using the
[net/http/pprof](https://golang.org/pkg/net/http/pprof) package. See `examples/helloworld` for an example.

```
$ curl http://localhost:6060/debug/pprof/profile?seconds=30 > pprof.out
$ pprof-upload -target=webserver
https://console.cloud.google.com/profiler/webserver;type=CPU?project=PROJECT
```

![Cloud Profiler Screenshot](https://i.imgur.com/JMUbzL9.png)

## Known issues

* pprof-upload should recognize profile type from the input file.
