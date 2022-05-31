package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func LineCounter(filePath string) (int, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := f.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func ReadCsv(filePath string) ([][]string, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, f.Sync()
}

func ReadCsvByLine(filePath string, line int) ([][]string, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 判断文件行数
	c, err := LineCounter(filePath)
	if err != nil {
		return nil, err
	}

	if line == 0 {
		line = c
	} else if line > c {
		return nil, fmt.Errorf("this file lins is: %v, is too short,except: %v", c, line)
	}

	r := csv.NewReader(f)
	d := [][]string{}
	//针对大文件，一行一行的读取文件
	for i := 0; i < line; i++ {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			return nil, err
		}
		// if err == io.EOF {
		// 	break
		// }
		d = append(d, row)
	}

	return d, f.Sync()
}

func WriteCsv(filePath string, records [][]string) error {

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.WriteAll(records)
	if err := w.Error(); err != nil {
		return err
	}

	return f.Sync()
}

func PathExists(path string) bool {

	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// if path not exist, then to create
func PathNEAC(path string) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(err)
		}
	}
}
