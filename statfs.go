package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// Display the file system type of a path given as argument
// usage: go run statfs.go /path/to/target
// examples:
//	$ go run statfs.go /tmp
//	/tmp has code 0x1021994 : TMPFS_MAGIC
//
//	$ go run statfs.go /boot
//	/boot has code 0xEF53 EXT4_SUPER_MAGIC
//
//	$ go run statfs.go /home
//	/home has code 0x9123683E unknown   <= should be BTRFS
//
//	$ go run statfs.go /some/remote/filesystem
//	/some/remote/filesystem has code 0xFF534D42 CIFS_MAGIC_NUMBER
func main() {
	fstypes := map[int64]string{
		0xadf5:     "ADFS_SUPER_MAGIC",
		0xADFF:     "AFFS_SUPER_MAGIC",
		0x42465331: "BEFS_SUPER_MAGIC",
		0x1BADFACE: "BFS_MAGIC",
		0xFF534D42: "CIFS_MAGIC_NUMBER",
		0x73757245: "CODA_SUPER_MAGIC",
		0x012FF7B7: "COH_SUPER_MAGIC",
		0x28cd3d45: "CRAMFS_MAGIC",
		0x1373:     "DEVFS_SUPER_MAGIC",
		0x00414A53: "EFS_SUPER_MAGIC",
		0x137D:     "EXT_SUPER_MAGIC",
		0xEF51:     "EXT2_OLD_SUPER_MAGIC",
		0xEF53:     "EXT4_SUPER_MAGIC",
		0x4244:     "HFS_SUPER_MAGIC",
		0xF995E849: "HPFS_SUPER_MAGIC",
		0x958458f6: "HUGETLBFS_MAGIC",
		0x9660:     "ISOFS_SUPER_MAGIC",
		0x72b6:     "JFFS2_SUPER_MAGIC",
		0x3153464a: "JFS_SUPER_MAGIC",
		0x137F:     "MINIX_SUPER_MAGIC",
		0x138F:     "MINIX_SUPER_MAGIC2",
		0x2468:     "MINIX2_SUPER_MAGIC",
		0x2478:     "MINIX2_SUPER_MAGIC2",
		0x4d44:     "MSDOS_SUPER_MAGIC",
		0x564c:     "NCP_SUPER_MAGIC",
		0x6969:     "NFS_SUPER_MAGIC",
		0x5346544e: "NTFS_SB_MAGIC",
		0x9fa1:     "OPENPROM_SUPER_MAGIC",
		0x9fa0:     "PROC_SUPER_MAGIC",
		0x002f:     "QNX4_SUPER_MAGIC",
		0x52654973: "REISERFS_SUPER_MAGIC",
		0x7275:     "ROMFS_MAGIC",
		0x517B:     "SMB_SUPER_MAGIC",
		0x012FF7B6: "SYSV2_SUPER_MAGIC",
		0x012FF7B5: "SYSV4_SUPER_MAGIC",
		0x01021994: "TMPFS_MAGIC",
		0x15013346: "UDF_SUPER_MAGIC",
		0x00011954: "UFS_MAGIC",
		0x9fa2:     "USBDEVICE_SUPER_MAGIC",
		0xa501FCF5: "VXFS_SUPER_MAGIC",
		0x012FF7B4: "XENIX_SUPER_MAGIC",
		0x58465342: "XFS_SUPER_MAGIC",
		0x012FD16D: "_XIAFS_SUPER_MAGIC",
	}
	var s syscall.Statfs_t
	err := syscall.Statfs(os.Args[1], &s)
	if err != nil {
		panic(err)
	}
	code := fmt.Sprintf("0x%s", strings.ToUpper(strconv.FormatInt(s.Type, 16)))
	fstype := "unknown"
	if _, ok := fstypes[s.Type]; ok {
		fstype = fstypes[s.Type]
	}
	fmt.Println(os.Args[1], "has code", code, fstype)
}
