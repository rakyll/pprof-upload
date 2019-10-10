# pprof-upload

[![Build Status](https://travis-ci.com/rakyll/pprof-upload.svg?token=Quf3mWAszVwsMXDMWxkm&branch=master)](https://travis-ci.com/rakyll/pprof-upload)

Uploads pprof files to Google Cloud Profiler.

## Requirements

* Enable [Google Cloud Profiler](https://cloud.google.com/profiler/).
* If running outside of Google Compute Engine, install
  [gcloud](https://cloud.google.com/sdk/gcloud/) if you haven't   
  and run `gcloud auth application-default login`.

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
[net/http/pprof](https://golang.org/pkg/net/http/pprof) package.

```
$ curl http://localhost:6060/debug/pprof/profile?seconds=30 > pprof.out
$ pprof-upload -target=webserver
https://console.cloud.google.com/profiler/webserver;type=CPU?project=PROJECT
```

![Cloud Profiler Screenshot](https://i.imgur.com/KRQcmZ5.png)

## Known issues

* pprof-upload should recognize profile type from the input file.
* pprof-upload resets profile data, it should be optional.
