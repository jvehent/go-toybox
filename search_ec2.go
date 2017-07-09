package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/stripe/aws-go/aws"
	"github.com/stripe/aws-go/gen/ec2"
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
		err  error
		conf conf
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

	// create a new client to EC2 api
	creds := aws.Creds(conf.Credentials.AccessKey, conf.Credentials.SecretKey, "")
	cli := ec2.New(creds, "us-east-1", nil)

	fireq := ec2.Filter{
		Name:   aws.String("private-ip-address"),
		Values: []string{"172.30.200.13"},
	}
	direq := ec2.DescribeInstancesRequest{
		Filters: []ec2.Filter{fireq},
	}
	resp, err := cli.DescribeInstances(&direq)
	if err != nil {
		panic(err)
	}
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Printf("%s\t%s", *instance.InstanceID, *instance.PrivateIPAddress)
			for _, tag := range instance.Tags {
				fmt.Printf("\t%s:%s", *tag.Key, *tag.Value)
			}
			fmt.Printf("\n")
		}
	}
}
