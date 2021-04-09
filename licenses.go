package main

// NOTE:
// This main package is a workaround for binary license scanning that forces transitive dependencies in
// Code we're distributing to be analyzed
import (
	_ "github.com/nitrictech/go-sdk/api"
	_ "github.com/nitrictech/go-sdk/faas"
)

func main() {}
