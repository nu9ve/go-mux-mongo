# go-mux-mongo project template
this project serves as a starting point for a full stack application working together with the [react-redux-router](https://github.com/nu9ve/react-redux-router) template.

## dependencies
- golang 1.17+
- [mongo db driver](https://docs.mongodb.com/manual/tutorial/install-mongodb-on-os-x/)

## running
in the project directory, you should:

### `go install`
install all the dependencies described in go.mod

### `go build cmd/main.go`
build the local package with starting point in cmd/main.go

### `./main`
run app, if you see a similar message you're golden

```shell
2022/02/21 09:46:25 Listening on port  8088
```