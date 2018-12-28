package main

import (
	"fmt"
	"github.com/ecray/avdb/app"
	"github.com/urfave/cli"
	"log"
	"net"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "avdb"
	app.Usage = "Ansible Variables Database"
	app.Version = "0.1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Eric Raymond",
			Email: "ecraymond@gmail.com",
		},
	}
	app.Action = server
	app.Flags = []cli.Flag{
		cli.StringFlag{
			EnvVar: "AVDB_ADDR",
			Name:   "addr",
			Usage:  "IP to bind",
			Value:  "127.0.0.1",
		},
		cli.StringFlag{
			EnvVar: "AVDB_PORT",
			Name:   "port",
			Usage:  "Port to bind socket",
			Value:  "3000",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func server(c *cli.Context) {
	conn := fmt.Sprintf("%s:%s", c.String("addr"), c.String("port"))

	// check socket in use
	_, err := net.Listen("tcp", conn)
	if err != nil {
		log.Fatal("Port already in use... SCRAM!!")
	}

	a := &app.App{}
	a.Initialize()
	log.Printf("Running on %s", conn)
	a.Run(conn)
}
