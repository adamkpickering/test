package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	err := doMain()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func doMain() error {
	fd, err := os.Open("index.html")
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer fd.Close()

	buf := &bytes.Buffer{}
	tee := io.TeeReader(fd, buf)

	fmt.Printf("before injection:\n\n")
	if _, err := io.Copy(os.Stdout, tee); err != nil {
		return fmt.Errorf("failed to print before to stdout: %w", err)
	}

	fmt.Printf("\n\nafter injection:\n\n")
	script := `<script>console.log("hello world")</script>`
	if err := injectScript(os.Stdout, buf, script); err != nil {
		return fmt.Errorf("failed to inject script: %w", err)
	}

	return nil
}

func injectScript(w io.Writer, r io.Reader, script string) error {
	document, err := html.Parse(r)
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}
	found := false
	for node := range document.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.Head {
			scriptNode := &html.Node{
				Type: html.RawNode,
				Data: script,
			}
			node.AppendChild(scriptNode)
			found = true
			break
		}
	}
	if !found {
		return errors.New("failed to find head element")
	}
	if err := html.Render(w, document); err != nil {
		return fmt.Errorf("failed to render modified document: %w", err)
	}
	return nil
}
