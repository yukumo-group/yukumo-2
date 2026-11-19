module yukumo-script/cli

go 1.27

require (
	github.com/spf13/cobra v1.10.2
	gopkg.in/yaml.v3 v3.0.1
	yukumo-script/core/api v0.0.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/tetratelabs/wazero v1.12.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	yukumo-script/core v0.0.0 // indirect
	yukumo-script/core/port v0.0.0 // indirect
	yukumo-script/runtime v0.0.0 // indirect
)

replace (
	yukumo-script/core => ../core
	yukumo-script/core/api => ../core/api
	yukumo-script/core/port => ../core/port
	yukumo-script/runtime => ../runtime
)
