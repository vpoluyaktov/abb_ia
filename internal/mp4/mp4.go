package mp4

import (
	"fmt"
	"os"

	"github.com/abema/go-mp4"
	"github.com/sunfish-shogi/bufseekio"
)

const (
	DataTypeBinary     = 0
	DataTypeStringUTF8 = 1
	DataTypeJPEG       = 14
	DataTypePNG        = 13
)

type Mp4 struct {
	fileName string
	mp4Tags  Mp4Tags
}

type Mp4Tags map[string]Mp4Tag
type Mp4Tag struct {
	Name     string
	Path     string
	DataType uint32
	Data     []byte
	Exists   bool
}

func NewMp4(fileName string) (*Mp4, error) {
	m4b := &Mp4{fileName: fileName}
	tags, err := m4b.GetMp4Tags()
	if err != nil {
		return nil, err
	}
	m4b.mp4Tags = tags
	return m4b, nil
}

func (m4b *Mp4) GetMp4Tags() (Mp4Tags, error) {
	if m4b.mp4Tags != nil {
		return m4b.mp4Tags, nil
	}
	inputFile, _ := os.Open(m4b.fileName)
	defer inputFile.Close()
	tags := make(Mp4Tags)
	r := bufseekio.NewReadSeeker(inputFile, 128*1024, 4)
	mp4.ReadBoxStructure(r, func(h *mp4.ReadHandle) (interface{}, error) {
		if h.BoxInfo.Context.UnderIlst && h.BoxInfo.Type != mp4.BoxTypeData() {
			tags[h.BoxInfo.Type.String()] = Mp4Tag{
				Name:   h.BoxInfo.Type.String(),
				Path:   getPath(h.Path),
				Exists: true,
			}
		} else if h.BoxInfo.Context.UnderIlstMeta && h.BoxInfo.Type == mp4.BoxTypeData() {
			tagName := h.Path[len(h.Path)-2].String()
			tag := tags[tagName]
			box, _, err := h.ReadPayload()
			if err != nil && box == nil {
				return nil, err
			}
			boxData := box.(*mp4.Data)
			tag.DataType = boxData.DataType
			tag.Data = boxData.Data
			tags[tagName] = tag
		}
		if h.BoxInfo.IsSupportedType() {
			h.Expand()
		}
		return nil, nil
	})
	m4b.mp4Tags = tags
	return m4b.mp4Tags, nil
}

func (m4b *Mp4) SetMp4Tag(tag *Mp4Tag) error {
	t := m4b.mp4Tags[tag.Name]
	if !t.Exists {
		t.Name = tag.Name
		t.DataType = tag.DataType
		t.Exists = false
		t.Data = tag.Data
	} else {
		t.Data = tag.Data
	}
	m4b.mp4Tags[tag.Name] = t
	return nil
}

func (m4b *Mp4) SetTag(name string, tag string) error {
	if len(name) != 4 {
		return fmt.Errorf("tag name must be 4 characters exactly")
	}
	t := m4b.mp4Tags[name]
	if !t.Exists {
		t.Name = name
		t.DataType = mp4.DataTypeStringUTF8
		t.Exists = false
		t.Data = []byte(tag)
	} else {
		t.Data = []byte(tag)
	}
	m4b.mp4Tags[name] = t
	return nil
}

func (m4b *Mp4) SetImage(imageData []byte, imageType uint32) error {
	t := m4b.mp4Tags["covr"]
	if !t.Exists {
		t.Name = "covr"
		t.DataType = imageType
		t.Exists = false
		t.Data = imageData
	} else {
		t.Data = imageData
	}
	m4b.mp4Tags["covr"] = t
	return nil
}

func (m4b *Mp4) Save() error {
	inputFileName := m4b.fileName
	outputFileName := inputFileName + ".tmp"

	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return fmt.Errorf("can't open %s: %v", inputFileName, err)
	}
	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("can't create temporary file %s: %v", outputFileName, err)
	}
	defer inputFile.Close()
	defer outputFile.Close()

	r := bufseekio.NewReadSeeker(inputFile, 128*1024, 4)
	w := mp4.NewWriter(outputFile)

	mp4.ReadBoxStructure(r, func(h *mp4.ReadHandle) (interface{}, error) {
		if !h.BoxInfo.IsSupportedType() {
			// copy all data for unsupported box types
			return nil, w.CopyBox(r, &h.BoxInfo)
		}

		// write moov box header
		_, err := w.StartBox(&h.BoxInfo)
		if err != nil {
			return nil, err
		}

		// read payload
		box, _, err := h.ReadPayload()
		if err != nil && box == nil {
			return nil, err
		}
		for _, tag := range m4b.mp4Tags {
			if h.BoxInfo.Type == mp4.BoxTypeIlst() && !tag.Exists {
				// create new tag
				w.StartBox(&mp4.BoxInfo{Type: mp4.StrToBoxType(tag.Name)}) // meta container
				w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeData()}) // data container
				dataContainer := &mp4.Data{
					DataType: tag.DataType,
					DataLang: 0x00,
					Data:     []byte(tag.Data),
				}
				mp4.Marshal(w, dataContainer, mp4.Context{UnderIlst: true, UnderIlstMeta: true})
				w.EndBox() // data container
				w.EndBox() // meta container
			} else if getPath(h.Path) == tag.Path+"data/"  && h.BoxInfo.Type == mp4.BoxTypeData() {
				// update existing tag
				boxData := box.(*mp4.Data)
				boxData.Data = []byte(tag.Data)
			}
		}

		// write box playload
		if _, err := mp4.Marshal(w, box, h.BoxInfo.Context); err != nil {
			return nil, err
		}
		// expand all of offsprings
		if _, err := h.Expand(); err != nil {
			return nil, err
		}
		// rewrite box size
		_, err = w.EndBox()
		return nil, err
	})

	inputFile.Close()
	outputFile.Close()
	// rename temporary file to final one
	os.Remove(inputFileName)
	os.Rename(outputFileName, inputFileName)
	return nil
}

func getPath(hPath mp4.BoxPath) string {
	path := ""
	for _, p := range hPath {
		path += p.String() + "/"
	}
	return path
}
