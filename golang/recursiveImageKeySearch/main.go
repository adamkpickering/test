package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	images := make([]string, 0)

	if err := extractImagesFromHelmTemplate(&images); err != nil {
		return fmt.Errorf("failed to extract images from helm template: %w", err)
	}
	fmt.Printf("images after helm template: %v\n", images)

	return nil
}

func extractImagesFromHelmTemplate(images *[]string) error {
	cmd := exec.Command("helm", "template", "test-test", "calico-operator/tigera-operator")
	templateOutput := &bytes.Buffer{}
	cmd.Stdout = templateOutput
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run template command: %w", err)
	}

	decoder := yaml.NewDecoder(templateOutput)
	for {
		var templateData any
		if err := decoder.Decode(&templateData); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("failed to parse template command output as yaml: %w", err)
		}

		err := extractImagesFromYaml(templateData, images)
		if err != nil {
			return fmt.Errorf("failed to extract images: %w", err)
		}
	}

	return nil
}

func extractImagesFromYaml(data any, images *[]string) error {
	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			if key == "image" {
				stringValue, ok := value.(string)
				if ok {
					*images = append(*images, stringValue)
				}
			} else {
				if err := extractImagesFromYaml(value, images); err != nil {
					return err
				}
			}
		}
	case []any:
		for _, item := range v {
			if err := extractImagesFromYaml(item, images); err != nil {
				return err
			}
		}
	}

	return nil
}
