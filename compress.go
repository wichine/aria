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
}

var assetsPath = "./hatch"
var assetsFileName = "assets.go"
var tempTemplateFile = ".assets.go.temp"

// compress aria assets and inject to assets.go
func InjectAssets() error {
	templateFileContent, err := ioutil.ReadFile(assetsFileName)
	if err != nil {
		return err
	}
	assetsContent := string(templateFileContent)
	if strings.Index(assetsContent, "TEMPLATE") < 0 {
		return fmt.Errorf("%s is not a correct template file.", assetsFileName)
	}
	err = ioutil.WriteFile(tempTemplateFile, templateFileContent, 0644)
	if err != nil {
		return fmt.Errorf("write temp file error: %s", err)
	}

	dirList, err := getDirList(assetsPath)
	if err != nil {
		return err
	}

	for _, dirName := range dirList {
		gzByte, err := compress(filepath.Join(assetsPath, dirName))
		if err != nil {
			return err
		}
		encodeGzContent := base64.StdEncoding.EncodeToString(gzByte)
		template := fmt.Sprintf("%s_TEMPLATE", strings.ToUpper(dirName))
		assetsContent = strings.Replace(assetsContent, template, encodeGzContent, -1)
	}

	err = ioutil.WriteFile(assetsFileName, []byte(assetsContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

// restore assets.go to original template
func RestoreAssets() error {
	_, err := os.Stat(tempTemplateFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("temp file <%s> not exsits", tempTemplateFile)
	}
	err = os.Rename(tempTemplateFile, assetsFileName)
	if err != nil {
		return fmt.Errorf("rename temp file error: %s", err)
	}
	return nil
}

// unpack assets from bytes stored in assets.go, argument "relativeDirNameInHatch" must be the dirname in aria/hatch.
func UnpackAssets(gzByte []byte, projectName, rootPath, relativeDirNameInHatch string) error {
	if gzByte == nil || len(gzByte) == 0 {
		return fmt.Errorf("gzByte must be a gz file content with base64 encoding")
	}

	gzContent, err := base64.StdEncoding.DecodeString(string(gzByte))
	if err != nil {
		return fmt.Errorf("decode base64 error: %s", err)
	}

	buf := bytes.NewReader(gzContent)
	err = decompress(buf, projectName, rootPath, relativeDirNameInHatch)
	if err != nil {
		return fmt.Errorf("decompress error: %s", err)
	}
	return nil
}

func compress(dir string) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	gw, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
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

func decompress(reader io.Reader, replaceRoot, rootPath, relativeDirNameInHatch string) error {
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

		buf := bytes.NewBuffer([]byte{})
		io.Copy(buf, tr)
		fileContent, err := ioutil.ReadAll(buf)
		if err != nil {
			return fmt.Errorf("read content from buffer error: %s", err)
		}
		newFileContent := bytes.Replace(fileContent, []byte(fmt.Sprintf(`"aria/hatch/%s/`, relativeDirNameInHatch)), []byte(fmt.Sprintf(`"%s/`, replaceRoot)), -1)
		err = ioutil.WriteFile(fileNameAbs, newFileContent, 0666)
		if err != nil {
			return fmt.Errorf("write content to file error: %s", err)
		}
	}
	return nil
}

func getDirList(dirpath string) ([]string, error) {
	f, err := os.Open(dirpath)
	if err != nil {
		return nil, fmt.Errorf("open dir %s error: %s", dirpath, err)
	}
	return f.Readdirnames(-1)
}
