// Package spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package spec

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	externalRef0 "github.com/trustbloc/vcs/pkg/restapi/v1/common"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	pathPrefix := path.Dir(pathToFile)

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(pathPrefix, "./common.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
	}
	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
