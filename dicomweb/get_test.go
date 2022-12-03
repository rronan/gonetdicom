package dicomweb

import (
	"fmt"
	"testing"

	"github.com/rronan/gonetdicom/dicomutil"
)

func Test_Get(t *testing.T) {
	url := getenv("MILVUE_URL", "")
	token := getenv("MILVUE_TOKEN", "")
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	dcm_slice, b, err := GetDicomWeb(url, headers)
	if err != nil {
		panic(err)
	}
	for _, dcm := range dcm_slice {
		study_instance_uid, sop_instance_uid, err := dicomutil.GetUIDs(dcm)
		if err != nil {
			panic(err)
		}
		fmt.Println(study_instance_uid)
		fmt.Println(sop_instance_uid)
	}
	fmt.Println(string(b))
}
