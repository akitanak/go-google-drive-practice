package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	drivev3 "google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	assetsDir         = "./assets"
	fileName          = "sample.csv"
	downloadURLPrefix = "https://drive.google.com/file/d/"
	bufsize           = 1024
)

func getService(ctx context.Context) (*drivev3.Service, error) {
	service, err := drivev3.NewService(ctx, option.WithScopes(drivev3.DriveScope))
	if err != nil {
		return nil, fmt.Errorf("unable to create Drive service: %w", err)
	}
	return service, nil
}

func main() {
	ctx := context.Background()
	service, err := getService(ctx)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(strings.Join([]string{assetsDir, fileName}, "/"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// upload file
	if _, err := service.Files.Create(&drivev3.File{Name: fileName}).Context(ctx).Media(f, googleapi.ContentType("text/csv")).Do(); err != nil {
		panic(err)
	}

	// list files
	files, err := service.Files.List().PageSize(1).Do()
	if err != nil {
		panic(err)
	}
	if len(files.Files) == 0 {
		panic("no files found")
	}
	listed := files.Files[0]
	fmt.Printf("%s (%s) url: %s%s\n", listed.Name, listed.Id, downloadURLPrefix, listed.Id)

	// download file
	res, err := service.Files.Get(listed.Id).Context(ctx).Download()
	if err != nil {
		panic(err)
	}
	body := res.Body
	defer body.Close()

	buf := make([]byte, bufsize)
	n, err := body.Read(buf)
	if err != nil {
		panic(err)
	}
	if bufsize < n {
		panic("bufsize is too small")
	}

	if err := ioutil.WriteFile(fileName, buf[:n], 0644); err != nil {
		panic(err)
	}
}
