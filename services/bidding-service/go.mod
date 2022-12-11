module github.com/PanGan21/bidding-service

go 1.18

replace (
	github.com/PanGan21/pkg/auth => ../../pkg/auth
	github.com/PanGan21/pkg/entity => ../../pkg/entity
	github.com/PanGan21/pkg/httpserver => ../../pkg/httpserver
	github.com/PanGan21/pkg/logger => ../../pkg/logger
	github.com/PanGan21/pkg/pagination => ../../pkg/pagination
	github.com/PanGan21/pkg/postgres => ../../pkg/postgres
	github.com/PanGan21/pkg/utils => ../../pkg/utils
)

require (
	github.com/PanGan21/pkg/httpserver v0.0.0-00010101000000-000000000000
	github.com/PanGan21/pkg/logger v0.0.0-00010101000000-000000000000
	github.com/PanGan21/pkg/postgres v0.0.0-00010101000000-000000000000
	github.com/gin-gonic/gin v1.8.1
	github.com/ilyakaznacheev/cleanenv v1.4.1
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/Masterminds/squirrel v1.5.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
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
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/rs/zerolog v1.28.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)