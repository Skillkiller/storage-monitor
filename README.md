# ssh-storage-monitor

A lightweight Go-based tool for monitoring disk usage on remote Linux servers via SSH.

## Features

- Secure connection to remote hosts via SSH
- Uses only a minimal set of commands (`ls`, `df`, `du`)
- Monitor storage usage per folder
- Monitor storage usage per filesystem
- Export metrics in Prometheus format
- On-demand statistics collection: metrics are gathered only when scraped
