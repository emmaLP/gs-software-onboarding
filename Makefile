proto_gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/grpc/proto/hackernews.proto

proto_mock:
	mockgen -package=proto -source=pkg/grpc/proto/hackernews_grpc.pb.go -destination=pkg/grpc/proto/mockClient.go

proto: proto_gen proto_mock

test_int_coverage:
	go test -tags integration ./... -coverprofile cover.out && cat cover.out | grep -v "mock" > coverexclude.out && go tool cover -func=coverexclude.out