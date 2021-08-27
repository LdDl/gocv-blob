# gocv-blob [![GoDoc](https://godoc.org/github.com/LdDl/gocv-blob?status.svg)](https://godoc.org/github.com/LdDl/gocv-blob)[![Sourcegraph](https://sourcegraph.com/github.com/LdDl/gocv-blob/-/badge.svg)](https://sourcegraph.com/github.com/LdDl/gocv-blob?badge)[![Go Report Card](https://goreportcard.com/badge/github.com/LdDl/gocv-blob)](https://goreportcard.com/report/github.com/LdDl/gocv-blob)[![GitHub tag](https://img.shields.io/github/tag/LdDl/gocv-blob.svg)](https://github.com/LdDl/gocv-blob/releases)
Blob tracking via GoCV package

## Table of Contents

- [About](#about)
- [Installation](#installation)
- [Usage](#usage)
- [Support](#support)
- [Thanks](#thanks)

## About
This small package implements basics of blob tracking: simple centroid and [Kalman filter](https://en.wikipedia.org/wiki/Kalman_filter)-based tracking

There are additional functions for checking if blob crossed horizontal (or even oblique) line.

## Installation

First of all you need OpenCV to be installed on your operation system. Also you need [GoCV](https://github.com/hybridgroup/gocv) package to be installed too. Please see ref. here https://github.com/hybridgroup/gocv#how-to-install

Then you are good to go with:
```shell
go get github.com/LdDl/gocv-blob/v2
## or (if you want to use legacy version):
## go get github.com/LdDl/gocv-blob
```

p.s. do not be worried when you see *can't load package: package github.com/LdDl/gocv-blob: no Go files....* - this is just warning.

## Usage

It's pretty straightforward (pseudocode'ish).

```go
// 1. Define global set of blobs
global_blobs = blob.NewBlobiesDefaults()

// 2. Define new blob objects
new_blob1 = blob.NewSimpleBlobie(image.Rectangle, how many points to store in track, class ID of object , class name of object)
new_blob2 = blob.NewSimpleBlobie(image.Rectangle, how many points to store in track, class ID of object , class name of object)
// You can use NewKalmanBlobie if needed

// 3. Append data to temporary set of blobs
tmp_blobs = []*blob.Blobie{new_blob1, new_blob2}

// 4. Compare blobs ()
global_blobs.MatchToExisting(tmp_blobs)

// 5. Repeat steps 2-4 every time when you find new objects on images. MatchToExisting() will update existing blobs and register new ones.
```

More informative example is here: [v2/array_tracker_test.go](v2/array_tracker_test.go)

**FOR LEGACY v1**:
<details>
<summary>Click to expand</summary>

```go
// 1. Define global set of blobs
global_blobs = blob.NewBlobiesDefaults()

// 2. Define new blob objects
new_blob1 = blob.NewBlobie(image.Rectangle, how many points to store in track, class ID of object , class name of object)
new_blob2 = blob.NewBlobie(image.Rectangle, how many points to store in track, class ID of object , class name of object)

// 3. Append data to temporary set of blobs
tmp_blobs = []*blob.Blobie{}
tmp_blobs = append(tmp_blobs, new_blob1)
tmp_blobs = append(tmp_blobs, new_blob2)

// 4. Compare blobs ()
global_blobs.MatchToExisting(tmp_blobs)

// 5. Repeat steps 2-4 every time when you find new objects on images. MatchToExisting() will update existing blobs and register new ones.
```
</details>

## Support

If you have troubles or questions please [open an issue](https://github.com/LdDl/gocv-blob/issues/new).

## Thanks
Big thanks to creators and developers of [GoCV](https://gocv.io/) for providing bindings to OpenCV
