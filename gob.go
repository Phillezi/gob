// Package gob is the root of the gob library.
//
// This library aims to make it possible to build go projects
// based on build recipies in go code.
//
// This is heavily inspired by how zig does it, with build.zig.
//
// # Usage
//
// You can use this library by making a go file that looks something like this:
//
// //go build ignore
// package main
//
// import "github.com/phillezi/gob"
//
//	func main() {
//		b := gob.New()
//		b.Add("", gob.Static())
//		b.Run()
//	}
package gob
