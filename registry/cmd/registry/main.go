package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	libreg "github.com/RackHD/neighborhood-manager/libreg/registry"
	"github.com/RackHD/neighborhood-manager/libreg/registry/consul"
	"github.com/RackHD/neighborhood-manager/registry"
	"github.com/spf13/viper"
)

var binaryName, buildDate, buildUser, commitHash, goVersion, osArch, releaseVersion string
var backendAddr, datacenter string

func init() {
	flag.StringVar(&backendAddr, "address", "127.0.0.1:8500", "Address:port of the backend store")
	flag.StringVar(&datacenter, "datacenter", "dc1", "Datacenter name of the backend storage")
}

func main() {
	log.Println(binaryName)
	log.Println("  Release version: " + releaseVersion)
	log.Println("  Built On: " + buildDate)
	log.Println("  Build By: " + buildUser)
	log.Println("  Commit Hash: " + commitHash)
	log.Println("  Go version: " + goVersion)
	log.Println("  OS/Arch: " + osArch)

	flag.Parse()

	consul.Register()
	registry, err := registry.NewRegistry(libreg.CONSUL, datacenter, backendAddr)
	if err != nil {
		log.Fatalf("Error creating new Service Registry: %s", err)
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("Service Registry Configuration Error: %s\n", err)
	}

	// Set viper configurations
	viper.SetConfigName("registry")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/infranm.d")
	viper.AddConfigPath(dir)
	viper.AddConfigPath("$GOPATH/bin")
	viper.AddConfigPath("$HOME")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Service Registry Configuration Error: %s\n", err)
	}

	log.Printf("Service Registry Configuration: %s\n", viper.ConfigFileUsed())

	// Parse all urns in config file and add to whitelist of search terms
	tags := viper.GetStringMapString("ssdp.tags")
	for tag := range tags {
		urns := viper.GetStringSlice("ssdp.tags." + tag)
		for _, urn := range urns {
			registry.AddSearchTerm(urn, tag)
		}
	}

	registry.Run()

	return
}
