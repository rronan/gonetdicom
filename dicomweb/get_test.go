package dicomweb

import (
	"fmt"
	"os"
	"testing"

	"github.com/rronan/gonetdicom/dicomutil"
	"github.com/suyashkumar/dicom"
)

func Test_Wado(t *testing.T) {
	url := getenv("MILVUE_API_URL", "") + "/v3/studies/1.2.826.0.1.3680044.0.0.0.20221228121333.16387?inference_command=smarturgences&?signed_url=false"
	token := getenv("MILVUE_TOKEN", "")
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	dcm_slice, _, err := Wado(url, headers)
	if err != nil {
		t.Fatal(err)
	}
	for _, dcm := range dcm_slice {
		study_instance_uid, series_instance_uid, sop_instance_uid, err := dicomutil.GetUIDs(dcm)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%s/%s/%s\n", study_instance_uid, series_instance_uid, sop_instance_uid)
	}
}

func Test_WadoToFile(t *testing.T) {
	url := getenv("MILVUE_API_URL", "") + "/v3/studies/1.2.826.0.1.3680044.0.0.0.20221228121333.16387?inference_command=smarturgences&signed_urls=false"
	token := getenv("MILVUE_TOKEN", "")
	headers := map[string]string{"x-goog-meta-owner": token, "Content-Type": "multipart/related; type=application/dicom"}
	dcm_path_slice, _, err := WadoToFile(url, headers, "../data/outdir")
	if err != nil {
		t.Fatal(err)
	}
	for _, dcm_path := range dcm_path_slice {
		fmt.Println(dcm_path)
		dcm, err := dicom.ParseFile(dcm_path, nil)
		if err != nil {
			t.Fatal(err)
		}
		err = os.Remove(dcm_path)
		if err != nil {
			t.Fatal(err)
		}
		study_instance_uid, series_instance_uid, sop_instance_uid, err := dicomutil.GetUIDs(&dcm)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%s/%s/%s\n", study_instance_uid, series_instance_uid, sop_instance_uid)
	}
}
