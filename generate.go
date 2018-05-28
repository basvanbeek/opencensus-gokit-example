package ocgokitexample

//go:generate protoc -I$GOPATH/src -I. services/device/transport/grpc/pb/svcdevice.proto --go_out=plugins=grpc:.
//go:generate protoc -I$GOPATH/src -I. services/qr/transport/pb/qr.proto --go_out=plugins=grpc:. --twirp_out=.
//go:generate protoc -I$GOPATH/src -I. services/event/transport/pb/event.proto --go_out=plugins=grpc:. --twirp_out=.
//go:generate go build -tags sqlite3 -o build/cli services/cli/main.go
//go:generate go build -tags sqlite3 -o build/elegantmonolith services/elegantmonolith/main.go
//go:generate go build -tags sqlite3 -o build/event services/event/cmd/main.go
//go:generate go build -tags sqlite3 -o build/QRGenerator services/qr/cmd/main.go
//go:generate go build -tags sqlite3 -o build/device services/device/cmd/main.go
//go:generate go build -tags sqlite3 -o build/frontend services/frontend/cmd/main.go
