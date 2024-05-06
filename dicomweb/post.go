package dicomweb

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"time"

	"github.com/suyashkumar/dicom"
)

func WriteMultipart(dcm_slice []*dicom.Dataset) (*[]byte, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	for _, dcm := range dcm_slice {
		part, err := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"application/dicom"}})
		if err != nil {
			return &[]byte{}, "", err
		}
		err = dicom.Write(part, *dcm, dicom.DefaultMissingTransferSyntax(), dicom.SkipVRVerification(), dicom.SkipValueTypeVerification())
		if err != nil {
			return &[]byte{}, "", err
		}
	}
	params := make(map[string]string)
	params["boundary"] = writer.Boundary()
	params["type"] = "application/dicom"
	content_type := mime.FormatMediaType("multipart/related", params)
	b := buf.Bytes()
	return &b, content_type, nil
}

func PostMultipart(url string, data *[]byte, headers map[string]string, timeout int) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(*data))
	if err != nil {
		return &http.Response{}, err
	}
	for key, element := range headers {
		req.Header.Set(key, element)
	}
	client := &http.Client{Timeout: time.Duration(timeout * 1e9)}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return &http.Response{}, &RequestError{StatusCode: resp.StatusCode, Err: errors.New(resp.Status)}
	}
	return resp, nil
}

func Stow(url string, dcm_slice []*dicom.Dataset, headers map[string]string, timeout int) (*http.Response, error) {
	b, content_type, err := WriteMultipart(dcm_slice)
	if err != nil {
		return &http.Response{}, err
	}
	headers["Content-Type"] = content_type
	return PostMultipart(url, b, headers, timeout)
}

func WriteMultipartFromFile(dcm_path_slice []string) (*[]byte, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	for _, dcm_path := range dcm_path_slice {

		data, err := os.Open(dcm_path)
		if err != nil {
			return &[]byte{}, "", err
		}
		defer data.Close()

		part, err := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"application/dicom"}})
		if err != nil {
			return &[]byte{}, "", err
		}
		_, err = io.Copy(part, data)
		if err != nil {
			return &[]byte{}, "", err
		}
	}
	params := make(map[string]string)
	params["boundary"] = writer.Boundary()
	params["type"] = "application/dicom"
	content_type := mime.FormatMediaType("multipart/related", params)
	b := buf.Bytes()
	return &b, content_type, nil
}

func StowFromFile(url string, dcm_path_slice []string, headers map[string]string, timeout int) (*http.Response, error) {
	b, content_type, err := WriteMultipartFromFile(dcm_path_slice)
	if err != nil {
		return &http.Response{}, err
	}
	headers["Content-Type"] = content_type
	return PostMultipart(url, b, headers, timeout)
}
