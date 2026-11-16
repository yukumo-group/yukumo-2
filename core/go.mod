module yukumo-script/core

go 1.27

require (
	gopkg.in/yaml.v3 v3.0.1
	yukumo-script/runtime v0.0.0
)

require (
	github.com/tetratelabs/wazero v1.12.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
)

replace yukumo-script/runtime => ../runtime
