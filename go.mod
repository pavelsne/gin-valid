module github.com/G-Node/gin-valid

go 1.12

require (
	github.com/G-Node/gin-cli v0.0.0-20190812205442-4a7c8da2e5c5
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/gogits/go-gogs-client v0.0.0-20190710002546-4c3c18947c15
	github.com/gogs/go-gogs-client v0.0.0-20190710002546-4c3c18947c15
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190717161051-705d9623b7c1
	github.com/go-xorm/core => xorm.io/core v0.6.3
)
