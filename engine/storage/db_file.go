package storage

import (
	"bufio"
	"io"
	"lms-db/engine/mmap"
	"os"
	"sync"
)

type FileAccess struct {
	path      string
	size      int64      // file size
	using     int32      // using routing number
	writeFile *os.File   // fd
	readFile  *mmap.MMap // mmap read file
	mutex     sync.RWMutex
}

func NewFileAccess(path string, using int32) (*FileAccess, error) {
	fa := &FileAccess{path: path, using: using}
	flag := os.O_WRONLY | os.O_CREATE | os.O_SYNC
	if file, err := os.OpenFile(path, flag, os.ModePerm); err == nil {
		fa.writeFile = file
		if fi, err := file.Stat(); err == nil {
			fa.size = fi.Size()
		} else {
			return fa, err
		}
	} else {
		return fa, err
	}
	if file, err := mmap.NewMMap(path); err != nil {
		return fa, err
	} else {
		fa.readFile = file
	}
	return fa, nil
}

// ReadAt read file
func (f *FileAccess) ReadAt(data []byte, offset int64) (int, error) {
	f.mutex.RLocker()
	defer f.mutex.RUnlock()
	return f.readFile.ReadAt(data, offset)
}

func (f *FileAccess) ReadOneLine() ([]byte, error) {
	if r, err := f.readBatchLine(1); err != nil {
		return nil, err
	} else {
		if len(r) == 0 {
			return []byte{}, nil
		} else {
			return r[0], nil
		}
	}
}

func (f *FileAccess) readBatchLine(count int) ([][]byte, error) {
	result := make([][]byte, 0, 0)
	var err error
	bf := bufio.NewReader(f.readFile)
	for i := count; i > 0; i-- {
		if r, err := bf.ReadBytes(byte('\n')); err == io.EOF {
			break
		} else if err != nil && err != io.EOF {
			return result, err
		} else {
			result = append(result, r)
		}
	}
	return result, err
}

func (f *FileAccess) readAllLine() ([][]byte, error) {
	result := make([][]byte, 0, 0)
	var err error
	bf := bufio.NewReader(f.readFile)
	for {
		if r, err := bf.ReadBytes(byte('\n')); err == io.EOF {
			break
		} else if err != nil && err != io.EOF {
			return result, err
		} else {
			result = append(result, r)
		}
	}
	return result, err
}

// WriteAt write file
func (f *FileAccess) WriteAt(data []byte, offset int64) (int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.writeFile.WriteAt(data, offset)
}

// WriteAppendEnd write file append end
func (f *FileAccess) WriteAppendEnd(data []byte) (int, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if n, err := f.writeFile.Seek(0, io.SeekEnd); err != nil {
		return -1, err
	} else {
		r := make([]byte, 0, len(data)+1)
		r = append(r, data...)
		r = append(r, []byte("\n")...)
		return f.writeFile.WriteAt(data, n)
	}
}

func (f *FileAccess) Close() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	var err error = nil
	err = f.readFile.Close()
	err = f.writeFile.Close()
	return err
}
