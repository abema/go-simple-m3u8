# go-simple-m3u8
------

[![Go Reference](https://pkg.go.dev/badge/github.com/abema/go-simple-m3u8.svg)](https://pkg.go.dev/github.com/abema/go-simple-m3u8)
![Test](https://github.com/abema/go-simple-m3u8/actions/workflows/test.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/abema/go-simple-m3u8/badge.svg)](https://coveralls.io/github/abema/go-simple-m3u8)
[![Go Report Card](https://goreportcard.com/badge/github.com/abema/go-simple-m3u8)](https://goreportcard.com/report/github.com/abema/go-simple-m3u8)

go-simple-m3u8 is a Go library for encoding and decoding M3U8 files.
This library retains all tags and their attributes using a map, ensuring that no tags are lost even if they are not explicitly supported.
While useful methods are provided, you can access all tags and attributes by their names if needed.