package info

// 通过 ldflags 在构建时注入，例如:
// go build -ldflags="-X main.Version=v1.0.0" -o Doubao-Susurro.exe .
var Version = "dev"
