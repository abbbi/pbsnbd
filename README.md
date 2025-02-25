# Proxmox Backup Server NBDKIT Plugin

Export your virtual machine disk backups from a Proxmox Backup Server via NBD.

To start an export, pass the plugin to nbdkit with required arguments:

```
 nbdkit --filter cow -f ./nbdkit-pbs-plugin.so \
    image=drive-scsi0.img \
    vmid=103 \ 
    timestamp="2025-02-25T10:32:54Z" \
    repo="root@pam@192.168.161.241:test" \
    fingerprint="EE:2F:53:29:F2:69:61:0C:D5:96:55:50:6B:2B:84:8E:EF:3E:A6:B0:CA:18:5C:F7:92:BA:54:71:15:56:83:5B" \
    password=test \
    namespace= \

 Connected to PBS version: [1.4.1 (UNKNOWN)]
 Attempt to open image [vm/103/2025-02-25T10:32:54Z/drive-scsi0.img]
 Successfully opened image [vm/103/2025-02-25T10:32:54Z/drive-scsi0.img]
```

The NBD device is by default reachable via localhost. The COW Filter allows
read/write operations to the disks without altering the original data. This
way, you can boot off the NBD device directly using Qemu:

```
 qemu-system-x86_64 -m 2000 -hda nbd://localhost -cpu host -enable-kvm
```

Or, as alternative, map the nbd device locally and access file systems:

```
 qemu-nbd  -c /dev/nbd5 nbd://localhost
 fdisk -l /dev/nbd5 
 Disk /dev/nbd5: 128 GiB, 137438953472 bytes, 268435456 sectors
 [..]
 Device      Boot   Start       End   Sectors   Size Id Type
 /dev/nbd5p1 *       2048    999423    997376   487M 83 Linux
 /dev/nbd5p2       999424   4999167   3999744   1.9G 82 Linux swap / Solaris
 /dev/nbd5p3      4999168 268433407 263434240 125.6G 83 Linux
```

Access via regular nbd tools:

```
 nbdinfo nbd://localhost
 protocol: newstyle-fixed without TLS, using structured packets
 export="":
        export-size: 137438953472 (128G)
        content: DOS/MBR boot sector
        uri: nbd://localhost:10809/
 [..]
```

# building

```
 sudo apt-get install nbdkit-plugin-dev
 make
```
