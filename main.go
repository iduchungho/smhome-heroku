/*
main.go file
Author: Ho Duc Hung
Start api: cd /src -> make run or go run main.go app.go
*/
package main

import appli "smhome/app"

func main() {
	// create application
	app := appli.GetApplication()
	// app run localhost:8080
	app.Run()
}
