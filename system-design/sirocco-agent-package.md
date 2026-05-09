To create a **Debian package (.deb)**, you basically prepare a folder that mimics the Linux filesystem (`/usr/bin`, `/etc/...`) + a `DEBIAN/control` file, then build it with `dpkg-deb`.

---

# 1) Install packaging tools

```bash
sudo apt update
sudo apt install dpkg-dev debhelper
```

---

# 2) Create package directory structure

Example package name: `sirocco-agent`

```bash
mkdir -p sirocco-agent_1.0.0_amd64/DEBIAN
mkdir -p sirocco-agent_1.0.0_amd64/usr/bin
mkdir -p sirocco-agent_1.0.0_amd64/etc/systemd/system
```

This folder name format is common:
`package_version_arch`

---

# 3) Put your binary inside the package

Assume you already built your Go agent:

```bash
go build -o sirocco-agent
```

Copy it:

```bash
cp sirocco-agent sirocco-agent_1.0.0_amd64/usr/bin/
chmod +x sirocco-agent_1.0.0_amd64/usr/bin/sirocco-agent
```

---

# 4) Add systemd service file (optional but realistic)

Create:

```bash
nano sirocco-agent_1.0.0_amd64/etc/systemd/system/sirocco-agent.service
```

Paste:

```ini
[Unit]
Description=Sirocco Agent
After=network.target

[Service]
ExecStart=/usr/bin/sirocco-agent
Restart=always
User=root

[Install]
WantedBy=multi-user.target
```

---

# 5) Create the control file

Create:

```bash
nano sirocco-agent_1.0.0_amd64/DEBIAN/control
```

Paste:

```text
Package: sirocco-agent
Version: 1.0.0
Section: base
Priority: optional
Architecture: amd64
Maintainer: Mohamed Amine Mensi <mensimohamedamine25@gmail.com>
Description: Sirocco Agent for MySQL sharding cluster
```

---

# 6) Add post-install script (recommended)

This enables and starts the service automatically.

Create:

```bash
nano sirocco-agent_1.0.0_amd64/DEBIAN/postinst
```

Paste:

```bash
#!/bin/bash
systemctl daemon-reload
systemctl enable sirocco-agent
systemctl restart sirocco-agent
exit 0
```

Make executable:

```bash
chmod 755 sirocco-agent_1.0.0_amd64/DEBIAN/postinst
```

---

# 7) Build the `.deb` package

```bash
dpkg-deb --build sirocco-agent_1.0.0_amd64
```

This produces:

```text
sirocco-agent_1.0.0_amd64.deb
```

---

# 8) Test installation

```bash
sudo dpkg -i sirocco-agent_1.0.0_amd64.deb
```

Check service:

```bash
systemctl status sirocco-agent
```

---

# 9) Remove package (test uninstall)

```bash
sudo apt remove sirocco-agent
```

or

```bash
sudo dpkg -r sirocco-agent
```

---

## Important notes

- The folder must contain `DEBIAN/control` or build fails.
- `Architecture` must match your machine (`amd64`, `arm64`, etc).
- For real apt install support, you must host it in an APT repository.

If you want, I can give you a **Makefile** that builds the Go binary + builds the `.deb` automatically.
