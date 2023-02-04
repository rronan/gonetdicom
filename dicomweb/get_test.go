package dicomweb

import (
	"fmt"
	"testing"

	"github.com/rronan/gonetdicom/dicomutil"
)

func Test_Wado(t *testing.T) {
	url := getenv("MILVUE_API_URL", "") + "/v3/studies/1.2.826.0.1.3680044.0.0.0.20221228121333.16387?inference_command=smarturgences"
	token := getenv("MILVUE_TOKEN", "")
	fmt.Println("url:", url, "token:", token)
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	dcm_slice, _, err := Wado(url, headers)
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
}
