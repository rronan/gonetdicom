package dicomweb

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/suyashkumar/dicom"
)

func Put(url string, dcm *dicom.Dataset, headers map[string]string) error {
	data := dicomutil.Dicom2Bytes(dcm)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(*data))
	if err != nil {
		return err
	}
	for key, element := range headers {
		req.Header.Set(key, element)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Status: %d", resp.StatusCode))
	}
	return err
}

func PutFromFile(url string, dcm_path string, headers map[string]string) error {
	data, err := os.Open(dcm_path)
	if err != nil {
		return err
	}
	defer data.Close()
	req, err := http.NewRequest("PUT", url, data)
	if err != nil {
		return err
	}
	for key, element := range headers {
		req.Header.Set(key, element)
	}
	stat, err := data.Stat()
	if err != nil {
		return err
	}
	req.ContentLength = stat.Size()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Status: %d", resp.StatusCode))
	}
	return err
}
