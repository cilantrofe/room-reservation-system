module github.com/Quizert/room-reservation-system/BookingSvc

go 1.23.3

require (
	github.com/Quizert/room-reservation-system/HotelSvc v0.0.0-20241221001323-5994fb8e8146
	github.com/Quizert/room-reservation-system/Libs v0.0.0-20241221001323-5994fb8e8146
	github.com/golang/mock v1.6.0
	github.com/jackc/pgx/v4 v4.18.3
	github.com/segmentio/kafka-go v0.4.47
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.10.0
	google.golang.org/grpc v1.69.2
)

require github.com/stretchr/testify v1.9.0 // indirect

require (
	github.com/Quizert/room-reservation-system/AuthSvc v0.0.0-20241222021346-c333398f5c3b
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/pierrec/lz4/v4 v4.1.16 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/protobuf v1.36.0 // indirect
)
