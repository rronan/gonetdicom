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

func MakePostRequest(url string, data *[]byte, headers map[string]string) error {
	req, err := http.NewRequest("POST", url, bytes.NewReader(*data))
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
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("HTTP Status: %d", resp.StatusCode))
	}
	return nil
}

func PostDicomWeb(url string, dcm_slice []*dicom.Dataset, headers map[string]string) error {
	b, content_type, err := WriteMultipart(dcm_slice)
	if err != nil {
		return err
	}
	headers["Content-Type"] = content_type
	err = MakePostRequest(url, b, headers)
	return err
}
