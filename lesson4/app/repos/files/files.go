package files

import (
	"os"
	"strings"
)

type File struct {
	Name      string `json:"name"`
	Extension string `json:"ext"`
	Size      int64  `json:"size"`
}

type Files struct {
	Dir string
}

func NewFiles(dir string) *Files {
	return &Files{
		Dir: dir,
	}
}

func (f Files) Search(extension string) ([]File, error) {
	dirFiles, err := os.ReadDir(f.Dir)
	if err != nil {
		return nil, err
	}
	files := make([]File, 0)
	for _, dirFile := range dirFiles {
		info, err := dirFile.Info()
		if err != nil {
			continue
		}

		arr := strings.Split(info.Name(), ".")
		file_ext := arr[len(arr)-1]

		// скрытые не выводим
		if string(info.Name()[0]) == "." {
			continue
		}

		// фильтр по разширению
		if extension != "" && extension != file_ext {
			continue
		}

		files = append(files, File{
			Name:      info.Name(),
			Extension: file_ext,
			Size:      info.Size(),
		})
	}
	return files, nil
}
