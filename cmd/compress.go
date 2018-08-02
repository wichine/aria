package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var exclude = map[string]bool{
	".DS_Store": true,
	".idea":     true,
	".git":      true,
	"cmd":       true,
	"readme.md": true,
	"img":       true,
}

// compress aria assets and inject to assets.go
func InjectAssets() error {
	content, err := compress("../../aria")
	if err != nil {
		return err
	}

	// ioutil.WriteFile("temp.tar.gz", content, 0666)

	encodeContent := base64.StdEncoding.EncodeToString(content)

	templateFileContent, err := ioutil.ReadFile("assets.go")
	if err != nil {
		return err
	}
	templateSlice := bytes.SplitN(templateFileContent, []byte("`"), -1)
	if len(templateSlice) != 3 {
		panic("assets.go is not a correct template!")
	}
	templateSlice[1] = []byte(encodeContent)

	newContent := bytes.Join(templateSlice, []byte("`"))
	err = ioutil.WriteFile("assets.go", newContent, 0644)
	if err != nil {
		return err
	}
	return nil
}

// restore assets.go to original template
func RestoreAssets() error {
	templateFileContent, err := ioutil.ReadFile("assets.go")
	if err != nil {
		return err
	}
	templateSlice := bytes.SplitN(templateFileContent, []byte("`"), -1)
	if len(templateSlice) != 3 {
		panic("assets.go is not in correct format!")
	}
	templateSlice[1] = []byte("template")

	newContent := bytes.Join(templateSlice, []byte("`"))
	err = ioutil.WriteFile("assets.go", newContent, 0644)
	if err != nil {
		return err
	}
	return nil
}

func UnpackAssets(gzByte []byte, projectName, rootPath string) error {
	if gzByte == nil || len(gzByte) == 0 {
		return fmt.Errorf("gzByte must be a gz file content with base64 encoding")
	}

	gzContent, err := base64.StdEncoding.DecodeString(string(gzByte))
	if err != nil {
		return fmt.Errorf("decode base64 error: %s", err)
	}

	// tf, err := ioutil.TempFile("./", ".tmp.")
	// if err != nil {
	// 	return fmt.Errorf("create temp file error: %s", err)
	// }
	// _, err = tf.Write(gzContent)
	// if err != nil {
	// 	return fmt.Errorf("write to temp file error: %s", err)
	// }
	// tf.Close()
	// buf, err := os.Open(tf.Name())
	// if err != nil {
	// 	return fmt.Errorf("open temp file error: %s", err)
	// }
	// defer buf.Close()
	buf := bytes.NewReader(gzContent)
	err = decompress(buf, projectName, rootPath)
	if err != nil {
		return fmt.Errorf("decompress error: %s", err)
	}
	return nil
}

func compress(dir string) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	// buf, err := ioutil.TempFile("./", ".tmp.")
	// if err != nil {
	// 	return nil, fmt.Errorf("create new temp gz file error: %s", err)
	// }
	// defer buf.Close()
	gw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	// gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("get abs of path [%s] error: %s", dir, err)
	}
	err = docompress(dir, tw, filepath.Dir(abs))
	if err != nil {
		return nil, err
	}
	tw.Close()
	gw.Close()

	content, err := ioutil.ReadAll(buf)
	// content, err := ioutil.ReadFile(buf.Name())
	if err != nil {
		return nil, err
	}
	return content, nil
}

func docompress(dir string, tw *tar.Writer, relativeBase string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("get abs of path [%s] error: %s", dir, err)
	}
	fileInfos, err := ioutil.ReadDir(abs)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		if exclude[fileInfo.Name()] {
			continue
		}
		subPath := filepath.Join(abs, fileInfo.Name())
		if fileInfo.IsDir() {
			docompress(subPath, tw, relativeBase)
			continue
		}
		f, err := os.Open(subPath)
		defer f.Close()
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(fileInfo, "")
		if err != nil {
			return err
		}
		header.Name, err = filepath.Rel(relativeBase, subPath)
		if err != nil {
			return err
		}
		fmt.Println(header.Name)

		err = tw.WriteHeader(header)
		if err != nil {
			return fmt.Errorf("write header error: %s", err)
		}

		_, err = io.Copy(tw, f)
		if err != nil {
			return fmt.Errorf("write file to gz file error: %s", err)
		}

	}
	return nil
}

func decompress(reader io.Reader, replaceRoot, rootPath string) error {
	gr, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("new gzip reader error: %s", err)
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("read header error: %s", err)
			}
		}
		fileName := header.Name
		if replaceRoot != "" {
			s := strings.Split(fileName, string(filepath.Separator))
			if len(s) > 0 {
				s[0] = replaceRoot
			}
			fileName = filepath.Join(s...)
		}
		if header.FileInfo().IsDir() {
			continue
		}
		fileNameAbs := filepath.Join(rootPath, fileName)
		fmt.Println(fileNameAbs)
		err = os.MkdirAll(filepath.Dir(fileNameAbs), 0755)
		if err != nil {
			return fmt.Errorf("make dir error: %s", err)
		}
		// f, err := os.Create(fileName)
		// if err != nil {
		// 	return err
		// }
		// io.Copy(f, tr)
		// f.Close()
		buf := bytes.NewBuffer([]byte{})
		io.Copy(buf, tr)
		fileContent, err := ioutil.ReadAll(buf)
		if err != nil {
			return fmt.Errorf("read content from buffer error: %s", err)
		}
		newFileContent := bytes.Replace(fileContent, []byte(`"aria/`), []byte(fmt.Sprintf(`"%s/`, replaceRoot)), -1)
		err = ioutil.WriteFile(fileNameAbs, newFileContent, 0666)
		if err != nil {
			return fmt.Errorf("write content to file error: %s", err)
		}
	}
	return nil
}
