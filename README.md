# foreman-builder (Lima Edition)

A macOS tool to help build Foreman development environments using Lima VMs, with the goal of making it feel as close to running locally as possible.

This tool sets up a Lima-based Foreman development environment on Rocky Linux 9.

## Requirements

- [Go](https://go.dev/doc/install)
- [Lima](https://lima-vm.io/) - Install with: `brew install lima`

## Architecture Support

This tool supports both:
- **Apple Silicon** (ARM64/aarch64) - M1, M2, M3, M4 Macs
- **Intel Macs** (x86_64/amd64)

Lima will automatically use the appropriate Rocky Linux 9 image for your architecture.

## Getting Started

### Installation

First install foreman-builder:
```bash
go install github.com/akshaykapse/foreman-builder-test@latest
```

Or build from source:
```bash
git clone https://github.com/akshaykapse/foreman-builder-test.git
cd foreman-builder-test
go build -o foreman-builder
# Optionally move to PATH
sudo mv foreman-builder /usr/local/bin/
```

### Commands

**Create a Lima instance**

To create a development instance, run the `create` command. This will guide you through the setup process and create a Lima VM with Foreman installed. Once finished, you can SSH into the environment.

```bash
foreman-builder create
```

After creation, access your instance:
```bash
limactl shell <instance-name>

# Inside the VM, start Foreman:
cd ~/foreman && bundle exec foreman start
```

Access Foreman at: `https://foreman-dev.example.com`  
Default credentials: `admin` / `changeme`

**List instances**

Display a list of Lima instances created by `foreman-builder`:
```bash
foreman-builder list
```

**Delete instance**

Delete an instance by name. You can only delete instances created with `foreman-builder`:
```bash
foreman-builder delete <instance-name>
```

## Configuration

You can customize the packages installed during setup by creating a `config.yml` file:

```yaml
packages:
  - vim
  - tmux
  - neovim
```

Place this file in the directory where you run `foreman-builder create`.

## Lima Features

- **Port Forwarding**: Ports 3000 (Foreman), 5432 (PostgreSQL), and 8443 (HTTPS) are automatically forwarded
- **File Sharing**: Your home directory is mounted in the VM for easy file access
- **SSH Access**: Use `limactl shell <instance-name>` to access the VM
- **Resource Configuration**: 4 CPUs, 8GB RAM, 100GB disk (configurable in the generated YAML)

## Troubleshooting

**Check Lima status:**
```bash
limactl list
```

**View Lima instance logs:**
```bash
limactl show-log <instance-name>
```

**Stop an instance:**
```bash
limactl stop <instance-name>
```

**Restart an instance:**
```bash
limactl start <instance-name>
```

## Contributing

Contributions are welcome! This fork adapts the original OrbStack-based tool to work with Lima.

## License

MIT License - see LICENSE file for details.



