package main

import (
	"flag"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func usage() {
	log.Printf("Usage: nats-consumers [-s server] [-creds file] [-nkey file] [-tlscert file] [-tlskey file] [-tlscacert file] stream-name subjects\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var nkeyFile = flag.String("nkey", "", "NKey Seed File")
	var tlsClientCert = flag.String("tlscert", "", "TLS client certificate file")
	var tlsClientKey = flag.String("tlskey", "", "Private key file for client certificate")
	var tlsCACert = flag.String("tlscacert", "", "CA certificate to verify peer against")
	var method = flag.String("method", "", "Method: add/update/delete on the consumer")
	var replicas = flag.Int("replicas", 1, "Replicas for stream")
	var showHelp = flag.Bool("h", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	args := flag.Args()
	if len(args) <= 0 {
		showUsageAndExit(1)
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Publisher")}

	if *userCreds != "" && *nkeyFile != "" {
		log.Fatal("specify -seed or -creds")
	}

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}

	// Use TLS client authentication
	if *tlsClientCert != "" && *tlsClientKey != "" {
		opts = append(opts, nats.ClientCert(*tlsClientCert, *tlsClientKey))
	}

	// Use specific CA certificate
	if *tlsCACert != "" {
		opts = append(opts, nats.RootCAs(*tlsCACert))
	}

	// Use Nkey authentication.
	if *nkeyFile != "" {
		opt, err := nats.NkeyOptionFromSeed(*nkeyFile)
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, opt)
	}

	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	stream := args[0]

	var subjects string
	if len(args) > 1 {
		subjects = args[1]
	}

	switch *method {
	case "add":
		js.AddStream(&nats.StreamConfig{
			Name:     stream,
			Subjects: []string{subjects},
			Replicas: *replicas,
		})
	case "update":
		js.UpdateStream(&nats.StreamConfig{
			Name:     stream,
			Replicas: *replicas,
		})
	case "delete":
		js.DeleteStream(stream)
	default:
		showUsageAndExit(1)
	}

}
