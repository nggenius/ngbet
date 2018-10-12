// Copyright 2017 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renders

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

var (
	cache               []*namedTemplate
	regularTemplateDefs []string
	lock                sync.Mutex
)

func getReDefineTag(delimsLeft string) *regexp.Regexp {
	return regexp.MustCompile(delimsLeft + "[ ]*define[ ]+\"([^\"]+)\"")
}

func getReTemplateTag(delimsLeft string) *regexp.Regexp {
	return regexp.MustCompile(delimsLeft + "[ ]*template[ ]+\"([^\"]+)\"")
}

type namedTemplate struct {
	Name string
	Src  string
}

// Load prepares and parses all templates from the passed basePath
func Load(opt Options) (map[string]*template.Template, error) {
	lock.Lock()
	defer lock.Unlock()

	// TODO: check if opt.Directories is a dir not a file

	l := &loader{
		Options:       opt,
		templates:     make(map[string]*template.Template),
		reTemplateTag: getReTemplateTag(opt.DelimsLeft),
		reDefineTag:   getReDefineTag(opt.DelimsLeft),
	}
	err := l.loadDir("")
	if err != nil {
		return nil, err
	}

	return l.templates, nil
}

// LoadWithFuncMap prepares and parses all templates from the passed basePath and injects
// a custom template.FuncMap into each template
func LoadWithFuncMap(opt Options) (map[string]*template.Template, error) {
	return Load(opt)
}

func alignTmplName(name string) string {
	name = strings.Replace(name, "\\\\", "/", -1)
	name = strings.Replace(name, "\\", "/", -1)
	return name
}

type loader struct {
	Options
	templates     map[string]*template.Template
	reTemplateTag *regexp.Regexp
	reDefineTag   *regexp.Regexp
}

func (l *loader) loadFile(dir, fileName string) error {
	if !l.IsExtMatch(fileName) {
		return nil
	}

	defer func() {
		cache = cache[0:0]
	}()

	rPath := path.Join(dir, fileName)
	if err := add(l.FileSystem, rPath, l.reTemplateTag); err != nil {
		return err
	}

	// Now we find all regular template definitions and check for the most recent definiton
	for _, t := range regularTemplateDefs {
		found := false
		defineIdx := 0

		// From the beginning (which should) most specifc we look for definitions
		for _, nt := range cache {
			nt.Src = l.reDefineTag.ReplaceAllStringFunc(nt.Src, func(raw string) string {
				parsed := l.reDefineTag.FindStringSubmatch(raw)
				name := parsed[1]
				if name != t {
					return raw
				}
				// Don't touch the first definition
				if !found {
					found = true
					return raw
				}

				defineIdx++

				return fmt.Sprintf(l.DelimsLeft+" define \"%s_invalidated_#%d\" "+l.DelimsRight, name, defineIdx)
			})
		}
	}

	var baseTmpl *template.Template

	for i, nt := range cache {
		var currentTmpl *template.Template
		if i == 0 {
			currentTmpl = template.New(nt.Name).Delims(l.DelimsLeft, l.DelimsRight)
			baseTmpl = currentTmpl
		} else {
			currentTmpl = baseTmpl.New(nt.Name).Delims(l.DelimsLeft, l.DelimsRight)
		}
		if len(l.Funcs) > 0 {
			template.Must(currentTmpl.Funcs(l.Funcs).Parse(nt.Src))
		} else {
			template.Must(currentTmpl.Parse(nt.Src))
		}
	}

	l.templates[alignTmplName(rPath)] = baseTmpl
	return nil
}

func (l *loader) loadDir(dir string) error {
	if dir == ".." {
		return nil
	}

	d, err := l.FileSystem.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	files, err := d.Readdir(0)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			err = l.loadDir(path.Join(dir, file.Name()))
			if err != nil {
				return err
			}
		} else if err = l.loadFile(dir, file.Name()); err != nil {
			return err
		}
	}

	return err
}

func add(fs http.FileSystem, path string, reTemplateTag *regexp.Regexp) error {
	// Get file content
	tplSrc, err := fileContent(path, fs)
	if err != nil {
		return err
	}

	tplName := alignTmplName(path)

	// Make sure template is not already included
	alreadyIncluded := false
	for _, nt := range cache {
		if nt.Name == tplName {
			alreadyIncluded = true
			break
		}
	}
	if alreadyIncluded {
		return nil
	}

	// Add to the cache
	nt := &namedTemplate{
		Name: tplName,
		Src:  tplSrc,
	}
	cache = append(cache, nt)

	// Check for any template block
	for _, raw := range reTemplateTag.FindAllString(nt.Src, -1) {
		parsed := reTemplateTag.FindStringSubmatch(raw)
		templatePath := parsed[1]
		ext := filepath.Ext(templatePath)
		if !strings.Contains(templatePath, ext) {
			regularTemplateDefs = append(regularTemplateDefs, templatePath)
			continue
		}

		// Add this template and continue looking for more template blocks
		add(fs, templatePath, reTemplateTag)
	}

	return nil
}

func isNil(a interface{}) bool {
	if a == nil {
		return true
	}
	aa := reflect.ValueOf(a)
	return !aa.IsValid() || (aa.Type().Kind() == reflect.Ptr && aa.IsNil())
}

func fileContent(path string, fs http.FileSystem) (string, error) {
	// Read the file content of the template
	file, err := fs.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	s := string(b)

	if len(s) < 1 {
		return "", errors.New("render: template file is empty")
	}

	return s, nil
}
