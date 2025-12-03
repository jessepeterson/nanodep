package apinext

//go:generate oa2js -o ErrorResponse.json ../../docs/openapi.yaml ErrorResponse
//go:generate oa2js -o BypassCodeResponse.json ../../docs/openapi.yaml BypassCodeResponse
//go:generate go-jsonschema -p $GOPACKAGE --tags json --only-models --output schema.go ErrorResponse.json BypassCodeResponse.json
//go:generate rm -f ErrorResponse.json BypassCodeResponse.json
