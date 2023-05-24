#! /bin/sh
set -ex

# temp start
find ./ -type f -name "*pb.go" | xargs rm -rf 
# temp end


protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative server/grpc/proto/auth/auth.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative server/grpc/proto/user/user.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative server/grpc/proto/blog/blog.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative server/grpc/proto/comment/comment.proto