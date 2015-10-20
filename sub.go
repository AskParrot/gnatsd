// Copyright 2012-2015 Apcera Inc. All rights reserved.
// +build ignore

package main

import (
	"flag"
	"log"
	"runtime"
	"strings"

	"github.com/AskParrot/nats"
)

func usage() {
	log.Fatalf("Usage: nats-sub [-s server] [--ssl] [-t] <subject> \n")
}

var index = 0

func printMsg(m *nats.Msg, i int) {
	index += 1
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display timestamps")
	var ssl = flag.Bool("ssl", false, "Use Secure Connection")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	opts := nats.DefaultOptions
	opts.Servers = strings.Split(*urls, ",")
	for i, s := range opts.Servers {
		opts.Servers[i] = strings.Trim(s, " ")
	}
	opts.Secure = *ssl
	opts.Token = "wC7sd4TMDCJvBvxayd9LsmU7eHpCWtk4rVbfHmM8xBtKBeVaq6ekzHrxRt5VF7aL6ZGpvELpAk3WYC9KuYF5LAFn7YrhRxaJC8S54cRS4ZBbP97Jn5Ze8f8ad9zsYNLfZ3ebEKSVrjMbzScsYxZFv7qgVwW2MbhjCFxHVYtb97Agp2JqpdrfJqqwaMgCEhv2tXguPpXuqqHKK4aT95EYhUr6mZ9Lw9VKUDqJYzd8aP6wb7yx7qtGpWwnCbMvPvDSa98fU45MexnqgLhr7q48JEhJ8ytG9BL5M8H7FCSrKLXNec9mLxBUMCHXYeWjWSvMaLxYvF65gqttH6hp3wupErAXG53YxcgcfA2EdxCxW2HSbFBnVSrQcdCBRaSXj6wFC4DMdsQh6LFnTeXh75L6xNsQD5CbuZmvNHAFgjBtFNheterceLvZCShmB4MAU9VVdZZQWRkg68U8HPtHbkmjJYLQLvJX5DxrDamHXrBgtFn7HDs8C39LFNCyTRHHREeJ"

	nc, err := opts.Connect()
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	subj, i := args[0], 0

	nc.Subscribe(subj, func(msg *nats.Msg) {
		i += 1
		printMsg(msg, i)
	})

	log.Printf("Listening on [%s]\n", subj)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	runtime.Goexit()
}