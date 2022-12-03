package dicomweb

import (
	"testing"

	"github.com/suyashkumar/dicom"
)

func Test_Put(t *testing.T) {
	url := getenv("MILVUE_URL", "")
	token := getenv("MILVUE_TOKEN", "")
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	dcm := dicom.Dataset{}
	err := Put(url, &dcm, headers)
	if err != nil {
		panic(err)
	}
}
