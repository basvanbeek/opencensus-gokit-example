package ocgokitexample

//go:generate protoc -I$GOPATH/src -I. services/device/transport/pb/svcdevice.proto --go_out=plugins=grpc:.
//go:generate protoc -I$GOPATH/src -I. services/qr/transport/pb/qr.proto --go_out=plugins=grpc:. --twirp_out=.
//go:generate protoc -I$GOPATH/src -I. services/event/transport/pb/event.proto --go_out=plugins=grpc:. --twirp_out=.
//go:generate go build -tags sqlite3 -o build/cli clients/cli/main.go
//go:generate go build -tags sqlite3 -o build/ocg-elegantmonolith services/elegantmonolith/main.go
//go:generate go build -tags sqlite3 -o build/ocg-event services/event/cmd/main.go
//go:generate go build -tags sqlite3 -o build/ocg-qrgenerator services/qr/cmd/main.go
//go:generate go build -tags sqlite3 -o build/ocg-device services/device/cmd/main.go
//go:generate go build -tags sqlite3 -o build/ocg-frontend services/frontend/cmd/main.go
