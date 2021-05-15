package bundle

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"

	opabundle "github.com/open-policy-agent/opa/bundle"
)

var NullBundle = opabundle.Bundle{}

//go:embed static/*
var content embed.FS

type Client struct {
	sync.RWMutex
	opaBundle        opabundle.Bundle
	bundled          bool
	opaBundleArchive []byte
	archived         bool
	revision         string
	archiveRevision  string
}

func NewClient() *Client {
	return &Client{
		bundled:  false,
		archived: false,
	}
}

func (c *Client) Get(data []byte, revision string) (opabundle.Bundle, error) {
	if !c.bundled || c.revision != revision {
		err := c.generate(data, revision)
		if err != nil {
			return NullBundle, err
		}
	}

	return c.opaBundle, nil
}

func (c *Client) GetArchive(data []byte, revision string) ([]byte, error) {
	if !c.archived || c.archiveRevision != revision {
		err := c.generate(data, revision)
		if err != nil {
			return nil, err
		}
	}

	err := c.generateArchive()
	if err != nil {
		return nil, err
	}

	return c.opaBundleArchive, nil
}

func (c *Client) generateArchive() error {
	c.Lock()
	defer c.Unlock()

	if !c.bundled {
		return errors.New("No bundle created, run Generate() before running GenerateArchive().")
	}

	if c.archived && c.opaBundle.Manifest.Revision == c.archiveRevision {
		return nil
	}

	var buf bytes.Buffer
	writer := opabundle.NewWriter(&buf).UseModulePath(true)

	err := writer.Write(c.opaBundle)
	if err != nil {
		return err
	}

	archive := buf.Bytes()
	c.opaBundleArchive = archive
	c.archived = true
	c.archiveRevision = c.opaBundle.Manifest.Revision

	return nil
}

func (c *Client) generate(data []byte, revision string) error {
	c.Lock()
	defer c.Unlock()

	if revision == c.revision && c.bundled {
		return nil
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}

	defer removeDir(tmpDir)

	err = listAndWriteStaticFiles(tmpDir)
	if err != nil {
		return err
	}

	err = writeDataFile(tmpDir, data)
	if err != nil {
		return err
	}

	b, err := newOpaBundle(tmpDir, revision)
	if err != nil {
		return err
	}

	c.opaBundle = b
	c.bundled = true
	c.revision = revision

	return nil
}

func writeDataFile(dir string, data []byte) error {
	dataFilePath := fmt.Sprintf("%s/data.json", dir)
	return os.WriteFile(dataFilePath, data, 0600)
}

func removeDir(dir string) {
	_ = os.RemoveAll(dir)
}

func listAndWriteStaticFiles(dir string) error {
	entries, err := content.ReadDir("static")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		err := writeStaticFiles(dir, entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeStaticFiles(dir string, entry fs.DirEntry) error {
	fileInfo, err := entry.Info()
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		fileName := fileInfo.Name()
		err := writeStaticFile(fileName, dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeStaticFile(fileName string, dir string) error {
	filePath := fmt.Sprintf("static/%s", fileName)
	byte, err := content.ReadFile(filePath)
	if err != nil {
		return err
	}

	newFilePath := fmt.Sprintf("%s/%s", dir, fileName)
	err = os.WriteFile(newFilePath, byte, 0600)
	if err != nil {
		return err
	}
	return nil
}

func newOpaBundle(dir string, revision string) (opabundle.Bundle, error) {
	loader := opabundle.NewDirectoryLoader(dir)
	reader := opabundle.NewCustomReader(loader).WithSkipBundleVerification(true)
	b, err := reader.Read()
	if err != nil {
		return NullBundle, err
	}

	b.Manifest.Revision = revision

	return b, nil
}
