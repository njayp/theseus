package fdl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Part struct {
	request  *http.Request
	fileName string
}

func NewPart(ctx context.Context, url string, start, end int, fileName string) (*Part, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
	req.Header.Set("Range", rangeHeader)

	return &Part{
		request:  req,
		fileName: fileName,
	}, nil
}

func (p *Part) downloadRequest(client *http.Client) error {
	resp, err := client.Do(p.request)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	out, err := os.Create(p.fileName)
	if err != nil {
		return fmt.Errorf("file creation error: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	return nil
}

// mergePart merges one part file into the final output file
func (p *Part) mergePart(w io.Writer) error {
	partFile, err := os.Open(p.fileName)
	if err != nil {
		return fmt.Errorf("error opening part file %s: %w", p.fileName, err)
	}
	defer partFile.Close()

	_, err = io.Copy(w, partFile)
	if err != nil {
		return fmt.Errorf("error merging part file %s: %w", p.fileName, err)
	}

	return p.removeFile()
}

// removePartFile removes the part file after merging
func (p *Part) removeFile() error {
	err := os.Remove(p.fileName)
	if err != nil {
		return fmt.Errorf("error removing part file %s: %w", p.fileName, err)
	}
	return nil
}
