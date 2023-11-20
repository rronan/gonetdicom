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

	"github.com/labstack/gommon/log"
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

func ReadMultipart(resp *http.Response) ([]*dicom.Dataset, error) {
	contentType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return []*dicom.Dataset{}, err
	}
	if contentType != "multipart/related" {
		return []*dicom.Dataset{}, &RequestError{StatusCode: 415, Err: errors.New("Invalid Content-Type:" + contentType)}
	}
	multipartReader := multipart.NewReader(resp.Body, params["boundary"])
	res := []*dicom.Dataset{}
	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return []*dicom.Dataset{}, err
		}
		if part.Header.Get("Content-type") != "application/dicom" {
			data, err := io.ReadAll(part)
			if err != nil {
				return []*dicom.Dataset{}, err
			}
			msg := fmt.Sprintf(
				"ReadMultipartToFile() received unknown Content-type: %s\n %s\n passing...",
				part.Header.Get("Content-type"),
				string(data),
			)
			log.Warn(msg)
		}
		data, err := io.ReadAll(part)
		if err != nil {
			return []*dicom.Dataset{}, err
		}
		dcm, err := dicomutil.Bytes2Dicom(data)
		if err != nil {
			return []*dicom.Dataset{}, err
		}
		res = append(res, dcm)
	}
	return res, nil
}

func Wado(url string, headers map[string]string) ([]*dicom.Dataset, *http.Response, error) {
	resp, err := Get(url, headers)
	if err != nil {
		return []*dicom.Dataset{}, resp, err
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "multipart/related") {
		return []*dicom.Dataset{}, resp, nil
	}
	dcm_slice, err := ReadMultipart(resp)
	return dcm_slice, resp, err
}

func ReadMultipartToFile(resp *http.Response, folder string) ([]string, error) {
	res := []string{}
	contentType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return res, err
	}
	if contentType != "multipart/related" {
		return res, &RequestError{StatusCode: 415, Err: errors.New("Invalid Content-Type:" + contentType)}
	}
	multipartReader := multipart.NewReader(resp.Body, params["boundary"])
	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return res, err
		}
		if part.Header.Get("Content-type") != "application/dicom" {
			data, err := io.ReadAll(part)
			if err != nil {
				return []string{}, err
			}
			msg := fmt.Sprintf(
				"ReadMultipartToFile() received unknown Content-type: %s\n %s\n passing...",
				part.Header.Get("Content-type"),
				string(data),
			)
			log.Warn(msg)
		}
		dcm_path := fmt.Sprintf("%s/%s", folder, dicomutil.RandomDicomName())
		f, err := os.Create(dcm_path)
		if err != nil {
			return res, err
		}
		defer f.Close()
		io.Copy(f, part)
		res = append(res, dcm_path)
	}
	return res, nil
}

func WadoToFile(url string, headers map[string]string, folder string) ([]string, *http.Response, error) {
	resp, err := Get(url, headers)
	if err != nil {
		return []string{}, resp, err
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/related") {
		return []string{}, resp, errors.New(fmt.Sprintf("Content-Type is not 'multipart/related' (%s)", contentType))
	}
	dcm_path_list, err := ReadMultipartToFile(resp, folder)
	return dcm_path_list, resp, err
}
