package dicomweb

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/suyashkumar/dicom"
)

func Get(url string, headers map[string]string) (*http.Response, error) {
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
		return &http.Response{}, &RequestError{StatusCode: resp.StatusCode, Err: errors.New(resp.Status)}
	}
	return resp, nil
}

func ReadMultipart(resp *http.Response) ([]*dicom.Dataset, []byte, error) {
	res := []*dicom.Dataset{}
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return res, []byte{}, err
	}
	multipartReader := multipart.NewReader(resp.Body, params["boundary"])
	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return res, []byte{}, err
		}
		data, err := io.ReadAll(part)
		if err != nil {
			return res, []byte{}, err
		}
		contentType := part.Header.Get("Content-type")
		if contentType == "application/json" {
			return res, data, err
		}
		if contentType != "application/dicom" {
			return res, []byte{}, &RequestError{StatusCode: 415, Err: errors.New("Invalid Content-Type:" + contentType)}
		}
		dcm, err := dicomutil.Bytes2Dicom(data)
		if err != nil {
			return res, []byte{}, err
		}
		res = append(res, dcm)
	}
	return res, []byte{}, nil
}

func Wado(url string, headers map[string]string) ([]*dicom.Dataset, []byte, error) {
	resp, err := Get(url, headers)
	if err != nil {
		return []*dicom.Dataset{}, []byte{}, err
	}
	if resp.StatusCode >= 400 {
		return []*dicom.Dataset{}, []byte{}, &RequestError{StatusCode: resp.StatusCode, Err: errors.New(resp.Status)}
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType == "application/json" {
		data, err := io.ReadAll(resp.Body)
		return []*dicom.Dataset{}, data, err
	}
	if !strings.HasPrefix(contentType, "multipart/related") {
		return []*dicom.Dataset{}, []byte{}, &RequestError{StatusCode: 415, Err: errors.New("Invalid Content-Type:" + contentType)}
	}
	dcm_slice, byte_slice, err := ReadMultipart(resp)
	return dcm_slice, byte_slice, err
}

func ReadMultipartToFile(resp *http.Response, folder string) ([]string, []byte, error) {
	res := []string{}
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return res, []byte{}, err
	}
	multipartReader := multipart.NewReader(resp.Body, params["boundary"])
	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return res, []byte{}, err
		}
		contentType := part.Header.Get("Content-type")
		if contentType == "application/json" {
			data, err := io.ReadAll(part)
			return res, data, err
		}
		if contentType != "application/dicom" {
			return res, []byte{}, &RequestError{StatusCode: 415, Err: errors.New("Invalid Content-Type:" + contentType)}
		}
		dcm_path := fmt.Sprintf("%s/%s", folder, dicomutil.RandomDicomName())
		f, err := os.Create(dcm_path)
		if err != nil {
			return res, []byte{}, err
		}
		defer f.Close()
		io.Copy(f, part)
		res = append(res, dcm_path)
	}
	return res, []byte{}, nil
}

func WadoToFile(url string, headers map[string]string, folder string) ([]string, []byte, error) {
	resp, err := Get(url, headers)
	if err != nil {
		return []string{}, []byte{}, err
	}
	if resp.StatusCode >= 400 {
		return []string{}, []byte{}, &RequestError{StatusCode: resp.StatusCode, Err: errors.New(resp.Status)}
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType == "application/json" {
		data, err := io.ReadAll(resp.Body)
		return []string{}, data, err
	}
	if !strings.HasPrefix(contentType, "multipart/related") {
		return []string{}, []byte{}, &RequestError{StatusCode: 415, Err: errors.New("Invalid Content-Type:" + contentType)}
	}
	dcm_path_slice, byte_slice, err := ReadMultipartToFile(resp, folder)
	return dcm_path_slice, byte_slice, err
}
