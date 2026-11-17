module yukumo-script/wasmtest

go 1.27

require (
	yukumo-script/core/api v0.0.0
	yukumo-script/core/port v0.0.0
)

require (
	gopkg.in/yaml.v3 v3.0.1 // indirect
	yukumo-script/core v0.0.0 // indirect
)

replace (
	yukumo-script/core => ../core
	yukumo-script/core/api => ../core/api
	yukumo-script/core/port => ../core/port
	yukumo-script/runtime => ../runtime
)
