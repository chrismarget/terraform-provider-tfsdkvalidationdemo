package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
	"terraform-provider-tfsdkvalidationdemo/tfsdkvalidationdemo"
)

func main() {
	err := providerserver.Serve(context.Background(), tfsdkvalidationdemo.New, providerserver.ServeOpts{
		Address: "example.com/chrismarget/tfsdkvalidationdemo",
	})
	if err != nil {
		log.Fatal(err)
	}
}
