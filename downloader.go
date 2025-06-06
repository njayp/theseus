package fdl

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const partSize = 10 * 1024 * 1024 // 10 MB

func Download(ctx context.Context, client *http.Client, url string, filename string) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	size, err := getFileSize(client, url)
	if err != nil {
		return err
	}
	slog.Info("File size", "size", size, "url", url)

	// Make array of parts based on the file size
	parts := []*Part{}
	for i := 0; i < size; i += partSize {
		partName := filename + ".part" + strconv.Itoa(i/partSize)
		part, err := NewPart(ctx, url, i, i+partSize-1, partName)
		if err != nil {
			return err
		}

		parts = append(parts, part)
	}
	slog.Info("Number of parts", "count", len(parts), "url", url)

	// Download parts concurrently
	wg := &sync.WaitGroup{}
	for _, part := range parts {
		wg.Add(1)

		go func(p *Part) {
			defer wg.Done()

			err := p.downloadRequest(client)
			if err != nil {
				cancel(err) // Cancel the context if any part fails
			}
		}(part)
	}
	wg.Wait()
	err = context.Cause(ctx)
	if err != nil {
		return err // Return the error if any part failed
	}
	slog.Info("All parts downloaded", "url", url)

	// Merge parts into the final file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, part := range parts {
		err := part.mergePart(file)
		if err != nil {
			slog.Error("Failed to merge part", "part", part.fileName, "error", err)
			return err
		}
	}
	slog.Info("All parts merged", "filename", filename)
	return nil
}

// getFileSize fetches the file size from Content-Length header
func getFileSize(client *http.Client, url string) (int, error) {
	resp, err := client.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return 0, err
	}

	return size, nil
}
