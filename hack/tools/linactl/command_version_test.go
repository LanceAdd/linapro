package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandRegistryIncludesVersion(t *testing.T) {
	registry := commandRegistry()
	if _, ok := registry["version"]; !ok {
		t.Fatalf("expected command %q to be registered", "version")
	}
}

func TestRunVersionUpdatesMetadataAndReadmeImages(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "apps", "lina-core", "manifest", "config", "metadata.yaml"), `framework:
  name: LinaPro
  version: "v0.3.0"
openapi:
  version: "v0.1.0"
`)
	writeFile(t, filepath.Join(root, "README.md"), `<img src="https://linapro.ai/img/logo.png" width="300" />
[![CI](https://example.com/badge.svg?style=flat)](https://example.com)
![Preview](https://linapro.ai/img/preview.webp?v=0.3.0)
`)
	writeFile(t, filepath.Join(root, "README.zh-CN.md"), `<img src='https://linapro.ai/img/zh-logo.png?old=1&v=0.3.0' />
![Preview](https://linapro.ai/img/zh-preview.webp)
`)

	var stdout bytes.Buffer
	application := newApp(&stdout, ioDiscard{}, strings.NewReader(""))
	application.root = root

	err := runVersion(context.Background(), application, commandInput{Params: map[string]string{"to": "v1.2.3"}})
	if err != nil {
		t.Fatalf("runVersion returned error: %v", err)
	}

	metadata := readFile(t, filepath.Join(root, "apps", "lina-core", "manifest", "config", "metadata.yaml"))
	if !strings.Contains(metadata, `  version: "v1.2.3"`) {
		t.Fatalf("framework.version was not updated:\n%s", metadata)
	}
	if !strings.Contains(metadata, `openapi:
  version: "v0.1.0"`) {
		t.Fatalf("non-framework version should not be changed:\n%s", metadata)
	}

	readme := readFile(t, filepath.Join(root, "README.md"))
	for _, fragment := range []string{
		`src="https://linapro.ai/img/logo.png?v=1.2.3"`,
		`https://example.com/badge.svg?style=flat&v=1.2.3`,
		`https://linapro.ai/img/preview.webp?v=1.2.3`,
	} {
		if !strings.Contains(readme, fragment) {
			t.Fatalf("README.md missing %q:\n%s", fragment, readme)
		}
	}

	chineseReadme := readFile(t, filepath.Join(root, "README.zh-CN.md"))
	for _, fragment := range []string{
		`src='https://linapro.ai/img/zh-logo.png?old=1&v=1.2.3'`,
		`https://linapro.ai/img/zh-preview.webp?v=1.2.3`,
	} {
		if !strings.Contains(chineseReadme, fragment) {
			t.Fatalf("README.zh-CN.md missing %q:\n%s", fragment, chineseReadme)
		}
	}
	if !strings.Contains(stdout.String(), "Updated framework.version to v1.2.3") {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}

func TestRunVersionRejectsInvalidTarget(t *testing.T) {
	application := newApp(ioDiscard{}, ioDiscard{}, strings.NewReader(""))
	application.root = t.TempDir()

	err := runVersion(context.Background(), application, commandInput{Params: map[string]string{"to": "1.2.3"}})
	if err == nil || !strings.Contains(err.Error(), "must match vMAJOR.MINOR.PATCH") {
		t.Fatalf("expected invalid target error, got: %v", err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}
