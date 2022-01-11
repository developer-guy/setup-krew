package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/hashicorp/go-getter"
)

var (
	downloadURLTemplate       = "https://github.com/kubernetes-sigs/krew/releases/download/%s/krew-%s_%s.tar.gz"
	downloadSha256URLTemplate = "https://github.com/kubernetes-sigs/krew/releases/download/%s/krew-%s_%s.tar.gz.sha256"
	version                   string
)

func main() {
	flag.StringVar(&version, "version", "", "the version of the krew release")
	flag.Parse()
	var err error
	if version == "latest" {
		version, err = getTheLatestVersion()
		if err != nil {
			log.Fatalf("could not found the latest release of kubernetes-sigs/krew: %v", err)
		}
	}

	// https://github.com/kubernetes-sigs/krew/releases/download/v0.4.2/krew-darwin_amd64.tar.gz.sha256
	releaseDownloadURL := fmt.Sprintf(downloadURLTemplate, version, runtime.GOOS, runtime.GOARCH)
	td, err := ioutil.TempDir("", "setup-krew")
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
	}
	fileName := fmt.Sprintf("krew-%s_%s.tar.gz", runtime.GOOS, runtime.GOARCH)
	err = installAndExtractTheTarGZFiles(td, releaseDownloadURL, fileName, true)
	if err != nil {
		log.Fatalf("could not install %s from URL %s: %v", fileName, releaseDownloadURL, err)
	}

	// https://github.com/kubernetes-sigs/krew/releases/download/v0.4.2/krew-darwin_amd64.tar.gz.sha256
	sha256FileDownloadURL := fmt.Sprintf(downloadSha256URLTemplate, version, runtime.GOOS, runtime.GOARCH)
	sha256fileName := fmt.Sprintf("krew-%s_%s.tar.gz.sha256", runtime.GOOS, runtime.GOARCH)
	err = installAndExtractTheTarGZFiles(td, sha256FileDownloadURL, sha256fileName, true)
	if err != nil {
		log.Fatalf("could not install %s from URL %s: %v", sha256fileName, sha256FileDownloadURL, err)
	}

	home := os.Getenv("HOME")
	installationPath := fmt.Sprintf("%s/%s/%s", home, ".setup-krew", "bin")
	_, err = exec.Command("mkdir", "-p", installationPath).CombinedOutput()
	if err != nil {
		log.Fatalf("could not create the home directory for krew: %v", err)
	}

	_, err = exec.Command("mv", filepath.Join(td, fmt.Sprintf("krew-%s_%s", runtime.GOOS, runtime.GOARCH)),
		filepath.Join(installationPath, "krew")).CombinedOutput()
	if err != nil {
		log.Fatalf("could not rename the binary: %v", err)
	}

	_, err = exec.Command("sh", "-c", fmt.Sprintf("echo \"%s\" >> %s", installationPath, os.Getenv("GITHUB_PATH"))).CombinedOutput()
	if err != nil {
		log.Fatalf("could not add binary to \"GITHUB_PATH\": %v", err)
	}
}

func installAndExtractTheTarGZFiles(dst string, downloadURL string, fileName string, verbose bool) error {
	u, err := url.Parse(downloadURL)
	if err != nil {
		return err
	}

	fg := new(getter.HttpGetter)
	binaryTarGzFile := filepath.Join(dst, fileName)
	err = fg.GetFile(binaryTarGzFile, u)
	if err != nil {
		return err
	}

	if strings.HasSuffix(fileName, ".tar.gz") {
		log.Printf(".tar.gz file found: %s, extracting..\n", fileName)
		tarDecompressor := new(getter.TarGzipDecompressor)
		err = tarDecompressor.Decompress(dst, binaryTarGzFile, true, 0600)
		if err != nil {
			return err
		}
	}

	if verbose {
		log.Println("verbose mode enabled")
		currentDirFiles, err := exec.Command("ls", "-latr", dst).CombinedOutput()
		if err != nil {
			return err
		}
		log.Println(string(currentDirFiles))
	}
	return nil
}

func getTheLatestVersion() (string, error) {
	client := github.NewClient(nil)
	t, _, err := client.Repositories.ListTags(context.Background(), "kubernetes-sigs", "krew", &github.ListOptions{})
	if err != nil {
		return "", err
	}

	if len(t) == 0 {
		return "", fmt.Errorf("could not any valid tag for kubernetes-sigs/krew")
	}

	return *t[0].Name, nil
}
