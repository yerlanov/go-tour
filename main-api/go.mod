module github.com/yerlanov/go-tour/main-api

go 1.20

require (
	github.com/gin-gonic/gin v1.9.0
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/google/uuid v1.3.0
	github.com/yerlanov/go-tour/grpctest2 v0.0.0-20230328064720-56bce79b6ba3
	github.com/yerlanov/go-tour/common v0.0.0-20230328064720-56bce79b6ba3
	go.mongodb.org/mongo-driver v1.9.0
	golang.org/x/oauth2 v0.6.0
	google.golang.org/appengine v1.6.7
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cloud.google.com/go/compute v1.19.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/bytedance/sonic v1.8.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.11.2 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.9 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230327215041-6ac7f18bb9d5 // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace (
	github.com/yerlanov/go-tour/grpctest2 => ../grpctest2
	github.com/yerlanov/go-tour/common => ../common
)