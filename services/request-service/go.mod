module github.com/PanGan21/request-service

go 1.20

replace (
	github.com/PanGan21/pkg/auth => ../../pkg/auth
	github.com/PanGan21/pkg/entity => ../../pkg/entity
	github.com/PanGan21/pkg/httpserver => ../../pkg/httpserver
	github.com/PanGan21/pkg/logger => ../../pkg/logger
	github.com/PanGan21/pkg/messaging => ../../pkg/messaging
	github.com/PanGan21/pkg/pagination => ../../pkg/pagination
	github.com/PanGan21/pkg/postgres => ../../pkg/postgres
	github.com/PanGan21/pkg/utils => ../../pkg/utils
)

require (
	github.com/PanGan21/pkg/auth v0.0.0-00010101000000-000000000000
	github.com/PanGan21/pkg/httpserver v0.0.0-00010101000000-000000000000
	github.com/PanGan21/pkg/logger v0.0.0-00010101000000-000000000000
	github.com/PanGan21/pkg/messaging v0.0.0-00010101000000-000000000000
	github.com/PanGan21/pkg/postgres v0.0.0-00010101000000-000000000000
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.0
	github.com/ilyakaznacheev/cleanenv v1.4.2
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/Masterminds/squirrel v1.5.3 // indirect
	github.com/PanGan21/pkg/entity v0.0.0-00010101000000-000000000000 // indirect
	github.com/PanGan21/pkg/utils v0.0.0-00010101000000-000000000000 // indirect
	github.com/bytedance/sonic v1.8.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.11.2 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.13 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/rs/zerolog v1.28.0 // indirect
	github.com/segmentio/kafka-go v0.4.38 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.9 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
