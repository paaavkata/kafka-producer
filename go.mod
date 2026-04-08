module producer

go 1.23.3

require (
	github.com/paaavkata/go-logger v0.1.1
	github.com/twmb/franz-go v1.18.1
)

replace github.com/paaavkata/go-logger => ../shared-libs/go-logger

require (
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/segmentio/kafka-go v0.4.47 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.9.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
