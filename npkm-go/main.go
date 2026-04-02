package main

import (
	"archive/zip"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v3"
)

type Playbook struct {
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Name       string      `yaml:"name"`
	GetUrl     *GetUrl     `yaml:"get_url,omitempty"`
	Copy       *Copy       `yaml:"copy,omitempty"`
	LineInFile *LineInFile `yaml:"lineinfile,omitempty"`
	Command    *Command    `yaml:"command,omitempty"`
	Shell      *Shell      `yaml:"shell,omitempty"`
	File       *File       `yaml:"file,omitempty"`
	Systemd    *Systemd    `yaml:"systemd,omitempty"`
	Git        *Git        `yaml:"git,omitempty"`
	Remove     *Remove     `yaml:"remove,omitempty"`
	Debug      *Debug      `yaml:"debug,omitempty"`
	Replace    *Replace    `yaml:"replace,omitempty"`
	Fail       *Fail       `yaml:"fail,omitempty"`
	Unzip      *Unzip      `yaml:"unzip,omitempty"`
	Move       *Move       `yaml:"move,omitempty"`
	Path       *PathTask   `yaml:"path,omitempty"`
}

type GetUrl struct {
	Url  string `yaml:"url"`
	Dest string `yaml:"dest"`
}

type Copy struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

type Move struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

type PathTask struct {
	Path string `yaml:"path"`
}

type LineInFile struct {
	Path   string `yaml:"path"`
	Regexp string `yaml:"regexp,omitempty"`
	Line   string `yaml:"line"`
}

type Command struct {
	Cmd string `yaml:"cmd"`
	Cwd string `yaml:"cwd,omitempty"`
}

type Shell struct {
	Cmd string `yaml:"cmd"`
	Cwd string `yaml:"cwd,omitempty"`
}

type File struct {
	Path  string      `yaml:"path"`
	State string      `yaml:"state"` // directory, touch, link, absent
	Src   string      `yaml:"src,omitempty"`
	Mode  os.FileMode `yaml:"mode,omitempty"`
}

type Systemd struct {
	Name    string `yaml:"name"`
	State   string `yaml:"state"` // started, stopped, restarted
	Enabled bool   `yaml:"enabled"`
}

type Git struct {
	Repo string `yaml:"repo"`
	Dest string `yaml:"dest"`
}

type Remove struct {
	Path string `yaml:"path"`
}

type Debug struct {
	Msg string `yaml:"msg"`
}

type Replace struct {
	Path    string `yaml:"path"`
	Regexp  string `yaml:"regexp"`
	Replace string `yaml:"replace"`
}

type Fail struct {
	Msg string `yaml:"msg"`
}

type Unzip struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <playbook.yml | http(s)://... | git repo>\n", os.Args[0])
		os.Exit(1)
	}

	source := os.Args[1]
	var data []byte
	var err error

	isGit := strings.HasSuffix(source, ".git") || strings.HasPrefix(source, "git://") || strings.HasPrefix(source, "git@")
	if isGit {
		tempDir, err := os.MkdirTemp("", "npkm-repo-*")
		if err != nil {
			fmt.Printf("Error creating temp dir: %v\n", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tempDir)

		fmt.Printf("Cloning %s into temporary directory...\n", source)
		_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
			URL: source,
		})
		if err != nil {
			fmt.Printf("Error cloning git repo: %v\n", err)
			os.Exit(1)
		}

		playbookPath := filepath.Join(tempDir, "playbook.yml")
		if _, err := os.Stat(playbookPath); os.IsNotExist(err) {
			playbookPath = filepath.Join(tempDir, "playbook.yaml")
		}

		data, err = os.ReadFile(playbookPath)
		if err != nil {
			fmt.Printf("Error reading playbook in git repo: %v\n", err)
			os.Exit(1)
		}

		os.Chdir(tempDir)
	} else if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		fmt.Printf("Downloading playbook from %s...\n", source)
		client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
		resp, err := client.Get(source)
		if err != nil {
			fmt.Printf("Error downloading playbook: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to download playbook, status: %s\n", resp.Status)
			os.Exit(1)
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading playbook response: %v\n", err)
			os.Exit(1)
		}
	} else {
		data, err = os.ReadFile(source)
		if err != nil {
			fmt.Printf("Error reading playbook: %v\n", err)
			os.Exit(1)
		}
	}

	var playbook Playbook
	if err := yaml.Unmarshal(data, &playbook); err != nil {
		fmt.Printf("Error parsing yaml: %v\n", err)
		os.Exit(1)
	}

	for _, task := range playbook.Tasks {
		fmt.Printf("TASK [%s]\n", task.Name)
		var err error
		if task.GetUrl != nil {
			err = executeGetUrl(task.GetUrl)
		} else if task.Copy != nil {
			err = executeCopy(task.Copy)
		} else if task.LineInFile != nil {
			err = executeLineInFile(task.LineInFile)
		} else if task.Command != nil {
			err = executeCommand(task.Command)
		} else if task.Shell != nil {
			err = executeShell(task.Shell)
		} else if task.File != nil {
			err = executeFile(task.File)
		} else if task.Systemd != nil {
			err = executeSystemd(task.Systemd)
		} else if task.Git != nil {
			err = executeGit(task.Git)
		} else if task.Remove != nil {
			err = executeRemove(task.Remove)
		} else if task.Debug != nil {
			executeDebug(task.Debug)
		} else if task.Replace != nil {
			err = executeReplace(task.Replace)
		} else if task.Fail != nil {
			err = fmt.Errorf("%s", task.Fail.Msg)
		} else if task.Unzip != nil {
			err = executeUnzip(task.Unzip)
		} else if task.Move != nil {
			err = executeMove(task.Move)
		} else if task.Path != nil {
			err = executePath(task.Path)
		} else {
			fmt.Println("  warning: unknown or missing module type")
			continue
		}

		if err != nil {
			fmt.Printf("  fatal: [%s] %v\n", task.Name, err)
			os.Exit(1)
		} else {
			fmt.Printf("  changed\n\n")
		}
	}
}

func executeGetUrl(spec *GetUrl) error {
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Get(spec.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	if err := os.MkdirAll(filepath.Dir(spec.Dest), 0755); err != nil {
		return err
	}

	out, err := os.Create(spec.Dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func executeCopy(spec *Copy) error {
	in, err := os.Open(spec.Src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(spec.Dest), 0755); err != nil {
		return err
	}

	out, err := os.Create(spec.Dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func executeLineInFile(spec *LineInFile) error {
	content, err := os.ReadFile(spec.Path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	replaced := false
	if spec.Regexp != "" {
		re, err := regexp.Compile(spec.Regexp)
		if err != nil {
			return fmt.Errorf("invalid regexp: %v", err)
		}
		for i, line := range lines {
			if re.MatchString(line) {
				lines[i] = spec.Line
				replaced = true
				break
			}
		}
	} else {
		for _, line := range lines {
			if line == spec.Line {
				replaced = true
				break
			}
		}
	}

	if !replaced {
		lines = append(lines, spec.Line)
	}

	finalContent := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(spec.Path, []byte(finalContent), 0644)
}

func executeCommand(spec *Command) error {
	parts := strings.Fields(spec.Cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = spec.Cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeShell(spec *Shell) error {
	cmd := exec.Command("sh", "-c", spec.Cmd)
	cmd.Dir = spec.Cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeFile(spec *File) error {
	switch spec.State {
	case "directory":
		if err := os.MkdirAll(spec.Path, 0755); err != nil {
			return err
		}
	case "touch":
		if err := os.MkdirAll(filepath.Dir(spec.Path), 0755); err != nil {
			return err
		}
		f, err := os.OpenFile(spec.Path, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		f.Close()
		currentTime := time.Now()
		if err := os.Chtimes(spec.Path, currentTime, currentTime); err != nil {
			return err
		}
	case "link":
		_ = os.Remove(spec.Path)
		if err := os.Symlink(spec.Src, spec.Path); err != nil {
			return err
		}
	case "absent":
		return os.RemoveAll(spec.Path)
	default:
		return fmt.Errorf("unknown file state: %s", spec.State)
	}

	if spec.Mode != 0 {
		if err := os.Chmod(spec.Path, spec.Mode); err != nil {
			return err
		}
	}
	return nil
}

func executeSystemd(spec *Systemd) error {
	if spec.Enabled {
		cmd := exec.Command("systemctl", "enable", spec.Name)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to enable: %v", err)
		}
	}
	if spec.State != "" {
		allowed := map[string]string{
			"started":   "start",
			"stopped":   "stop",
			"restarted": "restart",
		}
		action, ok := allowed[spec.State]
		if !ok {
			return fmt.Errorf("unknown systemd state: %s", spec.State)
		}
		cmd := exec.Command("systemctl", action, spec.Name)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to %s: %v", action, err)
		}
	}
	return nil
}

func executeGit(spec *Git) error {
	if _, err := os.Stat(filepath.Join(spec.Dest, ".git")); err == nil {
		repo, err := git.PlainOpen(spec.Dest)
		if err != nil {
			return err
		}
		w, err := repo.Worktree()
		if err != nil {
			return err
		}
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}
		return nil
	}

	_, err := git.PlainClone(spec.Dest, false, &git.CloneOptions{
		URL: spec.Repo,
	})
	return err
}

func executeRemove(spec *Remove) error {
	return os.RemoveAll(spec.Path)
}

func executeDebug(spec *Debug) {
	fmt.Printf("  msg: %s\n", spec.Msg)
}

func executeReplace(spec *Replace) error {
	content, err := os.ReadFile(spec.Path)
	if err != nil {
		return err
	}
	re, err := regexp.Compile(spec.Regexp)
	if err != nil {
		return fmt.Errorf("invalid regexp: %v", err)
	}
	newContent := re.ReplaceAll(content, []byte(spec.Replace))
	return os.WriteFile(spec.Path, newContent, 0644)
}

func executeUnzip(spec *Unzip) error {
	r, err := zip.OpenReader(spec.Src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(spec.Dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(spec.Dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func executeMove(spec *Move) error {
	if err := os.MkdirAll(filepath.Dir(spec.Dest), 0755); err != nil {
		return err
	}
	
	err := os.Rename(spec.Src, spec.Dest)
	if err == nil {
		return nil
	}

	// Fallback for cross-device link errors
	in, err := os.Open(spec.Src)
	if err != nil {
		return err
	}
	out, err := os.Create(spec.Dest)
	if err != nil {
		in.Close()
		return err
	}
	_, err = io.Copy(out, in)
	in.Close()
	out.Close()
	
	if err != nil {
		return err
	}
	
	return os.RemoveAll(spec.Src)
}

func executePath(spec *PathTask) error {
	newPath := spec.Path

	if runtime.GOOS == "windows" {
		// Option 1: Try PowerShell (often available, safe string handling)
		psCmd := fmt.Sprintf(`$oldPath = [Environment]::GetEnvironmentVariable('Path', 'User'); if (($oldPath -split ';') -notcontains '%s') { [Environment]::SetEnvironmentVariable('Path', $oldPath + ';%s', 'User') }`, newPath, newPath)
		if err := exec.Command("powershell", "-NoProfile", "-Command", psCmd).Run(); err == nil {
			return nil
		}

		// Option 2: Fallback to reg.exe (built-in Windows utility, available even without PowerShell)
		out, err := exec.Command("reg", "query", `HKCU\Environment`, "/v", "PATH").Output()
		if err == nil {
			outStr := string(out)
			if !strings.Contains(outStr, newPath) {
				var currentPath string
				lines := strings.Split(outStr, "\n")
				for _, line := range lines {
					if strings.Contains(line, "PATH") && (strings.Contains(line, "REG_SZ") || strings.Contains(line, "REG_EXPAND_SZ")) {
						parts := strings.Fields(line)
						if len(parts) >= 3 {
							idx := strings.Index(line, parts[1]) + len(parts[1])
							currentPath = strings.TrimSpace(line[idx:])
						}
					}
				}
				newFullPath := newPath
				if currentPath != "" {
					newFullPath = currentPath + ";" + newPath
				}
				if errAdd := exec.Command("reg", "add", `HKCU\Environment`, "/v", "PATH", "/t", "REG_EXPAND_SZ", "/d", newFullPath, "/f").Run(); errAdd == nil {
					return nil
				}
			} else {
				return nil // Already in path
			}
		}

		return fmt.Errorf("failed to update Windows PATH using both PowerShell and reg.exe")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	exportLine := fmt.Sprintf(`export PATH="%s:$PATH"`, newPath)
	filesToUpdate := []string{".bashrc", ".zshrc", ".profile", ".bash_profile"}
	
	updated := false
	for _, file := range filesToUpdate {
		rcPath := filepath.Join(home, file)
		if _, err := os.Stat(rcPath); err == nil {
			content, err := os.ReadFile(rcPath)
			if err == nil && !strings.Contains(string(content), exportLine) {
				f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_WRONLY, 0644)
				if err == nil {
					f.WriteString("\n" + exportLine + "\n")
					f.Close()
					updated = true
				}
			}
		}
	}

	if !updated {
		rcPath := filepath.Join(home, ".bashrc")
		if _, err := os.Stat(rcPath); os.IsNotExist(err) {
			os.WriteFile(rcPath, []byte(exportLine+"\n"), 0644)
		}
	}
	
	return nil
}
