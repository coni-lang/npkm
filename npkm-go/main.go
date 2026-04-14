package main

import (
	"archive/zip"
	"crypto/tls"
	"flag"
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

var Version string = "development"

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
	PowerShell *PowerShell `yaml:"powershell,omitempty"`
	Package    *Package    `yaml:"package,omitempty"`
	Cron       *Cron       `yaml:"cron,omitempty"`
	Archive    *Archive    `yaml:"archive,omitempty"`
	User       *User       `yaml:"user,omitempty"`
	Service    *Service    `yaml:"service,omitempty"`
	Template   *Template   `yaml:"template,omitempty"`
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

type PowerShell struct {
	Inline string   `yaml:"inline,omitempty"`
	File   string   `yaml:"file,omitempty"`
	Params []string `yaml:"params,omitempty"`
	Cwd    string   `yaml:"cwd,omitempty"`
}

type Package struct {
	Name  string `yaml:"name"`
	State string `yaml:"state"` // present, absent
}

type Cron struct {
	Name     string `yaml:"name"`
	Job      string `yaml:"job"`
	Schedule string `yaml:"schedule"` // e.g. "0 2 * * *"
	State    string `yaml:"state"` // present, absent
}

type Archive struct {
	Src    string `yaml:"src"`
	Dest   string `yaml:"dest"`
	Format string `yaml:"format"` // zip, tar
}

type User struct {
	Name  string `yaml:"name"`
	State string `yaml:"state"` // present, absent
}

type Service struct {
	Name    string `yaml:"name"`
	State   string `yaml:"state"` // started, stopped, restarted
	Enabled bool   `yaml:"enabled"`
}

type Template struct {
	Src  string            `yaml:"src"`
	Dest string            `yaml:"dest"`
	Vars map[string]string `yaml:"vars"` // For Go, normal maps work
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
	var versionFlag bool
	var helpFlag bool
	flag.BoolVar(&versionFlag, "v", false, "prints version (compiled at date)")
	flag.BoolVar(&helpFlag, "h", false, "shows help and supported tasks")
	
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <playbook.yml | directory | http(s)://... | git repo>\n\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("\nSupported Playbook Tasks:")
		fmt.Println("  get_url:      Download a file from HTTP/HTTPS.")
		fmt.Println("                { url: string, dest: string }")
		fmt.Println("  copy:         Copy a file from local source to destination.")
		fmt.Println("                { src: string, dest: string }")
		fmt.Println("  lineinfile:   Ensure a particular line is in a file, or replace an existing line using a regular expression.")
		fmt.Println("                { path: string, regexp?: string, line: string }")
		fmt.Println("  command:      Execute a command without going through a shell.")
		fmt.Println("                { cmd: string, cwd?: string }")
		fmt.Println("  shell:        Execute a command through the system shell.")
		fmt.Println("                { cmd: string, cwd?: string }")
		fmt.Println("  file:         Manage files, directories, and symlinks.")
		fmt.Println("                { path: string, state: string, src?: string, mode?: int }")
		fmt.Println("                states: directory, touch, link, absent")
		fmt.Println("  systemd:      Manage systemd services.")
		fmt.Println("                { name: string, state: string, enabled: bool }")
		fmt.Println("                states: started, stopped, restarted")
		fmt.Println("  git:          Clone or pull a git repository.")
		fmt.Println("                { repo: string, dest: string }")
		fmt.Println("  remove:       Remove a file or directory.")
		fmt.Println("                { path: string }")
		fmt.Println("  debug:        Print a message to the console.")
		fmt.Println("                { msg: string }")
		fmt.Println("  replace:      Replace all instances of a regular expression in a file.")
		fmt.Println("                { path: string, regexp: string, replace: string }")
		fmt.Println("  fail:         Fail the playbook execution with a message.")
		fmt.Println("                { msg: string }")
		fmt.Println("  unzip:        Extract a zip archive.")
		fmt.Println("                { src: string, dest: string }")
		fmt.Println("  move:         Move or rename a file or directory.")
		fmt.Println("                { src: string, dest: string }")
		fmt.Println("  path:         Add a directory to the system PATH environment variable.")
		fmt.Println("                { path: string }")
		fmt.Println("  powershell:   Execute a PowerShell script or inline command.")
		fmt.Println("                { inline?: string, file?: string, params?: []string, cwd?: string }")
		fmt.Println("  package:      Manage OS packages.")
		fmt.Println("  cron:         Manage crontab entries.")
		fmt.Println("  archive:      Compress files/directories.")
		fmt.Println("  user:         Manage OS users.")
		fmt.Println("  service:      Manage cross-platform background services.")
		fmt.Println("  template:     Deploy templated files replacing {{ key }} with Map vars.")
		fmt.Println("\nExample Playbook:")
		fmt.Println("  tasks:")
		fmt.Println("    - name: Ensure target directory exists")
		fmt.Println("      file:")
		fmt.Println("        path: /tmp/myapp")
		fmt.Println("        state: directory")
	}

	flag.Parse()

	if versionFlag {
		v := Version
		if v == "development" {
			if stat, err := os.Stat(os.Args[0]); err == nil {
				v = fmt.Sprintf("development (compiled %s)", stat.ModTime().Format(time.RFC3339))
			}
		}
		fmt.Printf("npkm version: %s\n", v)
		os.Exit(0)
	}

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	source := args[0]
	var data []byte
	var err error

	if info, statErr := os.Stat(source); statErr == nil && info.IsDir() {
		entries, err := os.ReadDir(source)
		if err != nil {
			fmt.Printf("Error reading directory: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Available playbooks in %s:\n", source)
		found := false
		for _, entry := range entries {
			if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".yml") || strings.HasSuffix(entry.Name(), ".yaml")) {
				fmt.Printf("  - %s\n", entry.Name())
				found = true
			}
		}
		if !found {
			fmt.Println("  (No .yml or .yaml files found)")
		}
		os.Exit(0)
	}

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
		} else if task.PowerShell != nil {
			err = executePowerShell(task.PowerShell)
		} else if task.Package != nil {
			err = executePackage(task.Package)
		} else if task.Cron != nil {
			err = executeCron(task.Cron)
		} else if task.Archive != nil {
			err = executeArchive(task.Archive)
		} else if task.User != nil {
			err = executeUser(task.User)
		} else if task.Service != nil {
			err = executeService(task.Service)
		} else if task.Template != nil {
			err = executeTemplate(task.Template)
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

func executePowerShell(spec *PowerShell) error {
	psBin := "powershell"
	if runtime.GOOS != "windows" {
		psBin = "pwsh"
	}
	
	args := []string{"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass"}
	if spec.Inline != "" {
		args = append(args, "-Command", spec.Inline)
	} else if spec.File != "" {
		args = append(args, "-File", spec.File)
		args = append(args, spec.Params...)
	} else {
		return fmt.Errorf("powershell task requires either 'inline' or 'file'")
	}

	cmd := exec.Command(psBin, args...)
	cmd.Dir = spec.Cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}


func executePackage(spec *Package) error {
	packages := []string{"brew", "apt-get", "yum", "choco"}
	var pkgCmd string
	for _, p := range packages {
		if err := exec.Command("which", p).Run(); err == nil {
			pkgCmd = p
			break
		}
	}
	if pkgCmd == "" && runtime.GOOS == "windows" {
		pkgCmd = "choco"
	} else if pkgCmd == "" {
		return fmt.Errorf("no supported package manager found")
	}

	installCmd := "install"
	if spec.State == "absent" {
		installCmd = "uninstall"
		if pkgCmd == "apt-get" || pkgCmd == "yum" {
			installCmd = "remove"
		}
	}
	
	args := []string{installCmd}
	if pkgCmd == "apt-get" || pkgCmd == "yum" || pkgCmd == "choco" {
		args = append(args, "-y")
	}
	args = append(args, spec.Name)
	cmd := exec.Command(pkgCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeCron(spec *Cron) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("cron task not yet supported on windows")
	}
	marker := fmt.Sprintf("# NPKM: %s", spec.Name)
	out, _ := exec.Command("crontab", "-l").Output()
	lines := strings.Split(string(out), "\n")
	var newLines []string
	skip := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" { continue }
		if line == marker {
			skip = true
			continue
		}
		if skip {
			skip = false
			continue
		}
		newLines = append(newLines, line)
	}

	if spec.State != "absent" {
		newLines = append(newLines, marker)
		newLines = append(newLines, fmt.Sprintf("%s %s", spec.Schedule, spec.Job))
	}
	newLines = append(newLines, "")
	
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(strings.Join(newLines, "\n"))
	return cmd.Run()
}

func executeArchive(spec *Archive) error {
	format := spec.Format
	if format == "" { format = "tar" }
	var cmd *exec.Cmd
	if format == "zip" {
		cmd = exec.Command("zip", "-r", spec.Dest, filepath.Base(spec.Src))
		cmd.Dir = filepath.Dir(spec.Src)
	} else {
		cmd = exec.Command("tar", "-czf", spec.Dest, "-C", filepath.Dir(spec.Src), filepath.Base(spec.Src))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeUser(spec *User) error {
	goos := runtime.GOOS
	if goos == "windows" {
		if spec.State == "absent" {
			return exec.Command("net", "user", spec.Name, "/delete").Run()
		}
		return exec.Command("net", "user", spec.Name, "/add").Run()
	} else if goos == "darwin" {
		if spec.State == "absent" {
			return exec.Command("sysadminctl", "-deleteUser", spec.Name).Run()
		}
		return exec.Command("sysadminctl", "-addUser", spec.Name).Run()
	} else {
		if spec.State == "absent" {
			return exec.Command("userdel", spec.Name).Run()
		}
		return exec.Command("useradd", spec.Name).Run()
	}
}

func executeService(spec *Service) error {
	goos := runtime.GOOS
	if goos == "windows" {
		action := "start"
		if spec.State == "stopped" { action = "stop" }
		return exec.Command("net", action, spec.Name).Run()
	} else if goos == "darwin" {
		action := "load"
		if spec.State == "stopped" { action = "unload" }
		return exec.Command("launchctl", action, spec.Name).Run()
	} else {
		action := "start"
		if spec.State == "stopped" { action = "stop" }
		if spec.State == "restarted" { action = "restart" }
		return exec.Command("systemctl", action, spec.Name).Run()
	}
}

func executeTemplate(spec *Template) error {
	content, err := os.ReadFile(spec.Src)
	if err != nil { return err }
	res := string(content)
	for k, v := range spec.Vars {
		res = strings.ReplaceAll(res, fmt.Sprintf("{{ %s }}", k), v)
	}
	return os.WriteFile(spec.Dest, []byte(res), 0644)
}
