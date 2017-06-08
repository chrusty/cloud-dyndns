package main

import (
	"flag"
	"time"

	logrus "github.com/Sirupsen/logrus"
	route53 "github.com/aws/aws-sdk-go/service/route53"
)

var (
	dnsTTL         = flag.Int("dnsttl", 300, "TTL for any DNS records created")
	hostName       = flag.String("hostname", "host.domain.com.", "The hostname to update")
	hostUpdateFreq = flag.Duration("hostupdate", "60m", "How often to update the record")
	hostedZoneId   = flag.String("zoneid", "XYWQJHASDJHG.", "The Route53 Zone-ID")
	awsRegion      = flag.String("awsregion", "eu-west-1", "The AWS region to connect to")
	logger         *logrus.Logger
)

func init() {
	// Parse the command-line arguments:
	flag.Parse()

	// Set up a logger:
	logger = logrus.Logger{
		Formatter: &logrus.JSONFormatter{},
	}

}

func main() {

	// Create a session to share configuration, and load external configuration.
	sess := session.Must(session.NewSession())

	// Create the service's client with the session.
	route53Service := route53.New(session.New(), aws.NewConfig().WithRegion(awsRegion))

	// A ticker to tell us when its time to update DNS:
	updateTimer := time.Tick(hostUpdateFreq)

	// Wait for the updateTimer to tell us its time to update the zone config:
	for {
		select {
		case <-updateTimer:
			logger.Debugf("Updating DNS record on schedule (%v)", hostUpdateFreq)
			updateRecord()
		}
	}

}

func updateRecord() {

	changeResourceRecordSetsInput := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				Action: "UPSERT",
				ResourceRecordSet: &route53.ResourceRecordSet{
					Name: hostname,
				},
			},
		},
		HostedZoneId: hostedZoneId,
	}

	err := route53Service.ChangeResourceRecordSets(changeResourceRecordSetsInput)
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Error("Couldn't update DNS record")
	} else {
		logger.WithFields(logrus.Fields{"hostname": hostname}).Error("Couldn't update DNS record")
	}
}
