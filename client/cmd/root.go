/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	bookpb "github.com/jy0t1/grpc-mongo-crud/proto"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cfgFile string
// client and context global variables
var client bookpb.BookServiceClient
var requestCtx context.Context
var requestOpts grpc.DialOption


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mongo_connect",
	Short: "Communicating to MongoDB via gRPC client",
	Long: `A gRPC client to connect to MongoDB to create, read, update, delete operation on collection.
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("root.go error ... 0001")
		fmt.Println(err)
		os.Exit(1)
	}
}


func init() {
	cobra.OnInitialize(initConfig)
    fmt.Println("Starting Book Service Client in cobra/cmd/root.go ")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mongo_connect.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

    // After Cobra root config init
	// We initialize the client
	fmt.Println("Starting Book Service Client in init() in root.go ....................")
	// Establish context to timeout if server does not respond
	requestCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	// Establish insecure grpc options (no TLS)
	requestOpts = grpc.WithInsecure()
	// Dial the server, returns a client connection
	conn, err := grpc.Dial("localhost:50051", requestOpts)
	if err != nil {
		log.Fatalf("Unable to establish client connection to localhost:50051: %v", err)
	}

	// defer posptones the execution of a function until the surrounding function returns
	// conn.Close() will not be called until the end of main()
	// The arguments are evaluated immeadiatly but not executed
	// defer conn.Close()

	// Instantiate the BookServiceClient with our client connection to the server
	client = bookpb.NewBookServiceClient(conn)
	fmt.Println("After bookpb.NewBookServiceClient in cmd/root.go ....................")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fmt.Println("Starting in initConfig() in cmd/root.go ....................")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mongo_connect" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mongo_connect")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	fmt.Println("Ending in initConfig() in cmd/root.go ....................")
}
