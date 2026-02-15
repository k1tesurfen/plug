package scaffold

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type PluginData struct {
	Name      string
	Namespace string
	Scripts   []string
}

func Generate(pluginPath string, items []string, scripts []string, templateRootDir string) error {
	// 1. Create plugin root
	if _, err := os.Stat(pluginPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", pluginPath)
	}
	if err := os.MkdirAll(pluginPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data := PluginData{
		Name:      filepath.Base(pluginPath),
		Namespace: toPascalCase(filepath.Base(pluginPath)),
		Scripts:   scripts,
	}

	// 2. Process items
	for _, itemRelPath := range items {
		// srcPath is the Absolute Path to the template file
		srcPath := filepath.Join(templateRootDir, itemRelPath)

		info, err := os.Stat(srcPath)
		if err != nil {
			return fmt.Errorf("template source not found: %s", srcPath)
		}

		if info.IsDir() {
			// RECURSION FIX:
			// We pass 'itemRelPath' so we know where this folder sits relative to templates
			if err := processDir(srcPath, pluginPath, templateRootDir, data); err != nil {
				return err
			}
		} else {
			// SINGLE FILE FIX:
			// We want "base/plugin.php" -> "my-plugin/plugin.php"
			// We DO NOT want "my-plugin/Users/name/.config/..."

			// We take just the filename (base) of the source
			destName := filepath.Base(itemRelPath)
			destPath := filepath.Join(pluginPath, destName)

			if err := processFile(srcPath, destPath, data); err != nil {
				return err
			}
		}
	}

	return nil
}

func processDir(srcDir, pluginRoot, templateRootDir string, data PluginData) error {
	// We walk inside the source directory (e.g. .../templates/base/src)
	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 1. Calculate the relative path from the specific item we are copying
		// If we are copying "base/src", and we find "base/src/Plugin.php",
		// relPath becomes "Plugin.php"
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		} // Skip the root folder itself

		// 2. Construct Destination
		// pluginRoot + "src" (from original item name) + relPath
		// Wait, if we just want to copy contents of "base/src" to "my-plugin/src",
		// we need to know the target folder name.

		// SIMPLIFICATION:
		// If config says "base/src", we want "my-plugin/src".
		// "base/src" -> base name is "src".
		targetFolder := filepath.Base(srcDir)
		destPath := filepath.Join(pluginRoot, targetFolder, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		return processFile(path, destPath, data)
	})
}

func processFile(srcPath, destPath string, data PluginData) error {

	fileName := filepath.Base(srcPath)

	// SPECIAL RULE: If the template is named "plugin.php", rename it to the plugin slug
	if fileName == "init.php" {
		fileName = data.Name + ".php"
	}

	destPath = filepath.Join(filepath.Dir(destPath), fileName)

	// Ensure parent dir exists (just in case)
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	tmplContent, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// Helper function map for templates (e.g. toLower, etc.)
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	tmpl, err := template.New(filepath.Base(srcPath)).Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", srcPath, err)
	}

	var out bytes.Buffer
	if err := tmpl.Execute(&out, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", srcPath, err)
	}

	return os.WriteFile(destPath, out.Bytes(), 0644)
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}
