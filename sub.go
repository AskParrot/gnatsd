package main

import (
	"errors"
	"flag"
	//"fmt"
	"github.com/AskParrot/go-socket.io"
	"github.com/AskParrot/nats"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	//"reflect"
	"strings"
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

	//socket
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	// var showTime = flag.Bool("t", false, "Display timestamps")
	// var ssl = flag.Bool("ssl", false, "Use Secure Connection")

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
	//opts.Secure = *ssl
	opts.Token = "wC7sd4TMDCJvBvxayd9LsmU7eHpCWtk4rVbfHmM8xBtKBeVaq6ekzHrxRt5VF7aL6ZGpvELpAk3WYC9KuYF5LAFn7YrhRxaJC8S54cRS4ZBbP97Jn5Ze8f8ad9zsYNLfZ3ebEKSVrjMbzScsYxZFv7qgVwW2MbhjCFxHVYtb97Agp2JqpdrfJqqwaMgCEhv2tXguPpXuqqHKK4aT95EYhUr6mZ9Lw9VKUDqJYzd8aP6wb7yx7qtGpWwnCbMvPvDSa98fU45MexnqgLhr7q48JEhJ8ytG9BL5M8H7FCSrKLXNec9mLxBUMCHXYeWjWSvMaLxYvF65gqttH6hp3wupErAXG53YxcgcfA2EdxCxW2HSbFBnVSrQcdCBRaSXj6wFC4DMdsQh6LFnTeXh75L6xNsQD5CbuZmvNHAFgjBtFNheterceLvZCShmB4MAU9VVdZZQWRkg68U8HPtHbkmjJYLQLvJX5DxrDamHXrBgtFn7HDs8C39LFNCyTRHHREeJ"

	nc, err := opts.Connect()
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	server.On("connection", func(so socketio.Socket) {

		//userid := nil
		//fmt.Print(reflect.TypeOf(so))
		token := so.Request().FormValue("token")
		if len(token) < 1 {
			so.Close()
		} else {

			userData, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				if jwt.GetSigningMethod("HS256") != token.Method {
					return nil, errors.New("Invalid signing algorithm")
				}
				return []byte("thisisreallybigandsecurekey"), nil
			})
			if err != nil {
				so.Close()
			}
			userid := userData.Claims["id"]
			if str, ok := userid.(string); ok {
				/* act on str */
				nc.Subscribe(str, func(msg *nats.Msg) {
					so.Emit(str, msg)
				})
			} else {
				/* not string */
				panic("OOPS")
			}
		}

		so.On("join", func(room string) {

			nc.Subscribe(room, func(msg *nats.Msg) {
				so.Emit(room, msg)
			})
		})

		// log.Println("on connection")
		// so.Join("chat")
		// so.On("chat message", func(msg string) {
		//     log.Println("emit:", so.Emit("chat message", msg))
		//     so.BroadcastTo("chat", "chat message", msg)
		// })

		so.On("disconnection", func() {

			//      	nc.Unsubscribe(room, func(msg *nats.Msg) {
			// 	so.Emit(room, msg)
			// })
			log.Println("on disconnect")
		})

	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	//http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":9999", nil))

	//log.Printf("Listening on [%s]\n", subj)
	// if *showTime {
	// 	log.SetFlags(log.LstdFlags)
	// }

	//runtime.Goexit()
}
