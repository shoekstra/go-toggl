//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	toggl "github.com/shoekstra/go-toggl"
)

func main() {
	client, err := toggl.NewClient(os.Getenv("TOGGL_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	me, _, err := client.Me.GetMe(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID:                %d\n", me.ID)
	fmt.Printf("Name:              %s\n", me.Fullname)
	fmt.Printf("Email:             %s\n", me.Email)
	fmt.Printf("Timezone:          %s\n", me.Timezone)
	fmt.Printf("Default workspace: %d\n", me.DefaultWorkspaceID)
}
