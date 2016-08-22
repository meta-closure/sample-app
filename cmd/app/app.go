package main

import (
	"os"
	"sample-app/app"

	log "github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	Port string `short:"p" long:"port" description:"A Port, default :8080"`
}

func main() {
	os.Exit(cli())
}

func cli() int {
	o := &opts{}
	_, err := flags.Parse(o)
	if err != nil {
		log.Printf("%s", err)
		return 1
	}

	if o.Port != "" {
		server(":8080")
	} else {
		server(":" + o.Port)
	}
	return 0
}

func server(p string) {
	if p != "" {
		log.Infof("Server start running with Port%s", p)
	} else {
		log.Info("Server start running with Port:8080")
	}
	app.Run(p)
}
