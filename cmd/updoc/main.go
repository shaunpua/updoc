package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shaunpua/updoc/internal/confluence"
)

func main() {
	// Load variables from .env if the file exists
	_ = godotenv.Load()

	// Show the env values so you can confirm they’re set
	base := os.Getenv("CONF_BASE")
	email := os.Getenv("CONF_EMAIL")
	// token := os.Getenv("CONF_TOKEN")
	fmt.Println("Env check:")
	fmt.Println("  CONF_BASE :", base)
	fmt.Println("  CONF_EMAIL:", email)
	// fmt.Println("  CONF_TOKEN:", strings.Repeat("*", len(token)))

	// Hard-coded page ID for now
	space := confluence.NewSpace("622593")

	// Optional flag to overwrite the body
	newBody := flag.String("body", "", "Replace page body with this HTML")
	flag.Parse()

	// ---- Read page ----
	page, err := confluence.GetPage(space.Client, space.PageID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n=== Title: %s ===\n\n", page.Title)
	fmt.Println(page.Body.Storage.Value)

	fmt.Println("\n=== Page properties ===")
	if len(page.Metadata.Properties) == 0 {
		fmt.Println("(none)")
	}
	for k, v := range page.Metadata.Properties {
		fmt.Printf("%s = %v\n", k, v.Value)
	}

	// ---- Write page (if -body supplied) ----
	if *newBody != "" {
		if err := confluence.UpdateBody(space.Client, page, *newBody); err != nil {
			log.Fatal(err)
		}
		fmt.Println("\nBody updated successfully ✔")
	}
}
