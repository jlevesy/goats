package instruction

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/jlevesy/goats/pkg/goats"
	"github.com/jlevesy/goats/pkg/testing"
)

// Builder builds an instruction.
type Builder func(cmd []string) (func(ctx context.Context, t *testing.T), error)

// Builders is a collection of builder mapped per instruction name.
type Builders map[string]Builder

// Resolve returns a raw cmd to a goats.Instruction
func (b Builders) Resolve(cmd []string) (goats.Instruction, error) {
	instructionName := cmd[0]

	builder, ok := b[instructionName]
	if !ok {
		return NewExec(cmd).Exec, nil
	}

	inst, err := builder(cmd)
	if err != nil {
		return nil, fmt.Errorf("unable to create instruction from command %q: %w", instructionName, err)
	}

	return inst, nil
}

// Symbols is the goats symbols exported in Yaegi.
var Symbols = map[string]map[string]reflect.Value{
	"github.com/jlevesy/goats/pkg/instruction": {
		"GetExecOutput": reflect.ValueOf(GetExecOutput),
	},
	"github.com/jlevesy/goats/pkg/testing": {
		"T": reflect.ValueOf((*testing.T)(nil)),
	},
}

type sourceFile struct {
	Name    string
	Content []byte
}

// LoadDynamic loads dynamic instructions from source.
func LoadDynamic(importPaths []string, builders Builders) error {
	files, err := ListSourceFiles(importPaths)
	if err != nil {
		return fmt.Errorf("unable to list source files: %w", err)
	}

	var (
		sources []sourceFile
		tags    []Tag
	)

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("unable to read source file %q: %w", file, err)
		}

		sources = append(sources, sourceFile{Name: file, Content: content})

		fileTags, err := ParseTags(content)
		if err != nil {
			return fmt.Errorf("unable to parse instruction tags for %q, %w", file, err)
		}

		tags = append(tags, fileTags...)
	}

	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)
	i.Use(Symbols)

	for _, source := range sources {
		if _, err = i.Eval(string(source.Content)); err != nil {
			return fmt.Errorf("unable to parse %q: %w", source.Name, err)
		}
	}

	for _, tag := range tags {
		if _, exist := builders[tag.Name]; exist {
			return fmt.Errorf("instruction %q already exists", tag.Name)
		}

		v, err := i.Eval(tag.BuilderName())
		if err != nil {
			return fmt.Errorf("unable to eval function %q: %w", tag.BuilderName(), err)
		}

		builder, ok := v.Interface().(func([]string) (func(ctx context.Context, t *testing.T), error))
		if !ok {
			return fmt.Errorf("function %q is not an instruction builder", tag.BuilderName())
		}

		builders[tag.Name] = builder
	}

	return nil
}

// ListSourceFiles returns a set of all imported source files.
func ListSourceFiles(importPaths []string) ([]string, error) {
	fileSet := make(map[string]struct{})

	for _, path := range importPaths {
		path = filepath.Clean(path)

		stat, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("unable to stat import path: %w", err)
		}

		if !stat.IsDir() {
			return nil, fmt.Errorf("invalid import path %q: must be an existing directory", path)
		}

		if !strings.HasSuffix(path, string(filepath.Separator)) {
			path += string(filepath.Separator)
		}

		// TODO handle recursive discovery ?
		fs, err := filepath.Glob(path + "*.go")
		if err != nil {
			return nil, fmt.Errorf("unable to lookup source files: %w", err)
		}

		for _, file := range fs {
			fileSet[file] = struct{}{}
		}
	}

	var files []string

	for file := range fileSet {
		files = append(files, file)
	}

	sort.Strings(files)

	return files, nil
}

// Tag represents metadata needed to discover instruction builders in go files.
type Tag struct {
	Name    string
	Builder string
	Package string
}

// BuilderName returns the name of the builder function.
func (t Tag) BuilderName() string {
	return fmt.Sprintf("%s.%s", t.Package, t.Builder)
}

var (
	tagRegexp     = regexp.MustCompile(`@instruction{name=(\w+),builder=(\w+)}`)
	packageRegexp = regexp.MustCompile(`package (\w+)`)
)

// ParseTags returns all instruction tags in given content.
func ParseTags(content []byte) ([]Tag, error) {
	p := packageRegexp.FindSubmatch(content)
	if p == nil {
		return nil, errors.New("invalid go file: package name not found")
	}

	packageName := string(p[1])

	var tags []Tag

	for _, match := range tagRegexp.FindAllSubmatch(content, -1) {
		tags = append(tags, Tag{Name: string(match[1]), Builder: string(match[2]), Package: packageName})
	}

	return tags, nil
}
