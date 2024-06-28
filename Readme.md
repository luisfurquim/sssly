Sssly is a simple wrapper for AWS S3 access.

```Go
package main

import (
	"io"
	"fmt"
	"github.com/luisfurquim/sssly"
)

func main() {
	var cli *sssly.Sssly
	var err error
	var wr io.WriteCloser
	var rd io.ReadCloser
	var dir []string
	var buf []byte

	cli, err = sssly.New(sssly.Opt{
		"region": "sa-east-1", // or "aws-cn-global", etc
		"credentials": "credentials", // filenames (comma separated) of credential files, see file format in comment below
		"profile": "myprofile", // profile name (section in the credential files)
		"bucket": "mybucket", // bucket name in S3
		"endpoint": "https://example.com", // URL of S3 endpoint
		"base-path": "temp/", // base pathname in S3 bucket, all keys provided in method callings will be appended to base-path
	})

/*

[myprofile]
aws_access_key_id = the_id_provided_by_the_S3_admin
aws_secret_access_key = secretsecretsecretsecretsecretsecretsecretsecret

[myotherprofile]
aws_access_key_id = the_other_id_provided_by_the_S3_admin
aws_secret_access_key = othersecretothersecretothersecretothersecret

*/

	if err != nil {
		// whatever
	}

	// List files in bucket
	dir, err = cli.Dir()
	if err != nil {
		// whatever
	}

	// Delete file(s) in bucket/base-path
	err = cli.Delete("some_junk_file.txt", ...)
	if err != nil {
		// whatever
	}

	// Create file in bucket/base-path a return its WriteCloser
	// Note: it bufferizes all content in memory, don't use to transfer huge files (use Upload for these cases)
	wr = cli.NewWriteCloser("some_new_file.txt")
	fmt.Fprintf(wr, "some file content")

	err = wr.Close() // MUST be called, otherwise no data is transfered to S3
	if err != nil {
		// whatever
	}

	// Uploads file to bucket/base-path, it does not buffers in memory, so it is suitable for huge files
	err = cli.Upload("s3-file.txt", "local-file.txt")
	if err != nil {
		// whatever
	}

	// Opens a file located in bucket/base-path and returns an io.Reader
	rd, err = cli.NewReadCloser("some_new_file.txt")
	if err != nil {
		// whatever
	}
	// Don't forget to close it!
	defer rd.Close()

	buf, err = io.ReadAll(rd)
	if err != nil {
		// whatever
	}

}

```