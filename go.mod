module github.com/G-Node/gin-valid

go 1.16

require (
	github.com/G-Node/gin-cli v0.0.0-20190819162807-7786caf50bbd
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/gogits/go-gogs-client v0.0.0-20190710002546-4c3c18947c15
	github.com/gogs/go-gogs-client v0.0.0-20190710002546-4c3c18947c15
	github.com/google/uuid v1.1.2
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pelletier/go-toml v1.5.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/yaml.v2 v2.2.8
)

replace (
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190717161051-705d9623b7c1
	github.com/go-xorm/core => xorm.io/core v0.6.3
)
