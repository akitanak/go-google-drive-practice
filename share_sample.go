package main

import (
	"context"
	"os"
	"strings"

	"github.com/akitanak/go-google-drive-practice/commons"
	drivev3 "google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

var (
	shareTestFolder = os.Getenv("SHARE_TEST_FOLDER")
	shareEmail      = os.Getenv("SHARE_EMAIL")
)

func main() {
	ctx := context.Background()
	service, err := commons.GetService(ctx)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(strings.Join([]string{commons.AssetsDir, commons.FileName}, "/"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// upload file
	file := &drivev3.File{
		Name: commons.FileName,
		Parents: []string{
			shareTestFolder,
		},
	}

	created, err := service.Files.Create(file).Context(ctx).Media(f, googleapi.ContentType("text/csv")).Do()
	if err != nil {
		panic(err)
	}

	permission := &drivev3.Permission{
		Type:         "user",
		Role:         "reader",
		EmailAddress: shareEmail,
	}

	service.Permissions.Create(created.Id, permission).Context(ctx).Do()
}
