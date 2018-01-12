package smbc

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var dir, file, remote, src, dst, p string
var a, auth bool

var getCMD = &cobra.Command{
	Use:   "get",
	Short: "Downloads file",
	Run: func(cmd *cobra.Command, args []string) {
		if dir == "" || remote == "" || p == "" {
			cmd.Help()
			return
		}

		smbc := New(remote)
		if auth {
			smbc.SetAuthCallback(askAuth)
		}

		f, err := smbc.FileReader(p)
		if err != nil {
			log.Println("failed to get remote file reader: ", err)
			return
		}

		name := filepath.Base(p)
		dst := filepath.Join(dir, name)
		dst, _ = filepath.Abs(dst)

		out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Fatalln(err.Error())
		}

		buf := make([]byte, 64*1024)
		for {
			time.Sleep(time.Second)
			r, err := f.Read(buf)

			if err != nil && err != io.EOF {
				log.Print(err)
				break
			}

			out.Write(buf[:r])
			if err == io.EOF {
				break
			}
		}
		out.Close()
		f.Close()
	},
}
var putCMD = &cobra.Command{
	Use:   "put",
	Short: "Uploads file to the remote samba file tree",
	Run: func(cmd *cobra.Command, args []string) {
		if file == "" || remote == "" || p == "" {
			cmd.Help()
			return
		}

		type UploadFunc func(in io.Reader, out io.Writer) error
		var uf UploadFunc
		uf = func(in io.Reader, out io.Writer) error {
			totalSent := 0
			buf := make([]byte, 64*1024)

			for {
				r, err := in.Read(buf)
				if err != nil && err != io.EOF {
					return err
				}

				if err == io.EOF {
					break
				}

				written, err := out.Write(buf[:r])
				if err != nil {
					log.Println(err)
					return err
				}

				totalSent += written
				line := fmt.Sprintf("\r%d bytes sent", totalSent)
				log.Printf("%30s", line)
			}
			return nil
		}

		name := filepath.Base(file)
		if strings.HasSuffix(p, "/") {
			p += name
		} else {
			p = p + "/" + name
		}

		log.Println("remote:", remote)
		log.Println("file:", file)
		log.Println("path:", p)
		smbc := New(remote)
		if auth {
			smbc.SetAuthCallback(askAuth)
		}

		fw, err := smbc.FileWriter(p, a)
		if err != nil {
			log.Println("failed to get remote file writer: ", err)
			return
		}

		in, err := os.Open(file)
		if err != nil {
			log.Fatalln(err.Error())
		}

		err = uf(in, fw)
		if err != nil {
			if err == syscall.EACCES && !auth {
				//restart with auth
				log.Println("retrying with authentication context...")
				smbc = New(remote)
				smbc.SetAuthCallback(askAuth)
				fw, err := smbc.FileWriter(p, a)
				if err != nil {
					log.Fatalln(err)
				}

				in.Seek(0, io.SeekStart)
				err = uf(in, fw)
				if err != nil {
					log.Fatalln(err)
				}

			} else {
				in.Close()
				fw.Close()
			}
		} else {
			log.Fatalln(err)
		}
	},
}
var statCMD = &cobra.Command{
	Use:   "stat",
	Short: "Get remote file info",
	Run: func(cmd *cobra.Command, args []string) {
		if remote == "" || p == "" {
			cmd.Help()
			return
		}

		smbc := New(remote)
		if auth {
			smbc.SetAuthCallback(askAuth)
		}

		st, err := smbc.Stat(p)
		if err != nil {
			log.Fatalln(err.Error())
		}

		log.Printf("\t%-25s\t:%-25d", "UUID", st.Uid)
		log.Printf("\t%-25s\t:%-25s", "Create date", time.Unix(int64(st.State), 0))
		log.Printf("\t%-25s\t:%-25s", "Edition date", time.Unix(int64(st.Edit), 0))
		log.Printf("\t%-25s\t:%-25s", "Last access date", time.Unix(int64(st.Access), 0))
		log.Printf("\t%-25s\t:%-25d", "Mode", st.Mode)
		log.Printf("\t%-25s\t:%-25d", "Size", st.Size)
		log.Printf("\t%-25s\t:%-25d", "Blocks", st.BlockCount)
		log.Printf("\t%-25s\t:%-25d", "Block size", st.BlockSize)
		log.Printf("\t%-25s\t:%-25d", "Dev", st.Dev)
		log.Printf("\t%-25s\t:%-25d", "Mode", st.Mode)
		log.Printf("\t%-25s\t:%-25d", "Gid", st.Gid)
		log.Printf("\t%-25s\t:%-25d", "INode", st.INode)
		log.Printf("\t%-25s\t:%-25d", "NLink", st.NLink)
		log.Printf("\t%-25s\t:%-25d", "RDev", st.RDev)
	},
}
var mvCMD = &cobra.Command{
	Use:   "mv",
	Short: "Rename file",
	Run: func(cmd *cobra.Command, args []string) {
		if remote == "" || src == "" || dst == "" {
			cmd.Help()
			return
		}

		smbc := New(remote)
		if auth {
			smbc.SetAuthCallback(askAuth)
		}

		err := smbc.Rename(src, dst)
		if err != nil {
			log.Fatalf("failed to rename %s to %s: %s\n", src, dst, err)
		}

		log.Printf("%s successfuly renamed to %s\n", src, dst)
	},
}
var mdCMD = &cobra.Command{
	Use:   "mkdir",
	Short: "Create folder on remote file tree",
	Run: func(cmd *cobra.Command, args []string) {
		if remote == "" || p == "" {
			cmd.Help()
			return
		}

		smbc := New(remote)
		if auth {
			smbc.SetAuthCallback(askAuth)
		}
		err := smbc.MakeDir(p, 777)
		if err != nil {
			log.Fatalln("error: ", err.Error())
		}
		log.Printf("%s directory created\n", remote)
	},
}
var rmCMD = &cobra.Command{
	Use:   "rm",
	Short: "Removed remote file or dir ",
	Run: func(cmd *cobra.Command, args []string) {
		if remote == "" || p == "" {
			cmd.Help()
			return
		}

		smbc := New(remote)
		if auth {
			smbc.SetAuthCallback(askAuth)
		}
		err := smbc.Unlink(p)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%s removed\n", remote)
	},
}
var lsCMD = &cobra.Command{
	Use:   "ls",
	Short: "Get directory children",
	Run: func(cmd *cobra.Command, args []string) {
		if remote == "" || p == "" {
			cmd.Help()
			return
		}

		smbc := New(remote)
		if auth {
			smbc.SetAuthCallback(askAuth)
		}
		ds, err := smbc.DirScan(p, true)
		if err != nil {
			log.Fatalln("openFile dir:", err)
		}

		var getTypeString = func(t os.FileMode) string {
			if t == os.ModeDir {
				return "directory"
			}
			return "file"
		}

		fmt.Printf("\t%-15s\t%-30s\t%-15s\t%s\n",
			"Type",
			"last access",
			"size",
			"Name",
		)

		entry, err := ds.Next()
		for entry != nil && err == nil {
			if entry.Info != nil {
				fmt.Printf("\t%-15s\t%-30s\t%-15d\t%s\n",
					getTypeString(entry.Mode),
					time.Unix(int64(entry.Info.Access), 0),
					entry.Info.Size,
					entry.Name,
				)
			} else {
				fmt.Printf("\t%-15s\t%-30s\t%-15d\t%s\n",
					getTypeString(entry.Mode),
					"",
					0,
					entry.Name,
				)
			}
			entry, err = ds.Next()
		}
		if err != nil && err.Error() != "EOF" {
			log.Println(err)
		}
		ds.Close()
	},
}
var testCMD = &cobra.Command{
	Use:   "test",
	Short: "Runs multiple routine doing the same things",
	Run: func(cmd *cobra.Command, args []string) {

		smbc := New("smb://192.168.0.5")
		var ls = func(p string) {
			//defer smbc.Close()
			ds, err := smbc.DirScan(p, true)
			if err != nil {
				log.Fatalln("openFile dir:", err)
			}
			var getTypeString = func(t os.FileMode) string {
				if t == os.ModeDir {
					return "directory"
				}
				return "file"
			}
			fmt.Printf("\t%-15s\t%-30s\t%-15s\t%s\n",
				"Type",
				"last access",
				"size",
				"Name",
			)
			entry, err := ds.Next()
			for entry != nil && err == nil {
				if entry.Info != nil {
					fmt.Printf("\t%-15s\t%-30s\t%-15d\t%s\n",
						getTypeString(entry.Mode),
						time.Unix(int64(entry.Info.Access), 0),
						entry.Info.Size,
						entry.Name,
					)
				} else {
					fmt.Printf("\t%-15s\t%-30s\t%-15d\t%s\n",
						getTypeString(entry.Mode),
						"",
						0,
						entry.Name,
					)
				}
				entry, err = ds.Next()
			}
			if err != nil && err.Error() != "EOF" {
				log.Println(err)
			}
			//ds.Close()
		}

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			ls("/Downloads/Jabar")
		}()
		go func() {
			defer wg.Done()
			ls("/Downloads/Jabar")
		}()
		wg.Wait()
		log.Println("done!")
		smbc.Close()

	},
}
var RootCMD = &cobra.Command{
	Use:   "smbc",
	Short: "SMB smbc",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	getCMD.PersistentFlags().StringVar(&remote, "remote", "", "Server host")
	getCMD.PersistentFlags().StringVar(&p, "path", "", "Remote file path")
	getCMD.PersistentFlags().StringVar(&dir, "dir", "", "Local directory path")
	getCMD.PersistentFlags().BoolVar(&auth, "auth", false, "With authentication")

	putCMD.PersistentFlags().StringVar(&remote, "remote", "", "Server host")
	putCMD.PersistentFlags().StringVar(&p, "path", "", "Remote file path")
	putCMD.PersistentFlags().StringVar(&file, "file", "", "Local file to upload")
	putCMD.PersistentFlags().BoolVar(&a, "append", false, "Append content")
	putCMD.PersistentFlags().BoolVar(&auth, "auth", false, "With authentication")

	mvCMD.PersistentFlags().StringVar(&remote, "remote", "", "Server host")
	mvCMD.PersistentFlags().StringVar(&p, "path", "", "Remote file path")
	mvCMD.PersistentFlags().StringVar(&src, "src", "", "old path")
	mvCMD.PersistentFlags().StringVar(&dst, "dst", "", "new path")
	mvCMD.PersistentFlags().BoolVar(&auth, "auth", false, "With authentication")

	statCMD.PersistentFlags().StringVar(&remote, "remote", "", "Server host")
	statCMD.PersistentFlags().StringVar(&p, "path", "", "Remote file path")
	statCMD.PersistentFlags().BoolVar(&auth, "auth", false, "With authentication")

	mdCMD.PersistentFlags().StringVar(&remote, "remote", "", "Server host")
	mdCMD.PersistentFlags().StringVar(&p, "path", "", "Remote file path")
	mdCMD.PersistentFlags().BoolVar(&auth, "auth", false, "With authentication")

	rmCMD.PersistentFlags().StringVar(&remote, "remote", "", "Server host")
	rmCMD.PersistentFlags().StringVar(&p, "path", "", "Remote file path")
	rmCMD.PersistentFlags().BoolVar(&auth, "auth", false, "With authentication")

	lsCMD.PersistentFlags().StringVar(&remote, "remote", "", "Server host")
	lsCMD.PersistentFlags().StringVar(&p, "path", "", "Remote file path")
	lsCMD.PersistentFlags().BoolVar(&auth, "auth", false, "With authentication")

	RootCMD.AddCommand(getCMD)
	RootCMD.AddCommand(putCMD)
	RootCMD.AddCommand(statCMD)
	RootCMD.AddCommand(mdCMD)
	RootCMD.AddCommand(rmCMD)
	RootCMD.AddCommand(lsCMD)
	RootCMD.AddCommand(mvCMD)
	RootCMD.AddCommand(testCMD)

	RootCMD.AddCommand(putCMD)
}

//go build -pkgdir /home/pydio/go/pkg/ -o gosmb main.go
