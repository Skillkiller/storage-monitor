# storage-monitor

A lightweight Go-based tool for monitoring disk usage on remote targets via rclone.  
Provides on-demand metrics in Prometheus format.

This project was inspired by [zekroTJA/ssh-storage-monitor](https://github.com/zekroTJA/ssh-storage-monitor).

## Features

- Connects to remote targets via rclone
- Monitors remote storage quotas via `rclone about --json` (single API call per remote)
- Monitors storage usage per folder via `rclone lsjson --recursive`
- Exports metrics in Prometheus format
- On-demand statistics collection: metrics are gathered only when scraped

## Note on rclone usage

### Remote quotas (`about_targets`)

For remote-level storage metrics (`total`, `used`, `free`, `trashed`, `other` in bytes), configure `about_targets` in `config.yml`. This uses `rclone about <remote> --json`, which only requires a single API call per remote (assuming the backend supports `about`).

Exported metrics:
- `rclone_about_total_bytes`
- `rclone_about_used_bytes`
- `rclone_about_free_bytes`
- `rclone_about_trashed_bytes`
- `rclone_about_other_bytes`
- `rclone_about_up`

### Folder details (`watched_dirs` / `recursive_watched_dirs`)

This tool uses the command `rclone lsjson --recursive` to gather detailed file information.

> ⚠️ **Warning:** This operation may trigger a large number of API calls, especially for remotes with many files or deeply nested directories.

However, it is required to accurately determine:
- the total number of files
- the time of the last modification
- the total size of all files
