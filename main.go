package main

import (
	"insighta-cli/cmd"
	_ "insighta-cli/cmd/profiles" // blank import triggers init()

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cmd.Execute()
}
