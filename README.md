# ssh-storage-monitor

A lightweight Go-based tool for monitoring disk usage on remote targets via rclone.

## Features

- connection to remote remote targets via rclone
- Monitor storage usage per folder
- Export metrics in Prometheus format
- On-demand statistics collection: metrics are gathered only when scraped

## Note on rclone usage

This tool uses the command `rclone lsjson --recursive` to gather detailed file information.

Please note that this operation may trigger a large number of API calls, especially for remotes with many files or deeply nested directories.  
However, it is required to accurately determine:
- the total number of files,
- the time of the last modification, and
- the total size of all files.
