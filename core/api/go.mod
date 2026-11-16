module yukumo-script/core/api

go 1.27

require (
	yukumo-script/core v0.0.0
	yukumo-script/runtime v0.0.0
)

require (
	github.com/tetratelabs/wazero v1.12.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	yukumo-script/core => ..
	yukumo-script/runtime => ../../runtime
)
