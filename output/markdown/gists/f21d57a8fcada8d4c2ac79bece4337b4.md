# [Golang Cross Compile Binary] #go #golang #compile #crosscompile

## Golang Cross Compile Binary.bash

```shell
env GOOS=darwin  GOARCH=amd64 go build -o fastly     ./cmd/fastly/main.go
env GOOS=windows GOARCH=amd64 go build -o fastly.exe ./cmd/fastly/main.go
```

