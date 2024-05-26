package jane

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"
)

const defaultWhence = 0

type DiskStore struct {
	File          *os.File
	WritePosition int
	KeyDir        map[string]KeyEntry
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || errors.Is(err, fs.ErrExist)
}

func NewDiskStore(filename string) (*DiskStore, error) {
	ds := &DiskStore{KeyDir: make(map[string]KeyEntry)}
	if isFileExist(filename) {
		if err := ds.initKeyDir(filename); err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	ds.File = file
	return ds, nil
}

func (d *DiskStore) Get(key string) (string, error) {
	kEntry, ok := d.KeyDir[key]
	if !ok {
		return "", fmt.Errorf("key %s does not exist", key)
	}

	if _, err := d.File.Seek(int64(kEntry.Position), defaultWhence); err != nil {
		return "", fmt.Errorf("error seeking to position: %v", err)
	}

	data := make([]byte, kEntry.TotalSize)
	if _, err := io.ReadFull(d.File, data); err != nil {
		return "", fmt.Errorf("read error: %v", err)
	}

	_, _, value := decodeKeyValue(data)
	return value, nil
}

func (d *DiskStore) Set(key string, value string) error {
	timestamp := uint32(time.Now().Unix())
	size, data := encodeKeyValue(timestamp, key, value)
	if err := d.write(data); err != nil {
		return err
	}

	d.KeyDir[key] = NewKeyEntry(timestamp, uint32(d.WritePosition), uint32(size))
	d.WritePosition += size
	return nil
}

func (d *DiskStore) Close() error {
	if err := d.File.Sync(); err != nil {
		return fmt.Errorf("error syncing file: %v", err)
	}
	if err := d.File.Close(); err != nil {
		return fmt.Errorf("error closing file: %v", err)
	}
	return nil
}

func (d *DiskStore) write(data []byte) error {
	if _, err := d.File.Write(data); err != nil {
		return fmt.Errorf("write error: %v", err)
	}
	if err := d.File.Sync(); err != nil {
		return fmt.Errorf("sync error: %v", err)
	}
	return nil
}

func (d *DiskStore) initKeyDir(existingFile string) error {
	file, err := os.Open(existingFile)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	for {
		header := make([]byte, headerSize)
		if _, err := io.ReadFull(file, header); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading header: %v", err)
		}

		timestamp, keySize, valueSize := decodeHeader(header)
		key := make([]byte, keySize)
		value := make([]byte, valueSize)

		if _, err := io.ReadFull(file, key); err != nil {
			return fmt.Errorf("error reading key: %v", err)
		}
		if _, err := io.ReadFull(file, value); err != nil {
			return fmt.Errorf("error reading value: %v", err)
		}

		totalSize := headerSize + keySize + valueSize
		d.KeyDir[string(key)] = NewKeyEntry(timestamp, uint32(d.WritePosition), totalSize)
		d.WritePosition += int(totalSize)
		fmt.Printf("loaded key=%s, value=%s\n", key, value)
	}

	return nil
}
