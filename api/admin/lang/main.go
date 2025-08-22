package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Colors for output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

// Config structure
type Config struct {
	MessagesDir      string
	TargetDir        string
	DefaultLocale    string
	SupportedLocales []string
	Strict           bool
}

// Conflict record
type Conflict struct {
	Key    string
	OldVal string
	NewVal string
	File   string
}

// cleanTargetDirectory removes old json files in targetDir
func cleanTargetDirectory(targetDir string) error {
	fmt.Printf("%sCleaning target directory...%s\n", ColorCyan, ColorReset)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("failed to read target directory: %w", err)
	}
	cleaned := false
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			if err := os.Remove(filepath.Join(targetDir, entry.Name())); err != nil {
				return fmt.Errorf("failed to delete %s: %w", entry.Name(), err)
			}
			fmt.Printf("%s✓ Deleted old file: %s%s\n", ColorGreen, entry.Name(), ColorReset)
			cleaned = true
		}
	}
	if !cleaned {
		fmt.Printf("%sTarget directory is already clean.%s\n", ColorBlue, ColorReset)
	}
	return nil
}

// getSupportedLocales scans dir for available locales
func getSupportedLocales(messagesDir, defaultLocale string) ([]string, error) {
	entries, err := os.ReadDir(messagesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages directory: %w", err)
	}
	var locales []string
	// Ensure default is first
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == defaultLocale {
			locales = append(locales, entry.Name())
		}
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		// 🚫 Skip the special output directory
		if name == "locales" {
			continue
		}
		if name != defaultLocale {
			locales = append(locales, name)
		}
	}
	return locales, nil
}

// readAndMergeJSON reads JSON file into map
func readAndMergeJSON(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}
	return parsed, nil
}

// mergeDeep merges source into target
func mergeDeep(target, source map[string]interface{}) map[string]interface{} {
	for k, v := range source {
		if existing, ok := target[k]; ok {
			tMap, tOk := existing.(map[string]interface{})
			sMap, sOk := v.(map[string]interface{})
			if tOk && sOk {
				target[k] = mergeDeep(tMap, sMap)
			} else {
				target[k] = v
			}
		} else {
			target[k] = v
		}
	}
	return target
}

// compareTranslationStructure returns missing keys
func compareTranslationStructure(defaultTranslations, translations map[string]interface{}, locale, keyPath string, logged map[string]bool) map[string]interface{} {
	missing := make(map[string]interface{})
	for key, defaultValue := range defaultTranslations {
		currentPath := key
		if keyPath != "" {
			currentPath = keyPath + "." + key
		}
		if _, exists := translations[key]; !exists {
			if !logged[currentPath] {
				fmt.Printf("%sWarning: Missing key \"%s\" in locale \"%s\".%s\n", ColorYellow, currentPath, locale, ColorReset)
				logged[currentPath] = true
			}
			missing[key] = defaultValue
		} else {
			defMap, ok1 := defaultValue.(map[string]interface{})
			trMap, ok2 := translations[key].(map[string]interface{})
			if ok1 && ok2 {
				nestedMissing := compareTranslationStructure(defMap, trMap, locale, currentPath, logged)
				if len(nestedMissing) > 0 {
					missing[key] = nestedMissing
				}
			}
		}
	}
	return missing
}

// checkConflicts checks for conflicting keys
func checkConflicts(target, source map[string]interface{}, path string, file string, conflicts *[]Conflict) {
	for key, newValue := range source {
		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}
		if oldValue, exists := target[key]; exists {
			tMap, tOk := oldValue.(map[string]interface{})
			sMap, sOk := newValue.(map[string]interface{})
			if tOk && sOk {
				checkConflicts(tMap, sMap, currentPath, file, conflicts)
			} else if oldValue != newValue {
				*conflicts = append(*conflicts, Conflict{
					Key:    currentPath,
					OldVal: fmt.Sprintf("%v", oldValue),
					NewVal: fmt.Sprintf("%v", newValue),
					File:   file,
				})
				target[key] = newValue
			}
		} else {
			target[key] = newValue
		}
	}
}

// processLocaleFiles merges all json files under a locale dir
func processLocaleFiles(locale, messagesDir string) (map[string]interface{}, error) {
	if locale == "locales" {
		return nil, fmt.Errorf("invalid locale directory: %q", locale)
	}
	dirPath := filepath.Join(messagesDir, locale)
	localeTranslations := make(map[string]interface{})
	var conflicts []Conflict
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}
		parsed, err := readAndMergeJSON(path)
		if err != nil {
			return err
		}
		checkConflicts(localeTranslations, parsed, "", path, &conflicts)
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, c := range conflicts {
		fmt.Printf("%sConflict: Key \"%s\" old=\"%s\" new=\"%s\" (from %s)%s\n",
			ColorRed, c.Key, c.OldVal, c.NewVal, c.File, ColorReset)
	}
	return localeTranslations, nil
}

// writeJSONFile writes map as json
func writeJSONFile(targetDir, filename string, data map[string]interface{}) error {
	outPath := filepath.Join(targetDir, filename)
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	if err := os.WriteFile(outPath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}
	fmt.Printf("%s✓ Wrote file: %s%s\n", ColorGreen, filename, ColorReset)
	return nil
}

// main entry
func main() {
	messagesDir := flag.String("messages", "", "Path to messages directory")
	defaultLocale := flag.String("default", "", "Default locale code (e.g., en)")
	strict := flag.Bool("strict", false, "Fail if any missing/extra keys detected")
	flag.Parse()

	if *messagesDir == "" || *defaultLocale == "" {
		log.Fatal("Usage: go run main.go -messages <dir> -default <locale> [--strict]")
	}

	locales, err := getSupportedLocales(*messagesDir, *defaultLocale)
	if err != nil {
		log.Fatal(err)
	}
	config := Config{
		MessagesDir:      *messagesDir,
		TargetDir:        *messagesDir + "/locales",
		DefaultLocale:    *defaultLocale,
		SupportedLocales: locales,
		Strict:           *strict,
	}

	if err := cleanTargetDirectory(config.TargetDir); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%sProcessing default locale %s...%s\n", ColorCyan, config.DefaultLocale, ColorReset)
	defaultTranslations, err := processLocaleFiles(config.DefaultLocale, config.MessagesDir)
	if err != nil {
		log.Fatal(err)
	}
	if err := writeJSONFile(config.TargetDir, config.DefaultLocale+".json", defaultTranslations); err != nil {
		log.Fatal(err)
	}

	logged := make(map[string]bool)
	missingSummary := make(map[string]map[string]interface{})
	for _, locale := range config.SupportedLocales {
		if locale == config.DefaultLocale {
			continue
		}
		fmt.Printf("%sProcessing locale %s...%s\n", ColorCyan, locale, ColorReset)
		translations, err := processLocaleFiles(locale, config.MessagesDir)
		if err != nil {
			log.Printf("%sSkipping locale %s due to error: %v%s\n", ColorRed, locale, err, ColorReset)
			continue
		}
		missing := compareTranslationStructure(defaultTranslations, translations, locale, "", logged)
		if len(missing) > 0 {
			missingSummary[locale] = missing
			mergeDeep(translations, missing)
		}
		if err := writeJSONFile(config.TargetDir, locale+".json", translations); err != nil {
			log.Fatal(err)
		}
	}

	if len(missingSummary) > 0 {
		fmt.Printf("%sSummary of missing keys:%s\n", ColorYellow, ColorReset)
		for locale, keys := range missingSummary {
			fmt.Printf("Locale %s missing keys: %+v\n", locale, keys)
		}
		if config.Strict {
			log.Fatal("Strict mode: missing keys detected, aborting.")
		}
	}
}
