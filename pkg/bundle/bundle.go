package bundle

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"fmt"
	"os"

	opabundle "github.com/open-policy-agent/opa/bundle"
)

var NullBundle = opabundle.Bundle{}

//go:embed static/*
var content embed.FS

func NewTarGzBundle(b opabundle.Bundle) ([]byte, error) {
	var buf bytes.Buffer
	writer := opabundle.NewWriter(&buf)

	err := writer.Write(b)
	if err != nil {
		return nil, err
	}

	var res []byte
	_, err = buf.Read(res)
	if err != nil {
		return nil, err
	}

	os.WriteFile("/tmp/test.tar.gz", res, 0600)

	return res, nil
}

func NewBundle(data []byte) (opabundle.Bundle, error) {
	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return NullBundle, err
	}

	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	entries, err := content.ReadDir("static")
	if err != nil {
		return NullBundle, err
	}

	for _, entry := range entries {
		fileInfo, err := entry.Info()
		if err != nil {
			return NullBundle, err
		}

		if !fileInfo.IsDir() {
			fileName := fileInfo.Name()
			filePath := fmt.Sprintf("static/%s", fileName)
			byte, err := content.ReadFile(filePath)
			if err != nil {
				return NullBundle, err
			}
			newFilePath := fmt.Sprintf("%s/%s", tmpDir, fileName)
			err = os.WriteFile(newFilePath, byte, 0600)
			if err != nil {
				return NullBundle, err
			}
		}
	}

	dataFilePath := fmt.Sprintf("%s/data.json", tmpDir)
	err = os.WriteFile(dataFilePath, data, 0600)
	if err != nil {
		return NullBundle, err
	}

	loader := opabundle.NewDirectoryLoader(tmpDir)
	reader := opabundle.NewCustomReader(loader).WithSkipBundleVerification(true)
	b, err := reader.Read()
	if err != nil {
		return NullBundle, err
	}

	h := sha256.New()
	_, err = h.Write(data)
	if err != nil {
		return NullBundle, err
	}

	dataHash := fmt.Sprintf("%x", h.Sum(nil))

	b.Manifest.Revision = dataHash

	return b, nil
}
