# storage-monitor

A lightweight Go-based tool for monitoring disk usage on remote targets via rclone.  
Provides on-demand metrics in Prometheus format.

This project was inspired by [zekroTJA/ssh-storage-monitor](https://github.com/zekroTJA/ssh-storage-monitor).

## Features

- Connects to remote targets via rclone
- Monitors storage usage per folder
- Exports metrics in Prometheus format
- On-demand statistics collection: metrics are gathered only when scraped

## Note on rclone usage

This tool uses the command `rclone lsjson --recursive` to gather detailed file information.

> ⚠️ **Warning:** This operation may trigger a large number of API calls, especially for remotes with many files or deeply nested directories.

However, it is required to accurately determine:
- the total number of files
- the time of the last modification
- the total size of all files
