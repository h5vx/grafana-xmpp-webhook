package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println(PackageName, PackageVersion, PackageCommitHash)

	configFilePath := flag.String("c", "config.toml", "Config file path")
	flag.Parse()

	err := ParseConfig(*configFilePath)

	if err != nil {
		log.Fatalf("Error while reading config: %v", err)
	}

	app := fiber.New(fiber.Config{Immutable: true})
	SetupHandlers(app)

	err = ConnectToXmpp()
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Fatal(app.Listen(GlobalConfig.Http.Listen_Address))
}
