package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type repository struct {
	params *params
	dir    string
}

func newRepository(params *params, dir string) *repository {
	return &repository{
		params: params,
		dir:    dir,
	}
}

func (r repository) clone(schema string) error {
	cmd := exec.Command(
		"git",
		"clone",
		"--depth=1",
		"-b", r.params.version,
		r.cloneURL(schema),
	)
	cmd.Dir = r.dir

	err := cmd.Start()
	if err != nil {
		return err
	}

	// kill process if token is invalid (wait password)
	var timer *time.Timer
	timer = time.AfterFunc(30*time.Second, func() {
		cmd.Process.Kill()
	})
	err = cmd.Wait()
	if err != nil {
		return err
	}
	timer.Stop()
	return err
}

func (r repository) diffArchive(dest, typ string) error {
	diff := r.diff()
	if len(diff) == 0 {
		return fmt.Errorf("can't find artifacts")
	}
	// TODO use native zip, tar, gzip package.
	switch typ {
	case "tar.gz":
		return r.targz(diff, dest+"."+typ)
	default:
		return nil
	}
}

func (r repository) targz(src []string, dest string) error {
	params := append(append([]string{}, "czf", dest), src...)
	cmd := exec.Command("tar", params...)
	cmd.Dir = r.pwd()
	return cmd.Run()
}

func (r repository) cleanWithDryRun() ([]byte, error) {
	cmd := exec.Command(
		"git",
		"clean",
		"--dry-run",
	)
	cmd.Dir = r.pwd()
	return cmd.Output()
}

func (r repository) diff() []string {
	out, _ := r.cleanWithDryRun()
	added := strings.Split(string(out), "\n")
	var files []string
	for _, a := range added {
		if a == "" || a == "Would remove Makefile" {
			continue
		}
		files = append(files, strings.TrimPrefix(a, "Would remove "))
	}
	return files
}

func (r repository) path() string {
	return filepath.Join(r.params.remote, r.params.owner(), r.params.repo)
}

func (r repository) pwd() string {
	return filepath.Join(r.dir, r.params.repo)
}

func (r repository) listRemote() error {
	cmd := exec.Command(
		"git",
		"ls-remote",
		"-h",
		r.cloneURL("https"),
		"HEAD",
	)
	cmd.Dir = r.dir

	err := cmd.Start()
	if err != nil {
		return err
	}

	// kill process if token is invalid (wait password)
	var timer *time.Timer
	timer = time.AfterFunc(5*time.Second, func() {
		cmd.Process.Kill()
	})
	err = cmd.Wait()
	if err != nil {
		return err
	}
	timer.Stop()
	return err
}

func (r repository) cloneURL(schema string) string {
	switch schema {
	case "https":
		token := r.params.token
		if token != "" {
			token = token + "@"
		}
		return fmt.Sprintf("https://%s%s/%s/%s.git", token, r.params.remote, r.params.owner(), r.params.repo)
	default:
		return ""
	}
}
