package dicomweb

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/suyashkumar/dicom"
)

func GetMultipart(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, err
	}
	for key, element := range headers {
		req.Header.Set(key, element)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return &http.Response{}, errors.New(fmt.Sprintf("HTTP Status: %d", resp.StatusCode))
	}
	return resp, nil
}

func ReadMultipart(resp *http.Response) ([]*dicom.Dataset, []byte, error) {
	if resp.StatusCode != http.StatusOK {
		return []*dicom.Dataset{}, []byte{}, errors.New(fmt.Sprintf("Status: %d", resp.StatusCode))
	}
	contentType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return []*dicom.Dataset{}, []byte{}, err
	}
	if contentType != "multipart/related" {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return []*dicom.Dataset{}, []byte{}, err
		}
		return []*dicom.Dataset{}, b, nil
	}
	multipartReader := multipart.NewReader(resp.Body, params["boundary"])
	res := []*dicom.Dataset{}
	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return []*dicom.Dataset{}, []byte{}, err
		}
		data, err := ioutil.ReadAll(part)
		if err != nil {
			return []*dicom.Dataset{}, []byte{}, err
		}
		dcm, err := dicomutil.Bytes2Dicom(data)
		if err != nil {
			return []*dicom.Dataset{}, []byte{}, err
		}
		res = append(res, dcm)
	}
	return res, []byte{}, nil
}

func Get(url string, headers map[string]string) ([]*dicom.Dataset, []byte, error) {
	resp, err := GetMultipart(url, headers)
	if err != nil {
		return []*dicom.Dataset{}, []byte{}, err
	}
	defer resp.Body.Close()
	res, b, err := ReadMultipart(resp)
	return res, b, err
}
