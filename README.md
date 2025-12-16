# penta

## Session Management

```bash
penta session create --name "my session" --workspace "/tmp"
penta session delete --name "my session"
penta session list
penta session current
penta session use --name "my session"
```

## Recon

### Hosts scan

```bash
penta scan hosts -m arp,icmp,tmp 10.0.0.0/24
penta scan hosts 10.0.0.0-255
penta scan hosts 10.0.0.0,10.0.0.1
penta scan hosts 10.0.0.0-255,10.0.1.1
penta scan hosts 10.0.0.0-255,10.0.1.0/24

penta scan --nmap -- -sn 10.0.0.0/24
penta scan 10.0.0.0/24 --nmap -- -sn 
```

### Port scan

```bash
penta scan ports -p 22,80
penta scan ports -p 1-1000
penta scan ports -p 22,80,1-1000
penta scan ports --use-nmap -- -sV -sC
```

### Web Enumeration

```bash
penta scan web -u http://10.0.0.1 
penta scan web --dns -u http://example.com
```

### Brute force

```bash
penta brute ssh -u user -P path/to/passwords.txt -h 10.0.0.1
penta brute ftp -u user -P path/to/passwords.txt -h 10.0.0.1
penta brute smb -u user -P path/to/passwords.txt -h 10.0.0.1
penta brute ssh -U path/to/usernames.txt -p password -h 10.0.0.1 --port 22
```

## Layout

```bash
~/.config/penta/        # config files (yaml, json, etc.)
~/.local/share/penta/   # long-term data (DB, history, reports)
~/.local/state/penta/   # mutable state (session store, locks, etc.)
~/.cache/penta/         # caches, can be nuked
```

### Project Structure

```bash
penta/
├── cmd/
│   └── penta/
│       ├── main.go          # cobra root cmd bootstrap
│       ├── root.go          # flags, global config, logging init
│       ├── scan.go          # `penta scan` (non-TUI)
│       ├── scan_tui.go      # `penta scan --tui` or `penta tui`
│       ├── report.go        # `penta report` (NDJSON → md/html)
│       └── version.go       # `penta version`
│
│   internal/
│   ├── engine/
│   │   ├── native.go          # was scan/native/engine.go
│   │   ├── nmap.go            # was scan/nmap/engine.go
│   │   └── options.go         # was scan/options.go (RunOptions, etc.)
│   ├── hosts/
│   │   ├── arp.go             # was scan/native/hosts/arp.go
│   │   ├── discovery.go
│   │   ├── icmp.go
│   │   ├── prober.go
│   │   └── tcp.go
│   ├── ports/
│   │   ├── scanner.go         # was scan/native/ports/scanner.go
    │   ├── tcp.go
    │   ├── udp.go
    │   ├── status.go
    │   ├── version.go
    │   └── service.go         # if you need helpers
│   │
│   ├── targets/
│   │   ├── parser.go        # targets.txt, stdin, CIDR, hostnames
│   │   └── scope.go         # scope.txt, in-scope checks
│   │
│   ├── hosts/
│   │   ├── probe.go         # host prober interface
│   │   ├── tcp_probe.go     # TCP "ping" / connect
│   │   ├── icmp_probe.go    # ICMP (privileged)
│   │   └── arp_probe.go     # ARP scan for LAN
│   │
│   ├── ports/
│   │   ├── scanner.go       # ScanPorts(...) per host
│   │   └── ranges.go        # parse "1-1024,443,8443"
│   │
│   ├── web/
│   │   ├── http_probe.go    # GET /, title, server, status
│   │   ├── tls_info.go      # TLS version, cipher, cert
│   │   └── robots.go        # robots.txt fetch (MVP or v2)
│   │
│   ├── model/
│   │   ├── finding.go       # Finding struct (your NDJSON schema)
│   │   └── event.go         # EventType, Event {Type, Finding, Err, Host}
│   │
│   ├── output/
│   │   ├── output.go        # Sink interface: HandleEvent, Close
│   │   ├── human/
│   │   │   └── human.go     # human-readable stdout (non-TUI)
│   │   ├── ndjson/
│   │   │   └── ndjson.go    # NDJSON sink (machine mode)
│   │   └── tui/
│   │       └── sink.go      # bridges events → Bubble Tea (tea.Msg)
│   │
│   ├── ui/
│   │   ├── model.go         # Bubble Tea model (state, Init, Update, View)
│   │   ├── components.go    # tables, status bars, keymap, styles
│   │   └── messages.go      # Msg types wrapping model.Event, errors, etc.
│   │
│   ├── report/
│   │   ├── reader.go        # read NDJSON → []Finding
│   │   ├── summarize.go     # group by host/check, counts, severities
│   │   └── markdown.go      # md renderer; later html.go for HTML
│   │
│   ├── config/
│   │   ├── config.go        # struct, default values
│   │   └── loader.go        # Viper: flags > env > file
│   │
│   └── log/
│       └── logger.go        # zerolog setup, human vs json logs
│
├── pkg/                     # public helpers (only if you want reuse)
│   └── penta/
│       └── client.go        # (optional) if you later expose engine as lib
│
├── testdata/                # local fixtures, http/tls test servers config
│
├── scripts/                 # dev scripts (lint, fmt, etc.)
│
├── go.mod
└── go.sum

```
