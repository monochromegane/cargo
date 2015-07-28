package cargo

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type repository struct {
	remote  string
	org     string
	user    string
	repo    string
	version string
	dir     string
}

func newRepository(remote, org, user, repo, version, dir string) *repository {
	return &repository{
		remote:  remote,
		org:     org,
		user:    user,
		repo:    repo,
		version: version,
		dir:     dir,
	}
}

func (r repository) clone(schema string) error {
	cmd := exec.Command(
		"git",
		"clone",
		"--depth=1",
		"-b", r.version,
		r.cloneURL(schema),
	)
	cmd.Dir = r.dir
	return cmd.Run()
}

func (r repository) diffArchive(typ, dest string) error {
	// TODO use native zip, tar, gzip package.
	switch typ {
	case "tar.gz":
		return r.targz(r.diff(), dest)
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
		if a == "" {
			continue
		}
		files = append(files, strings.TrimPrefix(a, "Would remove "))
	}
	return files
}

func (r repository) path() string {
	return filepath.Join(r.remote, r.owner(), r.repo)
}

func (r repository) pwd() string {
	return filepath.Join(r.dir, r.repo)
}

func (r repository) owner() string {
	var owner string
	if r.org != "" {
		owner = r.org + "/"
	}
	return owner + r.user
}

func (r repository) cloneURL(schema string) string {
	switch schema {
	case "https":
		return fmt.Sprintf("https://%s/%s/%s.git", r.remote, r.owner(), r.repo)
	default:
		return ""
	}
}