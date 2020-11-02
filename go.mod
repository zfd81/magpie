module github.com/zfd81/magpie

go 1.14

require (
	github.com/antonmedv/expr v1.8.9
	github.com/boltdb/bolt v1.3.1
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fatih/color v1.9.0
	github.com/golang/protobuf v1.4.1
	github.com/google/uuid v1.1.2
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	google.golang.org/grpc v1.27.0
	sigs.k8s.io/yaml v1.2.0 // indirect
	vitess.io/vitess v0.7.0
)

replace (
	github.com/golang/protobuf => github.com/golang/protobuf v1.4.3
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.25.0
)
