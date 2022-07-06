# go-multicam

Go language wrapper for the Euresys MultiCam SDK.

https://documentation.euresys.com/Products/MULTICAM/MULTICAM/Content/00_Home/home.htm

This package currently only works on Linux platforms.

To use it you must obtain a copy of the Euresys MultiCam SDK for Linux from Euresys directly.

## How To Use

Coming soon.

## Installing

Instructions will come later.


## Using with Docker

Download the latest Euresys SDK file, e.g. `multicam-linux-x86_64-6.18.3.4935.tar.gz` and put it into a directory named `multicam-linux` in this project. Note that this file will not be saved into the Github repo, you must download it yourself.

Once you have obtained the file, you can build the Docker image by running:


```
docker build -t multicam:latest .

```

If you are using a different version of the Multicam SDK, you can use the `MULTICAM_RELEASE` build argument like this::


```
docker build --build-arg MULTICAM_RELEASE=6.18.2.4781 -t multicam:latest .

```

Once you have built it, you can run it like this:


```
docker run --privileged multicam:latest

```

You should see output like this:

```
Driver was opened...
Boards detected: 2
FULL_XR_3489
FULL_XR_929
Done.
```


## Why it exists

Computer vision applications written in Go can easily take advantage of multi-core concurrency while also having access to packages such as GoCV (https://gocv.io). The `go-multicam` package now makes is possible for Go programs to connect to cameras and line scanners that are generally used for industrial applications.
