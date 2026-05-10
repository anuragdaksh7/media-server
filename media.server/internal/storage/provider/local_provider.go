package provider

import (
	"context"
	"errors"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

type LocalProvider struct {
	Root string
}

func NewLocalProvider(root string) *LocalProvider {
	return &LocalProvider{root}
}

func (p *LocalProvider) resolvePath(
	userPath string,
) (string, error) {

	// normalize slashes
	userPath = filepath.ToSlash(userPath)

	// remove leading slash
	userPath = strings.TrimPrefix(userPath, "/")

	// clean path
	clean := filepath.Clean(userPath)

	// prevent weird "." result
	if clean == "." {
		clean = ""
	}

	full := filepath.Join(
		p.Root,
		clean,
	)

	full = filepath.Clean(full)

	rel, err := filepath.Rel(
		p.Root,
		full,
	)

	if err != nil {
		return "", err
	}

	if strings.HasPrefix(rel, "..") {
		return "", errors.New("invalid path")
	}

	return full, nil
}

func (p *LocalProvider) List(
	ctx context.Context,
	path string,
) ([]FileInfo, error) {

	fullPath, err := p.resolvePath(path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	files := make([]FileInfo, 0)

	for _, entry := range entries {

		info, err := entry.Info()

		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Name:  entry.Name(),
			Path:  filepath.Join(path, entry.Name()),
			Size:  info.Size(),
			IsDir: entry.IsDir(),

			MimeType: mime.TypeByExtension(
				filepath.Ext(entry.Name()),
			),

			ModifiedAt: info.ModTime(),
		})
	}

	return files, nil
}

func (p *LocalProvider) Stat(
	ctx context.Context,
	path string,
) (*FileInfo, error) {

	fullPath, err := p.resolvePath(path)

	if err != nil {
		return nil, err
	}

	info, err := os.Stat(fullPath)

	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:  info.Name(),
		Path:  path,
		Size:  info.Size(),
		IsDir: info.IsDir(),

		MimeType: mime.TypeByExtension(
			filepath.Ext(info.Name()),
		),

		ModifiedAt: info.ModTime(),
	}, nil
}

func (p *LocalProvider) Read(
	ctx context.Context,
	path string,
) (io.ReadSeekCloser, error) {

	fullPath, err := p.resolvePath(path)

	if err != nil {
		return nil, err
	}

	return os.Open(fullPath)
}

func (p *LocalProvider) Write(
	ctx context.Context,
	path string,
	reader io.Reader,
) error {

	fullPath, err := p.resolvePath(path)

	if err != nil {
		return err
	}

	err = os.MkdirAll(
		filepath.Dir(fullPath),
		0755,
	)

	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, reader)

	return err
}

func (p *LocalProvider) Delete(
	ctx context.Context,
	path string,
) error {

	fullPath, err := p.resolvePath(path)

	if err != nil {
		return err
	}

	return os.RemoveAll(fullPath)
}

func (p *LocalProvider) Move(
	ctx context.Context,
	src string,
	dst string,
) error {

	srcPath, err := p.resolvePath(src)

	if err != nil {
		return err
	}

	dstPath, err := p.resolvePath(dst)

	if err != nil {
		return err
	}

	err = os.MkdirAll(
		filepath.Dir(dstPath),
		0755,
	)

	if err != nil {
		return err
	}

	return os.Rename(srcPath, dstPath)
}

func (p *LocalProvider) Copy(
	ctx context.Context,
	src string,
	dst string,
) error {

	srcPath, err := p.resolvePath(src)

	if err != nil {
		return err
	}

	dstPath, err := p.resolvePath(dst)

	if err != nil {
		return err
	}

	srcFile, err := os.Open(srcPath)

	if err != nil {
		return err
	}

	defer srcFile.Close()

	err = os.MkdirAll(
		filepath.Dir(dstPath),
		0755,
	)

	if err != nil {
		return err
	}

	dstFile, err := os.Create(dstPath)

	if err != nil {
		return err
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)

	return err
}

func (p *LocalProvider) Mkdir(
	ctx context.Context,
	path string,
) error {

	fullPath, err := p.resolvePath(path)

	if err != nil {
		return err
	}

	return os.MkdirAll(fullPath, 0755)
}
