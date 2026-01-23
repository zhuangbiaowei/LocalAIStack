package info

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/zhuangbiaowei/LocalAIStack/internal/config"
)

type BaseInfo struct {
	CollectedAt         string   `json:"collected_at"`
	OS                  string   `json:"os"`
	Arch                string   `json:"arch"`
	Kernel              string   `json:"kernel"`
	CPUModel            string   `json:"cpu_model"`
	CPUCores            int      `json:"cpu_cores"`
	MemoryTotal         string   `json:"memory_total"`
	GPU                 string   `json:"gpu"`
	DiskTotal           string   `json:"disk_total"`
	DiskAvailable       string   `json:"disk_available"`
	Hostname            string   `json:"hostname"`
	InternalIPs         []string `json:"internal_ips"`
	Docker              string   `json:"docker"`
	Podman              string   `json:"podman"`
	LocalAIStackVersion string   `json:"localaistack_version"`
	RuntimeCapabilities string   `json:"runtime_capabilities"`
}

func CollectBaseInfo(ctx context.Context) BaseInfo {
	if ctx == nil {
		ctx = context.Background()
	}
	info := BaseInfo{
		CollectedAt:         time.Now().Format(time.RFC3339),
		OS:                  runtime.GOOS,
		Arch:                runtime.GOARCH,
		CPUCores:            runtime.NumCPU(),
		LocalAIStackVersion: config.Version,
	}

	info.Kernel = kernelInfo(ctx)
	info.CPUModel = cpuModel(ctx)
	info.MemoryTotal = memoryTotal(ctx)
	info.GPU = gpuInfo(ctx)
	info.DiskTotal, info.DiskAvailable = diskInfo()
	info.Hostname = hostname()
	info.InternalIPs = internalIPs()
	info.Docker = runtimeAvailability(ctx, "docker")
	info.Podman = runtimeAvailability(ctx, "podman")
	info.RuntimeCapabilities = runtimeCapabilities(info.Docker, info.Podman)

	return info
}

func kernelInfo(ctx context.Context) string {
	stdout, stderr, err := runCommand(ctx, "uname", "-a")
	if err != nil {
		return formatUnknown("uname", err, stderr)
	}
	return strings.TrimSpace(stdout)
}

func cpuModel(ctx context.Context) string {
	switch runtime.GOOS {
	case "linux":
		file, err := os.Open("/proc/cpuinfo")
		if err != nil {
			return formatUnknown("/proc/cpuinfo", err, "")
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(strings.ToLower(line), "model name") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return formatUnknown("/proc/cpuinfo", err, "")
		}
		return "unknown: cpu model not found"
	case "darwin":
		stdout, stderr, err := runCommand(ctx, "sysctl", "-n", "machdep.cpu.brand_string")
		if err != nil {
			return formatUnknown("sysctl", err, stderr)
		}
		return strings.TrimSpace(stdout)
	default:
		return "unknown: unsupported OS"
	}
}

func memoryTotal(ctx context.Context) string {
	switch runtime.GOOS {
	case "linux":
		file, err := os.Open("/proc/meminfo")
		if err != nil {
			return formatUnknown("/proc/meminfo", err, "")
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					value := fields[1]
					return fmt.Sprintf("%s kB", value)
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return formatUnknown("/proc/meminfo", err, "")
		}
		return "unknown: meminfo missing MemTotal"
	case "darwin":
		stdout, stderr, err := runCommand(ctx, "sysctl", "-n", "hw.memsize")
		if err != nil {
			return formatUnknown("sysctl", err, stderr)
		}
		return fmt.Sprintf("%s bytes", strings.TrimSpace(stdout))
	default:
		return "unknown: unsupported OS"
	}
}

func gpuInfo(ctx context.Context) string {
	switch runtime.GOOS {
	case "linux":
		stdout, stderr, err := runCommand(ctx, "nvidia-smi", "--query-gpu=name", "--format=csv,noheader")
		if err == nil {
			gpu := strings.TrimSpace(stdout)
			if gpu != "" {
				return gpu
			}
		}

		stdout, stderr, err = runCommand(ctx, "lspci")
		if err != nil {
			return formatUnknown("lspci", err, stderr)
		}
		var matches []string
		for _, line := range strings.Split(stdout, "\n") {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "vga") || strings.Contains(lower, "3d") {
				matches = append(matches, strings.TrimSpace(line))
			}
		}
		if len(matches) == 0 {
			return "unknown: no GPU entries"
		}
		return strings.Join(matches, "; ")
	default:
		return "unknown: unsupported OS"
	}
}

func diskInfo() (string, string) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return formatUnknown("statfs", err, ""), formatUnknown("statfs", err, "")
	}
	total := stat.Blocks * uint64(stat.Bsize)
	available := stat.Bavail * uint64(stat.Bsize)
	return formatBytes(total), formatBytes(available)
}

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return formatUnknown("hostname", err, "")
	}
	return name
}

func internalIPs() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return []string{formatUnknown("interfaces", err, "")}
	}

	var ips []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := extractIP(addr)
			if ip == "" {
				continue
			}
			ips = append(ips, ip)
		}
	}

	if len(ips) == 0 {
		return []string{"unknown: no internal IPs"}
	}
	return ips
}

func runtimeAvailability(ctx context.Context, runtimeName string) string {
	stdout, stderr, err := runCommand(ctx, runtimeName, "version")
	if err != nil {
		return formatUnknown(runtimeName, err, stderr)
	}
	trimmed := strings.TrimSpace(stdout)
	if trimmed == "" {
		return "available"
	}
	return fmt.Sprintf("available: %s", trimmed)
}

func runtimeCapabilities(dockerStatus, podmanStatus string) string {
	return fmt.Sprintf("docker=%s; podman=%s", dockerStatus, podmanStatus)
}

func extractIP(addr net.Addr) string {
	switch v := addr.(type) {
	case *net.IPNet:
		if v.IP == nil {
			return ""
		}
		if v.IP.IsLoopback() {
			return ""
		}
		if v.IP.To4() != nil {
			return v.IP.String()
		}
	case *net.IPAddr:
		if v.IP == nil || v.IP.IsLoopback() {
			return ""
		}
		if v.IP.To4() != nil {
			return v.IP.String()
		}
	}
	return ""
}

func runCommand(ctx context.Context, name string, args ...string) (string, string, error) {
	if errors.Is(ctx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return "", "", ctx.Err()
	}

	cmd := exec.CommandContext(ctx, name, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func formatUnknown(source string, err error, stderr string) string {
	if err == nil {
		return "unknown"
	}
	message := fmt.Sprintf("unknown: %s", err.Error())
	if source != "" {
		message = fmt.Sprintf("unknown (%s): %s", source, err.Error())
	}
	stderr = strings.TrimSpace(stderr)
	if stderr != "" {
		message = fmt.Sprintf("%s; stderr: %s", message, stderr)
	}
	return message
}

func formatBytes(value uint64) string {
	const unit = 1024
	if value < unit {
		return fmt.Sprintf("%d B", value)
	}
	div, exp := uint64(unit), 0
	for n := value / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(value)/float64(div), "KMGTPE"[exp])
}
