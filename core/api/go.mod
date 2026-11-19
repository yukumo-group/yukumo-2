module yukumo-script/core/api

go 1.27

require (
	github.com/invopop/jsonschema v0.13.0
	yukumo-script/core v0.0.0
	yukumo-script/core/port v0.0.0
	yukumo-script/runtime v0.0.0
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/tetratelabs/wazero v1.12.0 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	golang.org/x/sys v0.44.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	yukumo-script/core => ..
	yukumo-script/core/port => ../port
	yukumo-script/runtime => ../../runtime
)
