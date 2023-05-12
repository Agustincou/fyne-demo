package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	GoldPriceOrgChartURL = "https://goldprice.org/charts/gold_3d_b_o_%s_x.png"
	DownloadedFileName = "gold_3d_b_o_x.png"
)

type ImageDownloader interface {
	Download() error
}

type HTTPImageDownloader struct {
	ImageDownloader
	_URL     string
	fileName string
	client   *http.Client
}

func NewHTTPImageDownloader(client *http.Client, baseURL string, fileName string) ImageDownloader {
	return &HTTPImageDownloader{
		_URL:     baseURL,
		client:   client,
		fileName: fileName,
	}
}

func (h *HTTPImageDownloader) Download() error {
	response, err := h.client.Get(fmt.Sprintf(GoldPriceOrgChartURL, Currency))
	if err != nil {
		log.Printf("error requesting %s\n", h._URL)

		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("received wrong response code when downloading image")
	}

	readedBytes, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		log.Printf("error reading rersponse bytes %s\n", h._URL)

		return err
	}

	img, _, imgErr := image.Decode(bytes.NewReader(readedBytes))
	if imgErr != nil {
		return imgErr
	}

	out, outErr := os.Create(fmt.Sprintf("./%s", h.fileName))
	if outErr != nil {
		return outErr
	}

	if err := png.Encode(out, img); err != nil {
		return err
	}

	return nil
}
