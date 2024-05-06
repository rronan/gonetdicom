package dicomweb

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/suyashkumar/dicom"
)

func Put(url string, dcm *dicom.Dataset, headers map[string]string, timeout int) error {
	data := dicomutil.Dicom2Bytes(dcm)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(*data))
	if err != nil {
		return err
	}
	for key, element := range headers {
		req.Header.Set(key, element)
	}
	client := &http.Client{Timeout: time.Duration(timeout * 1e9)}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &RequestError{StatusCode: resp.StatusCode, Err: errors.New(resp.Status)}
	}
	return nil
}

func PutFromFile(url string, dcm_path string, headers map[string]string, timeout int) error {
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
	client := &http.Client{Timeout: time.Duration(timeout * 1e9)}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return &RequestError{StatusCode: resp.StatusCode, Err: errors.New(resp.Status)}
	}
	return nil
}
