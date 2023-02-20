package dicomutil

import (
	"fmt"
	"testing"

	"github.com/suyashkumar/dicom"
)

var DICOM_PATH = "../data/study/1.2.276.0.7230010.3.1.4.0.78767.1672226121.633599.dcm"

func Test_TrimTag(t *testing.T) {
	tag := "[1.2.3.4]"
	res := TrimTag(tag)
	if res != "1.2.3.4" {
		panic(res)
	}
}

func Test_GetUIDs(t *testing.T) {
	dcm, err := dicom.ParseFile(DICOM_PATH, nil)
	if err != nil {
		panic(err)
	}
	study_instance_uid, series_instance_uid, sop_instance_uid, err := GetUIDs(&dcm)
	if err != nil {
		panic(err)
	}
	fmt.Println(study_instance_uid, series_instance_uid, sop_instance_uid)
}

func Test_ParseFileUIDs(t *testing.T) {
	study_instance_uid, series_instance_uid, sop_instance_uid, err := ParseFileUIDs(DICOM_PATH)
	if err != nil {
		panic(err)
	}
	fmt.Println(study_instance_uid, series_instance_uid, sop_instance_uid)
}

func Test_Dicom2Bytes(t *testing.T) {
	dcm, err := dicom.ParseFile(DICOM_PATH, nil)
	if err != nil {
		panic(err)
	}
	Dicom2Bytes(&dcm)
}

func Test_Bytes2Dicom(t *testing.T) {
	dcm, err := dicom.ParseFile(DICOM_PATH, nil)
	if err != nil {
		panic(err)
	}
	bytes := Dicom2Bytes(&dcm)
	_, err = Bytes2Dicom(*bytes)
}

func Test_RandomDicomName(t *testing.T) {
	_ = RandomDicomName()
}
