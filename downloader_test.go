package theseus

import (
	"context"
	"net/http"
	"os"
	"testing"

	mw "github.com/njayp/middleware/client"
	"github.com/njayp/middleware/client/limiter"
)

func TestDownload(t *testing.T) {
	ctx := context.Background()
	client := &http.Client{
		Transport: mw.BuildTransport(limiter.WithLimiter(limiter.WithCount(1))),
	}
	url := "https://download.samplelib.com/mp4/sample-30s.mp4"
	filename := "test_video.mp4"

	err := Download(ctx, client, url, filename)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Downloaded file does not exist: %v", err)
	}

	// Clean up the downloaded file
	//os.Remove(filename)
}
