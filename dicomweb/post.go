package dicomweb

import (
	"bytes"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"

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

func PostMultipart(url string, data *[]byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(*data))
	if err != nil {
		return &http.Response{}, err
	}
	for key, element := range headers {
		req.Header.Set(key, element)
	}
	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	if r.StatusCode != http.StatusOK {
		return r, errors.New(fmt.Sprintf("HTTP Status: %d", r.StatusCode))
	}
	return r, nil
}

func Stow(url string, dcm_slice []*dicom.Dataset, headers map[string]string) (*http.Response, error) {
	b, content_type, err := WriteMultipart(dcm_slice)
	if err != nil {
		return &http.Response{}, err
	}
	headers["Content-Type"] = content_type
	return PostMultipart(url, b, headers)
}
