package dicomutil

import (
	"bytes"
	"strings"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/frame"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func TrimTag(s string) string {
	s = strings.TrimPrefix(s, "[")
	return strings.TrimSuffix(s, "]")
}

func GetUIDs(dataset *dicom.Dataset) (string, string, error) {
	element, err := dataset.FindElementByTag(tag.SOPInstanceUID)
	if err != nil {
		return "", "", err
	}
	sop_instance_uid := element.Value.String()
	element, err = dataset.FindElementByTag(tag.StudyInstanceUID)
	if err != nil {
		return "", "", err
	}
	study_instance_uid := element.Value.String()
	return TrimTag(study_instance_uid), TrimTag(sop_instance_uid), nil
}

func ParseDataset(dcm_path string) (string, string, error) {
	dataset, err := dicom.ParseFile(dcm_path, nil)
	if err != nil {
		return "", "", err
	}
	return GetUIDs(&dataset)
}

func Bytes2Dicom(b []byte) (*dicom.Dataset, error) {
	reader := bytes.NewReader(b)
	dcm, err := dicom.Parse(reader, int64(len(b)), nil)
	return &dcm, err
}

func Dicom2Bytes(dcm *dicom.Dataset) *[]byte {
	buf := new(bytes.Buffer)
	dicom.Write(buf, *dcm, dicom.DefaultMissingTransferSyntax(), dicom.SkipVRVerification(), dicom.SkipValueTypeVerification())
	b := buf.Bytes()
	return &b
}

var NULL_PIXEL, _ = dicom.NewElement(tag.PixelData, dicom.PixelDataInfo{
	IsEncapsulated: false,
	Frames: []frame.Frame{
		{
			Encapsulated: false,
			NativeData: frame.NativeFrame{
				BitsPerSample: 0,
				Rows:          0,
				Cols:          0,
				Data:          [][]int{{}},
			},
		},
	},
})
