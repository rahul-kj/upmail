package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rahul-kj/upmail/email"
	"github.com/sourcegraph/checkup"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = "upmail - %s\n"
	// VERSION is the binary version.
	VERSION = "v0.1.0"
)

var (
	configFile   = kingpin.Flag("config", "checkup.json config file location").Default("checkup.json").OverrideDefaultFromEnvar("CONFIG_FILE").String()
	recipient    = kingpin.Flag("recipient", "recipient for email notifications").OverrideDefaultFromEnvar("RECIPIENT_EMAIL").Required().String()
	interval     = kingpin.Flag("interval", "check interval (ex. 5ms, 10s, 1m, 3h)").Default("10m").OverrideDefaultFromEnvar("INTERVAL").String()
	smtpServer   = kingpin.Flag("smtp-server", "SMTP server for email notifications").OverrideDefaultFromEnvar("SMTP_SERVER").Required().String()
	smtpSender   = kingpin.Flag("sender", "SMTP default sender email address for email notifications").OverrideDefaultFromEnvar("SENDER_EMAIL").Required().String()
	smtpUsername = kingpin.Flag("smtp-username", "SMTP server username").OverrideDefaultFromEnvar("SMTP_USERNAME").String()
	smtpPassword = kingpin.Flag("smtp-password", "SMTP server password").OverrideDefaultFromEnvar("SMTP_PASSWORD").String()
	debug        = kingpin.Flag("debug", "run in debug mode").Default("false").OverrideDefaultFromEnvar("DEBUG").Bool()
)

func main() {
	kingpin.Version(VERSION)
	kingpin.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var ticker *time.Ticker
	// On ^C, or SIGTERM handle exit.
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		for sig := range s {
			ticker.Stop()
			logrus.Infof("Received %s, exiting.", sig.String())
			os.Exit(0)
		}
	}()

	configBytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	var c checkup.Checkup
	err = json.Unmarshal(configBytes, &c)
	if err != nil {
		log.Fatal(err)
	}

	n := email.Notifier{
		Recipient: *recipient,
		Server:    *smtpServer,
		Sender:    *smtpSender,
		Auth: smtp.PlainAuth(
			"",
			*smtpUsername,
			*smtpPassword,
			strings.SplitN(*smtpServer, ":", 2)[0],
		),
	}
	c.Notifier = n

	// parse the duration
	dur, err := time.ParseDuration(*interval)
	if err != nil {
		logrus.Fatalf("parsing %s as duration failed: %v", *interval, err)
	}

	logrus.Infof("Starting checks that will send emails to: %s", *recipient)
	ticker = time.NewTicker(dur)

	for range ticker.C {
		if _, err := c.Check(); err != nil {
			logrus.Warnf("check failed: %v", err)
		}
	}
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
