package smbc

import (
	"errors"
	"io"
	"os"
	"sync"
	"syscall"
	"unsafe"
	//"bufio"
	//"strings"
	"fmt"
	"log"
	"strings"
)

/*
#cgo LDFLAGS: -lsmbclient
#cgo CFLAGS: -I/usr/include/samba-4.0
#include "smbc.h"
*/
import "C" // DO NOT CHANGE THE POSITION OF THIS IMPORT

// SmbType
type SmbcType int

const (
	SMBC_WORKGROUP     = C.SMBC_WORKGROUP
	SMBC_FILE_SHARE    = C.SMBC_FILE_SHARE
	SMBC_PRINTER_SHARE = C.SMBC_PRINTER_SHARE
	SMBC_COMMS_SHARE   = C.SMBC_COMMS_SHARE
	SMBC_IPC_SHARE     = C.SMBC_IPC_SHARE
	SMBC_DIR           = C.SMBC_DIR
	SMBC_FILE          = C.SMBC_FILE
	SMBC_LINK          = C.SMBC_LINK
)

type Dirent struct {
	Type    SmbcType
	Comment string
	Name    string
}

type NotifyAction struct {
	action   uint32
	filename *C.char
}

type inoT int
type modeT int
type nlinkT int
type uidT int
type gidT int
type devT int32
type offT int32
type blksizeT int32
type blkcntT uint32
type timeT int64

type Stat struct {
	/*struct stat {
		dev_t     st_dev;      ID du périphérique contenant le fichier
		ino_t     st_ino;      Numéro inœud
		mode_t    st_mode;     Protection
		nlink_t   st_nlink;    Nb liens matériels
		uid_t     st_uid;      UID propriétaire
		gid_t     st_gid;      GID propriétaire
		dev_t     st_rdev;     ID périphérique (si fichier spécial)
		off_t     st_size;     Taille totale en octets
		blksize_t st_blksize;  Taille de bloc pour E/S
		blkcnt_t  st_blocks;   Nombre de blocs alloués
		time_t    st_atime;    Heure dernier accès
		time_t    st_mtime;    Heure dernière modification
		time_t    st_ctime;    Heure dernier changement état
	};*/
	Dev        devT
	INode      inoT
	Mode       modeT
	NLink      nlinkT
	Uid        uidT
	Gid        gidT
	RDev       devT
	Size       offT
	BlockSize  blksizeT
	BlockCount blkcntT
	Access     timeT
	Edit       timeT
	State      timeT
}

// *sigh* even with libsmbclient-4.0 the library is not MT safe,
// e.g. smbc_init_context from multiple threads crashes

var glock = sync.Mutex{}

func lock() {
	glock.Lock()
}

func unlock() {
	glock.Unlock()
}

func askAuth(server_name, share_name string) (string, string, string) {
	workgroup := ""
	username := "tran"
	password := "Mcdanol"

	fmt.Printf("auth for %s %s\n", server_name, share_name)
	fmt.Println("Workgroup:", workgroup)
	fmt.Println("Username:", username)
	fmt.Println("Password:", password)
	return workgroup, username, password
}

func setEcho(terminal_echo_enabled bool) {
	/*var cmd *exec.Cmd
	if terminal_echo_enabled {
		cmd = exec.Command("stty", "-F", "/dev/tty", "echo")
	} else {
		cmd = exec.Command("stty", "-F", "/dev/tty", "-echo")
	}
	cmd.Run()*/
}

//************************************************************************** CONTEXT
type Smbc struct {
	//mutex        sync.Mutex // libsmbclient is not thread safe
	url          string
	ctx          *C.SMBCCTX
	authCallback *func(string, string) (string, string, string)
}

func New(url string) *Smbc {
	s := new(Smbc)
	s.ctx = C.smbc_new_context()
	C.smbc_init_context(s.ctx)
	s.url = url
	return s
}

func (c *Smbc) pathUrl(p string) string {

	if !strings.HasPrefix(c.url, "smb://") {
		c.url = "smb://" + c.url
	}
	return c.url + p
}

func (c *Smbc) SetAuthCallback(fn func(string, string) (string, string, string)) {
	lock()
	defer unlock()
	C.my_smbc_init_auth_callback(c.ctx, unsafe.Pointer(&fn))
	// we need to store it in the Smbc struct to ensure its not garbage
	// collected later (I think)
	c.authCallback = &fn
}

func (c *Smbc) Destroy() error {
	return c.Close()
}

func (c *Smbc) Close() error {
	// FIXME: is there a more elegant way for this c.lock.Lock() that
	//        needs to be part of every function? python decorator to
	//        the rescue :)
	lock()
	defer unlock()
	log.Println("Freeing context...")
	var err error
	if c.ctx != nil {
		// 1 would mean we force the destroy
		res := C.smbc_free_context(c.ctx, C.int(1))
		if res < 0 {
			e := C.my_errno()
			err = syscall.Errno(e)
		}
		c.ctx = nil
	}
	return err
}

func (c *Smbc) GetDebug() int {
	lock()
	defer unlock()
	return int(C.smbc_getDebug(c.ctx))
}

func (c *Smbc) SetDebug(level int) {
	lock()
	defer unlock()

	C.smbc_setDebug(c.ctx, C.int(level))
}

func (c *Smbc) GetUser() string {
	lock()
	defer unlock()

	return C.GoString(C.smbc_getUser(c.ctx))
}

func (c *Smbc) SetUser(user string) {
	lock()
	defer unlock()

	C.smbc_setUser(c.ctx, C.CString(user))
}

func (c *Smbc) GetWorkGroup() string {
	lock()
	defer unlock()

	return C.GoString(C.smbc_getWorkgroup(c.ctx))
}

func (c *Smbc) SetWorkGroup(wg string) {
	lock()
	defer unlock()

	C.smbc_setWorkgroup(c.ctx, C.CString(wg))
}

func (c *Smbc) openFile(p string, flags int, mode modeT) (*C.SMBCFILE, error) {
	url := c.pathUrl(p)
	sf := C.my_smbc_open(c.ctx, C.CString(url), C.int(flags), C.mode_t(mode))
	if sf == nil {
		return nil, syscall.Errno(C.my_errno())
	}
	return sf, nil
}

func (c *Smbc) MakeDir(p string, mode modeT) error {
	lock()
	defer unlock()
	url := c.pathUrl(p)
	res := C.my_smbc_mkdir(c.ctx, C.CString(url), C.mode_t(mode))
	if res < 0 {
		return syscall.Errno(C.my_errno())
	}
	return nil
}

func (c *Smbc) Stat(p string) (*Stat, error) {
	var stat C.c_stat
	url := c.pathUrl(p)
	r := C.my_smbc_stat(c.ctx, C.CString(url), &stat)
	if r == 0 {
		return &Stat{
			Dev:        devT(stat.st_dev),
			INode:      inoT(stat.st_ino),
			Mode:       modeT(stat.st_mode),
			NLink:      nlinkT(stat.st_nlink),
			Uid:        uidT(stat.st_uid),
			Gid:        gidT(stat.st_gid),
			RDev:       devT(stat.st_rdev),
			Size:       offT(stat.st_size),
			BlockSize:  blksizeT(stat.st_blksize),
			BlockCount: blkcntT(stat.st_blocks),
			Access:     timeT(stat.st_atim.tv_sec),
			Edit:       timeT(stat.st_mtim.tv_sec),
			State:      timeT(stat.st_ctim.tv_sec),
		}, nil
	}
	return nil, syscall.Errno(C.my_errno())
}

func (c *Smbc) Create(p string, mode modeT) error {
	lock()
	defer unlock()
	url := c.pathUrl(p)
	fd := C.my_smbc_creat(c.ctx, C.CString(url), C.mode_t(mode))
	if fd == nil {
		return syscall.Errno(C.my_errno())
	}
	return nil
}

func (c *Smbc) FileReader(p string) (*SmbcFileContent, error) {
	lock()
	defer unlock()

	f, err := c.openFile(p, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	fc := &SmbcFileContent{
		client:  c,
		offset:  -1,
		mode:    os.O_RDONLY,
		path:    p,
		fd:      f,
		authErr: false,
	}
	return fc, nil
}

func (c *Smbc) FileWriter(p string, append bool) (*SmbcFileContent, error) {

	mode := os.O_WRONLY | os.O_CREATE
	if append {
		mode |= os.O_APPEND
	} else {
		mode |= os.O_TRUNC
	}

	f, err := c.openFile(p, mode, 0666)
	if err != nil {
		return nil, err
	}

	fc := &SmbcFileContent{
		client:  c,
		offset:  -1,
		mode:    mode,
		path:    p,
		fd:      f,
		authErr: false,
	}
	return fc, nil
}

func (c *Smbc) DirScan(p string, stat bool) (*SmbcDirScanner, error) {
	lock()
	defer unlock()
	url := c.pathUrl(p)
	d := C.my_smbc_opendir(c.ctx, C.CString(url))
	if d == nil {
		return nil, syscall.Errno(C.my_errno())
	}

	return &SmbcDirScanner{
		client:   c,
		path:     p,
		fd:       d,
		withStat: stat,
	}, nil
}

func (c *Smbc) Unlink(p string) error {
	lock()
	defer unlock()
	url := c.pathUrl(p)
	res := C.my_smbc_unlink(c.ctx, C.CString(url))
	if res < 0 {
		return syscall.Errno(C.my_errno())
	}
	return nil
}

func (c *Smbc) Rename(oldUrl string, p string) error {
	lock()
	defer unlock()
	url := c.pathUrl(p)
	//nCtx := C.smbc_new_context()
	//C.smbc_init_context(nCtx)
	//res := C.my_smbc_rename(c.ctx, C.CString(oldUrl), nCtx, C.CString(nUrl))
	res := C.my_smbc_rename(c.ctx, C.CString(oldUrl), c.ctx, C.CString(url))
	if res < 0 {
		return syscall.Errno(C.my_errno())
	}
	//c.ctx = nCtx
	return nil
}

func (c *Smbc) WatchChanges(url string, recursive bool, handler EventHandler) error {
	return nil
}

//************************************************************************** DIRECTORY
type SmbcDirScanner struct {
	fd       *C.SMBCFILE
	client   *Smbc
	path     string
	withStat bool
}

type DirEntry struct {
	Mode    os.FileMode
	Name    string
	Comment string
	Info    *Stat
}

func (ds *SmbcDirScanner) Next() (*DirEntry, error) {
	lock()
	cDirEnt := C.my_smbc_readdir(ds.client.ctx, ds.fd)
	if cDirEnt == nil {
		unlock()
		return nil, io.EOF
	}
	unlock()

	typ := SmbcType(cDirEnt.smbc_type)
	de := &DirEntry{
		Name:    C.GoStringN(&cDirEnt.name[0], C.int(cDirEnt.namelen)),
		Comment: C.GoStringN(cDirEnt.comment, C.int(cDirEnt.commentlen)),
	}

	switch typ {
	case SMBC_DIR:
		de.Mode = os.ModeDir
	case SMBC_FILE_SHARE:
		de.Mode = os.ModeDir
	case SMBC_FILE:
		de.Mode = os.ModeType
	case SMBC_LINK:
		de.Mode = os.ModeSymlink
	}

	if ds.withStat {
		childPath := ds.path + "/" + de.Name
		de.Info, _ = ds.client.Stat(childPath)
	}
	return de, nil
}

func (ds *SmbcDirScanner) Close() error {
	lock()
	defer unlock()
	res := C.my_smbc_closedir(ds.client.ctx, ds.fd)
	if res < 0 {
		return syscall.Errno(C.my_errno())
	}
	return nil
}

//************************************************************************** STREAM
type SmbcFileContent struct {
	path      string
	client    *Smbc
	fd        *C.SMBCFILE
	offset    int64
	totalSize int64
	closed    bool
	mode      int
	authErr   bool
}

func (fc *SmbcFileContent) Mode() int {
	return fc.mode
}

func (fc *SmbcFileContent) Read(b []byte) (int, error) {
	if fc.mode&os.O_RDONLY != os.O_RDONLY {
		return -1, errors.New("not in read mode")
	}
	if fc.closed {
		return -1, io.ErrClosedPipe
	}

	if fc.offset > 0 && fc.offset == fc.totalSize {
		return 0, io.EOF
	}

	var e error

	lock()
	defer unlock()
	cCount := C.my_smbc_read(fc.client.ctx, fc.fd, unsafe.Pointer(&b[0]), C.size_t(len(b)))
	if cCount == 0 {
		return 0, io.EOF
	}

	if cCount == -1 {
		e = syscall.Errno(C.my_errno())
	}

	fc.offset += int64(cCount)
	return int(cCount), e
}

func (fc *SmbcFileContent) Write(b []byte) (int, error) {
	if fc.mode&os.O_WRONLY != os.O_WRONLY {
		return -1, errors.New("not in write mode")
	}

	if fc.closed {
		return -1, io.ErrClosedPipe
	}

	if fc.offset > 0 && fc.offset == fc.totalSize {
		return 0, io.EOF
	}

	var e error

	lock()
	defer unlock()

	cCount := C.my_smbc_write(fc.client.ctx, fc.fd, unsafe.Pointer(&b[0]), C.size_t(len(b)))
	if cCount == 0 {
		return 0, io.EOF
	}

	if cCount == -1 {
		e = syscall.Errno(C.my_errno())
	}

	fc.offset += int64(cCount)
	return int(cCount), e
}

func (fc *SmbcFileContent) Seek(offset int64, whence int) (int64, error) {
	lock()
	defer unlock()
	newOffset := C.my_smbc_lseek(fc.client.ctx, fc.fd, C.off_t(offset), C.int(whence))
	if newOffset < 0 {
		return fc.offset, syscall.Errno(newOffset)
	}
	fc.offset = int64(newOffset)
	return int64(newOffset), nil
}

func (fc *SmbcFileContent) Close() error {
	lock()
	defer unlock()

	fc.closed = true
	d := C.my_smbc_close(fc.client.ctx, fc.fd)
	if d < 0 {
		return syscall.Errno(C.my_errno())
	}
	return nil
}

//************************************************************************** EVENT
type SmbcChange struct {
	Action int
	Path   string
}

type EventHandler interface {
	Handle(*SmbcChange)
}

//************************************************************************** EXPORTED FUNCTIONS

//export GoAuthCallbackHelper
func GoAuthCallbackHelper(
	o unsafe.Pointer,
	server_name,
	share_name *C.char,
	domain_out *C.char,
	domain_len C.int,
	username_out *C.char,
	username_len C.int,
	password_out *C.char,
	password_len C.int) {

	goFn := *(*func(server_name, share_name string) (string, string, string))(o)
	domain, user, pw := goFn(C.GoString(server_name), C.GoString(share_name))
	C.strncpy(domain_out, C.CString(domain), C.size_t(domain_len))
	C.strncpy(username_out, C.CString(user), C.size_t(username_len))
	C.strncpy(password_out, C.CString(pw), C.size_t(password_len))
}

/*//export GoChangeNotifyCallback
func GoChangeNotifyCallback(action []NotifyAction, privateData unsafe.Pointer) {

}*/
