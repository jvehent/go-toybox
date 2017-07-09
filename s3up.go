package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/stripe/aws-go/aws"
	"github.com/stripe/aws-go/gen/s3"
	"os"
)

// conf takes an AWS configuration from a file in ~/.awsgo
// example:
//
// [credentials]
//    accesskey = "AKI...."
//    secretkey = "mw0...."
//
type conf struct {
	Credentials struct {
		AccessKey string
		SecretKey string
	}
}

func main() {
	var (
		err         error
		conf        conf
		bucket      string = "testawsgo" // change to your convenience
		fd          *os.File
		contenttype string = "binary/octet-stream"
	)
	// obtain credentials from ~/.awsgo
	credfile := os.Getenv("HOME") + "/.awsgo"
	_, err = os.Stat(credfile)
	if err != nil {
		fmt.Println("Error: missing credentials file in ~/.awsgo")
		os.Exit(1)
	}
	err = gcfg.ReadFileInto(&conf, credfile)
	if err != nil {
		panic(err)
	}

	// create a new client to S3 api
	creds := aws.Creds(conf.Credentials.AccessKey, conf.Credentials.SecretKey, "")
	cli := s3.New(creds, "us-east-1", nil)

	// open the file to upload
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <inputfile>\n", os.Args[0])
		os.Exit(1)
	}
	fi, err := os.Stat(os.Args[1])
	if err != nil {
		fmt.Printf("Error: no input file found in '%s'\n", os.Args[1])
		os.Exit(1)
	}
	fd, err = os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	// create a bucket upload request and send
	objectreq := s3.PutObjectRequest{
		ACL:           aws.String("public-read"),
		Bucket:        aws.String(bucket),
		Body:          fd,
		ContentLength: aws.Integer(int(fi.Size())),
		ContentType:   aws.String(contenttype),
		Key:           aws.String(fi.Name()),
	}
	_, err = cli.PutObject(&objectreq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("%s\n", "https://s3.amazonaws.com/"+bucket+"/"+fi.Name())
	}

	// list the content of the bucket
	//listreq := s3.ListObjectsRequest{
	//	Bucket: aws.StringValue(&bucket),
	//}
	//listresp, err := cli.ListObjects(&listreq)
	//if err != nil {
	//	panic(err)
	//}
	//if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//} else {
	//	fmt.Printf("Content of bucket '%s': %d files\n", bucket, len(listresp.Contents))
	//	for _, obj := range listresp.Contents {
	//		fmt.Println("-", *obj.Key)
	//	}
	//}
}
