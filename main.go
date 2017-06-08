package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	logrus "github.com/Sirupsen/logrus"
	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	route53 "github.com/aws/aws-sdk-go/service/route53"
)

const (
	ipCheckURL = "http://curlmyip.org"
)

var (
	awsRegion      = flag.String("awsregion", "eu-west-1", "The AWS region to connect to")
	debug          = flag.Bool("debug", false, "Debug logging")
	dnsTTL         = flag.Int64("ttl", 900, "TTL for any DNS records created")
	hostedZoneId   = flag.String("zoneid", "XYWQJHASDJHG.", "The Route53 Zone-ID")
	hostName       = flag.String("hostname", "host.domain.com.", "The hostname to update")
	hostUpdateFreq = flag.Duration("frequency", 60*time.Minute, "How often to update the record")
	route53Service *route53.Route53
)

func init() {
	// Parse the command-line arguments:
	flag.Parse()

	// Set up a logger:
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {

	// Create the service's client with the session.
	logrus.Debug("Connecting to AWS")
	route53Service = route53.New(session.New(), aws.NewConfig().WithRegion(*awsRegion))

	updateRecord()

	// A ticker to tell us when its time to update DNS:
	updateTimer := time.Tick(*hostUpdateFreq)

	// Wait for the updateTimer to tell us its time to update the zone config:
	for {
		select {
		case <-updateTimer:
			logrus.Infof("Updating DNS record on schedule (%v)", hostUpdateFreq)
			updateRecord()
		}
	}

}

func updateRecord() error {

	logrus.Debugf("Updating DNS record (%v)", hostUpdateFreq)

	// Find out what our IP address is:
	ipCheckResponse, err := http.Get(ipCheckURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Couldn't check our IP address")
		return err
	}
	ourIPAddress, err := ioutil.ReadAll(ipCheckResponse.Body)
	defer ipCheckResponse.Body.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Couldn't read ipCheckResponse.Body")
		return err
	}

	// Make a ChangeResourceRecordSetsInput:
	changeResourceRecordSetsInput := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				&route53.Change{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: hostName,
						TTL:  dnsTTL,
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							&route53.ResourceRecord{
								Value: aws.String(string(ourIPAddress)),
							},
						},
					},
				},
			},
		},
		HostedZoneId: hostedZoneId,
	}

	// Make the ChangeResourceRecordSets request:
	_, err = route53Service.ChangeResourceRecordSets(changeResourceRecordSetsInput)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Couldn't update DNS record")
		return err
	}

	logrus.WithFields(logrus.Fields{"hostname": *hostName, "address": string(ourIPAddress), "ttl": *dnsTTL, "zone_id": *hostedZoneId}).Infof("Updated DNS record")

	return nil
}
