package sssly

import (
	"bytes"
	"errors"
	"github.com/luisfurquim/goose"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Opt map[string]interface{}

type Sssly struct {
	Bucket string
	BasePath string
	Client *s3.Client
}

type WriteCloser struct {
	bytes.Buffer
	cli *Sssly
	key string
}

type GooseG struct {
	Init goose.Alert
	Storage goose.Alert
}

var Goose GooseG = GooseG{
	Init: goose.Alert(2),
	Storage: goose.Alert(2),
}

var ErrOptionRequiredRegion error = errors.New("Option required: region")
var ErrOptionWrongTypeRegion error = errors.New("Option wrong type: region")

var ErrOptionWrongTypeHttpClient error = errors.New("Option wrong type: http client")
var ErrOptionWrongTypeHttpTransport error = errors.New("Option wrong type: http transport")

var ErrOptionRequiredCredentials error = errors.New("Option required: credentials")
var ErrOptionWrongTypeCredentials error = errors.New("Option wrong type: credentials")

var ErrOptionRequiredProfile error = errors.New("Option required: profile")
var ErrOptionWrongTypeProfile error = errors.New("Option wrong type: profile")

var ErrOptionRequiredBucket error = errors.New("Option required: bucket")
var ErrOptionWrongTypeBucket error = errors.New("Option wrong type: bucket")

var ErrOptionRequiredEndpoint error = errors.New("Option required: endpoint")
var ErrOptionWrongTypeEndpoint error = errors.New("Option wrong type: endpoint")

var ErrOptionRequiredBasePath error = errors.New("Option required: base path")
var ErrOptionWrongTypeBasePath error = errors.New("Option wrong type: base path")

