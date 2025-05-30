# Go: Cross Compile Binary 

[View original Gist on GitHub](https://gist.github.com/Integralist/f21d57a8fcada8d4c2ac79bece4337b4)

**Tags:** #go #compiler #build

## Golang Cross Compile Binary.bash

```shell
env GOOS=darwin  GOARCH=amd64 go build -o fastly     ./cmd/fastly/main.go
env GOOS=windows GOARCH=amd64 go build -o fastly.exe ./cmd/fastly/main.go
```

