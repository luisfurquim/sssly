package sssly

import (
	"context"
	"strings"
	"net/http"
	"crypto/tls"
	"path/filepath"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func New(opt Opt) (*Sssly, error) {
	var cli *Sssly
	var hcli *http.Client
	var tr *http.Transport
	var cfg aws.Config
	var err error
	var region string
	var op, op2 interface{}
	var credFiles string
	var ok bool
	var endpoint string
	var profile string

	if op, ok = opt["region"] ; !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionRequiredRegion)
		return nil, ErrOptionRequiredRegion
	}

	if region, ok = op.(string); !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionWrongTypeRegion)
		return nil, ErrOptionWrongTypeRegion
	}

	if op, ok = opt["http-client"] ; !ok {
		if op2, ok = opt["http-transport"] ; !ok {
			tr = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
		} else {
			if tr, ok = op2.(*http.Transport); !ok {
				Goose.Init.Logf(1, "Error %s", ErrOptionWrongTypeHttpTransport)
				return nil, ErrOptionWrongTypeHttpTransport
			}
		}
		hcli = &http.Client{
			Transport: tr,
		}
	} else {
		if hcli, ok = op.(*http.Client); !ok {
			Goose.Init.Logf(1, "Error %s", ErrOptionWrongTypeHttpClient)
			return nil, ErrOptionWrongTypeHttpClient
		}
	}

	if op, ok = opt["credentials"] ; !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionRequiredCredentials)
		return nil, ErrOptionRequiredCredentials
	}

	if credFiles, ok = op.(string); !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionWrongTypeCredentials)
		return nil, ErrOptionWrongTypeCredentials
	}

	if op, ok = opt["profile"] ; !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionRequiredProfile)
		return nil, ErrOptionRequiredProfile
	}

	if profile, ok = op.(string); !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionWrongTypeProfile)
		return nil, ErrOptionWrongTypeProfile
	}

	cli = &Sssly{}

	if op, ok = opt["bucket"] ; !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionRequiredBucket)
		return nil, ErrOptionRequiredBucket
	}

	if cli.Bucket, ok = op.(string); !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionWrongTypeBucket)
		return nil, ErrOptionWrongTypeBucket
	}

	if op, ok = opt["endpoint"] ; !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionRequiredEndpoint)
		return nil, ErrOptionRequiredEndpoint
	}

	if endpoint, ok = op.(string); !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionWrongTypeEndpoint)
		return nil, ErrOptionWrongTypeEndpoint
	}

	if op, ok = opt["base-path"] ; !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionRequiredBasePath)
		return nil, ErrOptionRequiredBasePath
	}

	if cli.BasePath, ok = op.(string); !ok {
		Goose.Init.Logf(1, "Error %s", ErrOptionWrongTypeBasePath)
		return nil, ErrOptionWrongTypeBasePath
	}

	cli.BasePath = filepath.Clean(cli.BasePath) + "/"

	cfg, err = config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithHTTPClient(hcli),
		config.WithSharedCredentialsFiles(strings.Split(credFiles,",")),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		Goose.Init.Logf(1, "Error initializing config: %s", err)
		return nil, err
	}

	cli.Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = &endpoint
		o.UsePathStyle = true // for automatic path setup on keys
	})

	return cli, nil
}

