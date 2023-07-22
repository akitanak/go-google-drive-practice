package commons

import (
	"context"
	"fmt"

	drivev3 "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	AssetsDir         = "./assets"
	FileName          = "sample.csv"
	DownloadURLPrefix = "https://drive.google.com/file/d/"
	Bufsize           = 1024
)

func GetService(ctx context.Context) (*drivev3.Service, error) {
	service, err := drivev3.NewService(ctx, option.WithScopes(drivev3.DriveScope))
	if err != nil {
		return nil, fmt.Errorf("unable to create Drive service: %w", err)
	}
	return service, nil
}
