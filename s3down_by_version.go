// sample code to download content from S3 bucket by version
//
// export AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY in your environment
//
// then run $ go run s3down_by_version.go us-east-1 mybucket myfolder/mykey zG77EWreHU9Jqswz1QA6mEKP4QaeUlLa
//					  { region} {bucket} {..doc key...} {.........version ID...........}
//	{
//	  AcceptRanges: "bytes",
//	  Body: buffer(0xc20836cf00),
//	  ContentLength: 887,
//	  ContentType: "application/octet-stream",
//	  ETag: "\"0e2a5820fdb608ff3ff8f0bdda6ee378\"",
//	  LastModified: 2015-05-13 13:54:28 +0000 UTC,
//	  Metadata: {
//	  },
//	  ServerSideEncryption: "AES256",
//	  VersionID: "zG77EWreHU9Jqswz1QA6mEKP4QaeUlLa"
//	}
//	-----BEGIN RSA PRIVATE KEY-----
//	MIICXAIBAAKBgQCsqIsQOkcNQamH+JHJL8vQI0OjvyPnL3oN0uLnCSKU5e166L+W
//	...
//	-----END RSA PRIVATE KEY-----
//
package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <region> <bucket> <key> (<version>)\n", os.Args[0])
		os.Exit(1)
	}
	svc := s3.New(aws.NewConfig().WithRegion(os.Args[1]).WithMaxRetries(1))

	params := &s3.GetObjectInput{
		Bucket: aws.String(os.Args[2]), // Required
		Key:    aws.String(os.Args[3]), // Required
	}
	if len(os.Args) > 4 {
		params.VersionID = aws.String(os.Args[4])
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS error with Code, Message, and original error (if any)
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			// This case should never be hit, the SDK should always return an
			// error which satisfies the awserr.Error interface.
			fmt.Println(err.Error())
		}
	}

	// Pretty-print the response data.
	fmt.Println(awsutil.Prettify(resp))

	// Print the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", body)
}
