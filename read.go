package unzip

import (
	"context"
	"errors"
)

// ReadFileByName todo:多个文件并发操作
// 通过文件名远程读取指定文件，返回key:path value:fileContent(可能存在同名，以不同key:path区分)
func (zr *ZipReader) ReadFileByName(c context.Context, files []string) (fileContent map[string][]byte, err error) {
	if zr == nil {
		return nil, ErrZipReader
	}
	if zr.directory == nil || len(zr.directory.children) == 0 {
		return nil, ErrZipReaderDirectory
	}
	fileContent = make(map[string][]byte)
	for _, f := range files {
		retFiles := zr.findFileNode(zr.directory.children[0], f)
		for _, rf := range retFiles {
			fileStream, err := zr.readFile(c, rf.file)
			if err != nil {
				return nil, err
			}
			fileContent[rf.filePath] = fileStream
		}
	}
	return fileContent, nil
}

// ReadFileByPath 通过完整路径+文件名远程读取指定文件
func (zr *ZipReader) ReadFileByPath(c context.Context, filePath string) (fileContent []byte, err error) {
	if zr == nil {
		return nil, ErrZipReader
	}
	if zr.directory == nil || len(zr.directory.children) == 0 {
		return nil, ErrZipReaderDirectory
	}
	retFiles := zr.findFileNodeByPath(zr.directory.children[0], filePath)
	if retFiles == nil {
		return nil, errors.New("file not found")
	}
	fileStream, err := zr.readFile(c, retFiles.file)
	if err != nil {
		return nil, err
	}
	return fileStream, nil
}
