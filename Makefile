




cross-%:
	GOOS=linux CGO_ENABLED=0 GOARCH=arm64  go build -mod=vendor -o iam .\\cmd\\server\\main.go
