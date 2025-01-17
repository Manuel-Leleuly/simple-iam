package main

import (
	"github.com/Manuel-Leleuly/simple-iam/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {

}
