package dicomweb

import (
	"testing"

	"github.com/suyashkumar/dicom"
)

func Test_Put(t *testing.T) {
	url := "http://localhost:8000/foo.dcm"
	DICOM_PATH := "../data/study/1.2.276.0.7230010.3.1.4.0.78767.1672226121.633599.dcm"
	dcm, err := dicom.ParseFile(DICOM_PATH, nil)
	if err != nil {
		panic(err)
	}
	err = Put(url, &dcm, map[string]string{})
	if err != nil {
		panic(err)
	}
}
