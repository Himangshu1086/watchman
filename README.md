# watchman

## How to Install :
Install the package 
```
go install github.com/Himangshu1086/watchman
```

## How To Use :
```
go run watchman.go <port> <server-filepath>
```
eg : `go run watchman.go 8000 cmd/user-service.go`

## Important :
**Add Go bin directory to your PATH**:
By default, Go installs binaries to the `GOPATH/bin` directory. You need to add this directory to your PATH. If you are using `zsh`, you can do this by adding the following line to your `.zshrc` file:
```
export PATH=$PATH:$(go env GOPATH)/bin
```
**Reload your .zshrc file**:
After updating your `.zshrc` file, you need to reload it:
```
source ~/.zshrc
```

## Alternate way :
- Copy the watchman.go file and add it into the source directory
- And run ``go run watchman.go <port> <server-filepath>``
