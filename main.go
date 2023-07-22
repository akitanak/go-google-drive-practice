package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/akitanak/go-google-drive-practice/commons"
	drivev3 "google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
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
	if _, err := service.Files.Create(&drivev3.File{Name: commons.FileName}).Context(ctx).Media(f, googleapi.ContentType("text/csv")).Do(); err != nil {
		panic(err)
	}

	// list files
	files, err := service.Files.List().Context(ctx).PageSize(1).Do()
	if err != nil {
		panic(err)
	}
	if len(files.Files) == 0 {
		panic("no files found")
	}
	listed := files.Files[0]
	fmt.Printf("%s (%s) url: %s%s\n", listed.Name, listed.Id, commons.DownloadURLPrefix, listed.Id)

	// download file
	res, err := service.Files.Get(listed.Id).Context(ctx).Download()
	if err != nil {
		panic(err)
	}
	body := res.Body
	defer body.Close()

	buf := make([]byte, commons.Bufsize)
	n, err := body.Read(buf)
	if err != nil {
		panic(err)
	}
	if commons.Bufsize < n {
		panic("bufsize is too small")
	}

	if err := ioutil.WriteFile(commons.FileName, buf[:n], 0644); err != nil {
		panic(err)
	}
}
