package find

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

type Find struct {
	allowPaths []string
	allowTerms []string
	blockPaths []string
	blockTerms []string
	debug      bool
	fileType   string
	pattern    []string
}

func New(debug bool) *Find {
	return &Find{
		debug: debug,
	}
}

func (f *Find) AddAllowPaths(paths ...string) {
	if len(paths) <= 0 {
		return
	}

	f.allowPaths = append(f.allowPaths, paths...)
	f.generatePattern()
}

func (f *Find) AddAllowTerms(terms ...string) {
	if len(terms) <= 0 {
		return
	}

	f.allowTerms = append(f.allowTerms, terms...)
	f.generatePattern()
}

func (f *Find) AddBlockPaths(paths ...string) {
	if len(paths) <= 0 {
		return
	}

	f.blockPaths = append(f.blockPaths, paths...)
	f.generatePattern()
}

func (f *Find) AddBlockTerms(terms ...string) {
	if len(terms) <= 0 {
		return
	}

	f.blockTerms = append(f.blockTerms, terms...)
	f.generatePattern()
}

func (f *Find) Check(path string) ([]string, error) {
	unique := make(map[string]string)

	defaults := []string{path}
	defaults = append(defaults, f.pattern...)

	for _, path := range f.blockPaths {
		if f.debug {
			log.Printf("Blocking path: %s\n", path)
		}

		parameters := make([]string, len(defaults))
		copy(parameters, defaults)
		parameters = append(parameters, "-iwholename", path)

		err := f.performCheck(unique, parameters...)
		if err != nil {
			return []string{}, err
		}
	}

	for _, term := range f.blockTerms {
		if f.debug {
			log.Printf("Blocking term: %s\n", term)
		}

		parameters := make([]string, len(defaults))
		copy(parameters, defaults)
		parameters = append(parameters, "-iname", term)

		err := f.performCheck(unique, parameters...)
		if err != nil {
			return []string{}, err
		}
	}

	// If there are no block paths or types, but there is a file type, perform check
	if len(f.blockPaths) <= 0 && len(f.blockTerms) <= 0 && f.fileType != "" {
		err := f.performCheck(unique, defaults...)
		if err != nil {
			return []string{}, err
		}
	}

	found := make([]string, 0)

	if len(unique) > 0 {
		for _, item := range unique {
			found = append(found, item)
		}
	}

	return found, nil
}

func (f *Find) SetType(fileType string) {
	if fileType == "" {
		return
	}

	f.fileType = fileType
	f.generatePattern()
}

func (f *Find) generatePattern() {
	if f.debug {
		log.Println("Regenerating config...")
	}

	f.pattern = make([]string, 0)

	if f.fileType != "" {
		if f.debug {
			log.Printf("Setting type: %s\n", f.fileType)
		}

		f.pattern = append(f.pattern, "-type", f.fileType)
	}

	for _, path := range f.allowPaths {
		if f.debug {
			log.Printf("Allowing path: %s\n", path)
		}

		f.pattern = append(f.pattern, "-not", "-iwholename", path)
	}

	for _, term := range f.allowTerms {
		if f.debug {
			log.Printf("Allowing term: %s\n", term)
		}

		f.pattern = append(f.pattern, "-not", "-iname", term)
	}
}

func (f *Find) performCheck(unique map[string]string, parameters ...string) error {
	if f.debug {
		log.Printf("Running: /usr/bin/find %s\n\n", parameters)
	}

	cmd := exec.Command("/usr/bin/find", parameters...)

	found := make([]string, 0)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdoutpipe: %v", err)
	}

	scanner := bufio.NewScanner(stdout)

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to start find: %v", err)
	}

	for scanner.Scan() {
		found = append(found, scanner.Text())
	}

	if scanner.Err() != nil {
		err = cmd.Process.Kill()
		if err != nil {
			return fmt.Errorf("unable to kill scanner: %v", err)
		}

		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("unable to run find: %v", err)
		}

		return fmt.Errorf("unable to run scanner: %v", scanner.Err())
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("unable to run find: %v", err)
	}

	if len(found) > 0 {
		for _, path := range found {
			unique[path] = path
		}
	}

	return nil
}
