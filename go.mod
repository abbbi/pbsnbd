module pbsnbd

go 1.24.0

require github.com/elbandi/go-proxmox-backup-client v0.0.0-20240901114840-dcb8a5fcb96c // indirect

replace libguestfs.org/nbdkit => ./nbdkit

require libguestfs.org/nbdkit v1.0.0
