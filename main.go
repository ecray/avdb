package main

import (
	"github.marqeta.com/ecray/avdb/app"
)

func main() {
	app := &app.App{}
	app.Initialize()
	app.Run(":3000")
}

/*
var flags = []cli.Flag{
    cli.StringFlag{
        EnvVar: "AVDB_ADDR",
        Name:   "addr",
        Usage:  "IP to bind",
        Value:  "0.0.0.0",
    },
    cli.StringFlag{
        EnvVar: "AVDB_PORT",
        Name:   "port",
        Usage:  "Port to bind socket",
        Value:  "8080",
    }


func server(c *cli.Context) {
    addr := c.String("addr")
    port := c.String("port")

    // check socket in use
    conn, err := net.Listen("tcp", fmt.Sprintf("%s:%s", addr, port))
    if err != nil {
        log.Fatal("Port already in use... SCRAM!!")
}
*/
