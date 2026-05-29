package app

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func (a *App) updateSelf() error {
	latest, err := a.latestReleaseTag()
	if err != nil {
		return err
	}
	current := normalizeVersion(a.version)
	target := normalizeVersion(latest)
	if current == target {
		fmt.Printf("gacha is already up to date (%s).\n", a.version)
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return err
	}

	fmt.Printf("Updating gacha %s -> %s\n", a.version, target)
	tmpDir, err := os.MkdirTemp("", "gacha-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, "gacha.tar.gz")
	url := releaseAssetURL(latest)
	if err := a.downloadFile(url, archivePath); err != nil {
		return err
	}
	if err := extractTarGz(archivePath, tmpDir); err != nil {
		return err
	}
	newBinary := filepath.Join(tmpDir, "gacha")
	if err := os.Chmod(newBinary, 0o755); err != nil {
		return err
	}

	backup := exe + ".old"
	_ = os.Remove(backup)
	if err := os.Rename(exe, backup); err != nil {
		return fmt.Errorf("cannot replace %s: %w", exe, err)
	}
	if err := copyFile(newBinary, exe, 0o755); err != nil {
		_ = os.Rename(backup, exe)
		return err
	}
	_ = os.Remove(backup)

	fmt.Printf("Updated %s to %s.\n", exe, target)
	return nil
}

func (a *App) latestReleaseTag() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/dkstm95/gacha/releases/latest", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "gacha/"+a.version)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("GitHub release check failed: %s", resp.Status)
	}
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	if release.TagName == "" {
		return "", fmt.Errorf("GitHub latest release did not include a tag")
	}
	return release.TagName, nil
}

func normalizeVersion(value string) string {
	return strings.TrimPrefix(strings.TrimSpace(value), "v")
}

func releaseAssetURL(tag string) string {
	return fmt.Sprintf("https://github.com/dkstm95/gacha/releases/download/%s/gacha-%s.tar.gz", tag, targetTriple())
}

func (a *App) downloadFile(url string, destination string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "gacha/"+a.version)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.ReadFrom(resp.Body)
	return err
}

func extractTarGz(archivePath string, destinationDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		name := filepath.Clean(header.Name)
		if name != "gacha" {
			continue
		}
		target := filepath.Join(destinationDir, "gacha")
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, tarReader)
		closeErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		return nil
	}
}

func copyFile(source string, destination string, mode fs.FileMode) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Chmod(destination, mode)
}

func targetTriple() string {
	return runtime.GOOS + "-" + runtime.GOARCH
}
