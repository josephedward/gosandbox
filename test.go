package main

import(
	"testing"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/joho/godotenv"
	"golang.design/x/clipboard"
	// "io"
	"log"
	"os"
	// "reflect"
	"time"
	"fmt"
)


func LoginTest(){
	
	// Create a new page
	url := "https://learn.acloud.guru/cloud-playground/cloud-sandboxes"
	fmt.Println(url)

	//load env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Could not load .env file - Err: %s", err)
	}
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	//call the login function
	Login(url, username, password)
	
}
