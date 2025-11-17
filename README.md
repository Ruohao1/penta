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
