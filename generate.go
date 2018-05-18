package ocgokitexample

//go:generate protoc -I=$GOPATH/src -I=. device/transport/grpc/pb/svcdevice.proto --go_out=plugins=grpc:.
//go:generate protoc -I=$GOPATH/src -I=. qr/transport/grpc/pb/svcqr.proto --go_out=plugins=grpc:.
//go:generate go build -tags sqlite3 -o build/elegantmonolith elegantmonolith/main.go
//go:generate go build -tags sqlite3 -o build/device device/cmd/main.go
//go:generate go build -tags sqlite3 -o build/frontend frontend/cmd/main.go
//go:generate go build -tags sqlite3 -o build/QRGenerator qr/cmd/main.go
