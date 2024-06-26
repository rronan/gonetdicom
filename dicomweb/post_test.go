package dicomweb

import (
	"testing"

	"github.com/suyashkumar/dicom"
)

func Test_Stow(t *testing.T) {
	url := getenv("MILVUE_API_URL", "") + "/v3/studies?signed_url=false"
	token := getenv("MILVUE_TOKEN", "")
	headers := map[string]string{
		"x-goog-meta-owner": token,
		"Content-Type":      "multipart/related; type=application/dicom",
		"Accept":            "application/json",
	}
	var DICOM_PATH_SLICE = []string{
		"../data/study/1.2.276.0.7230010.3.1.4.0.78767.1672226121.633599.dcm",
		"../data/study/1.2.276.0.7230010.3.1.4.0.78767.1672226121.633601.dcm",
	}
	dcm_slice := []*dicom.Dataset{}
	for _, path := range DICOM_PATH_SLICE {
		dcm, err := dicom.ParseFile(path, nil)
		if err != nil {
			t.Fatal(err)
		}
		dcm_slice = append(dcm_slice, &dcm)
	}
	resp, err := Stow(url, dcm_slice, headers, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func Test_StowFromFile(t *testing.T) {
	url := getenv("MILVUE_API_URL", "") + "/v3/studies?signed_url=false"
	token := getenv("MILVUE_TOKEN", "")
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	var DICOM_PATH_SLICE = []string{
		"../data/study/1.2.276.0.7230010.3.1.4.0.78767.1672226121.633599.dcm",
		"../data/study/1.2.276.0.7230010.3.1.4.0.78767.1672226121.633601.dcm",
	}
	resp, err := StowFromFile(url, DICOM_PATH_SLICE, headers, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}
