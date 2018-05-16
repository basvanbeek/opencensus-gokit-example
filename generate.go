package ocgokitexample

//go:generate protoc -I=$GOPATH/src -I=. qr/transport/grpc/pb/svcqr.proto --go_out=plugins=grpc:.
//go:generate go build -o build/elegantmonolith elegantmonolith/main.go
//go:generate go build -o build/frontend frontend/cmd/main.go
//go:generate go build -o build/QRGenerator qr/cmd/main.go
