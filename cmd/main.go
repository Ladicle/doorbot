package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Ladicle/doorbot/pkg/sensor"
	"github.com/Ladicle/doorbot/pkg/slack"
	"github.com/spf13/pflag"
	"github.com/tarm/serial"
	"k8s.io/klog"
)

const (
	defaultTarget     = "/dev/tty.usbserial-MW4WB9CP"
	defaultBaud       = 115200
	defaultBufferSize = 128

	slackURLKey = "SLACK_URL"
)

var (
	lineEnd = []byte("\r\n")
)

type Options struct {
	Target string
	Baud   int

	SlackURL string
}

func init() {
	klog.InitFlags(flag.CommandLine)
	flag.Set("logtostderr", "true")
}

func main() {
	var opts Options

	flag.StringVar(&opts.Target, "target", defaultTarget, "The target serial device name")
	flag.IntVar(&opts.Baud, "baud", defaultBaud, "The baud rate for serial communication")
	flags := pflag.NewFlagSet("", pflag.ExitOnError)
	flags.Parse(os.Args)

	if err := opts.complete(); err != nil {
		log.Fatal(err)
	}
	log.Fatal(opts.run())
}

func (o *Options) complete() error {
	slackURL := os.Getenv(slackURLKey)
	if slackURL == "" {
		return fmt.Errorf("%s is required environment value", slackURLKey)
	}
	return nil
}

func (o *Options) run() error {
	c := &serial.Config{
		Name: o.Target,
		Baud: o.Baud,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	defer s.Close()

	log.Println("start to read serial data...")

	var (
		line bytes.Buffer
		buf  = make([]byte, defaultBufferSize)
	)
	for {
		n, err := s.Read(buf)
		if err != nil {
			return err
		}
		line.Write(buf[:n])

		if !isLineEnd(line.Bytes()) {
			continue
		}
		data := line.String()
		line.Reset()

		s, err := sensor.Parse(data)
		if err != nil {
			log.Printf("ERROR: failed to parse data: %v", err)
		}
		switch v := s.(type) {
		case sensor.CloserSensor:
			if v.State == sensor.NoMagnet {
				slack.SendMsg(o.SlackURL, "the door opened")
			}
		default:
			log.Printf("ERROR: %v is unknown type", v)
		}
	}
}

func isLineEnd(data []byte) bool {
	return bytes.HasSuffix(data, lineEnd)
}
