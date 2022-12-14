package dicomweb

import (
	"testing"

	"github.com/suyashkumar/dicom"
)

func Test_Post(t *testing.T) {
	url := getenv("MILVUE_URL", "")
	token := getenv("MILVUE_TOKEN", "")
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	dcm_slice := []*dicom.Dataset{{}, {}}
	err := Post(url, dcm_slice, headers)
	if err != nil {
		panic(err)
	}
}
