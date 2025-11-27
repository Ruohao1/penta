# Scan

## Hosts

It search for up hosts in a list of ips, it can take a list of ip ranges or list, of ip addresses or cidr ranges.

### Commands

```bash
penta scan hosts 10.0.0.0/24
```

### Methods

#### Unprivileged

By default, it will perform a tcp scan on default ports (22, 80, 443).
Possible outcomes:

- A successful scan means the host is up
- `ECONNREFUSED`, `ECONNRESET` means the host is up

But if timed out, it might be filtered
You can still pass the -m arp flag to look at the neighbors table to see if the host was cached as up.
As unprivileged, it cannot perform neither icmp ping nor arp ping.

#### Privileged
