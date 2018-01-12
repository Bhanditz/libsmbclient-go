// +build ignore

package main

import (
	smbc "github.com/pydio/libsmbclient-go"
)

//func getStat(client *libsmbclient.Smbc, uri string) {
//	st, err := client.Stat(uri)
//	if err != nil {
//		log.Println("Response:", err.Error())
//	} else {
//		log.Println("Stat response:", st)
//	}
//}
//
//func smbMkdir(client *libsmbclient.Smbc, duri string) {
//	err := client.MakeDir(duri, 777)
//	if err != nil {
//		log.Println("error: ", err.Error())
//	} else {
//		log.Println("dir created")
//	}
//}
//
//func openSmbdir(client *libsmbclient.Smbc, duri string) {
//	dh, _ := client.OpenDir(duri)
//	/*if err != nil {
//		log.Print("open dir failed:", err)
//		return
//	}*/
//	var getTypeString= func(t libsmbclient.SmbcType) string{
//
//		/*SMBC_WORKGROUP     SmbcType = C.SMBC_WORKGROUP
//		SMBC_FILE_SHARE             = C.SMBC_FILE_SHARE
//		SMBC_PRINTER_SHARE          = C.SMBC_PRINTER_SHARE
//		SMBC_COMMS_SHARE            = C.SMBC_COMMS_SHARE
//		SMBC_IPC_SHARE              = C.SMBC_IPC_SHARE
//		SMBC_DIR                    = C.SMBC_DIR
//		SMBC_FILE                   = C.SMBC_FILE
//		SMBC_LINK                   = C.SMBC_LINK*/
//
//		switch t{
//		case libsmbclient.SMBC_FILE:
//			return "file"
//		case libsmbclient.SMBC_DIR:
//			return "directory"
//		case libsmbclient.SMBC_LINK:
//			return "link"
//		default:
//			return "unknown"
//		}
//	}
//
//	log.Printf("%15s\t%15s\t%15s\t%15s\t%s",
//
//		"Type",
//		"a time",
//		"e time",
//		"size",
//		"Name",
//	)
//	for {
//		dirent, err := dh.ReadDir()
//		if err != nil {
//			break
//		}
//		furi := duri + "/" + dirent.Name
//		st, _ := client.Stat(furi)
//		if st != nil {
//			log.Printf("%15s\t%15d\t%15d\t%15d\t%s",
//				getTypeString(dirent.Type),
//				st.Access,
//				st.Edit,
//				st.Size,
//				dirent.Name,
//			)
//		}
//	}
//	dh.CloseDir()
//}
//
//func openSmbfile(client *libsmbclient.Smbc, furi string) {
//	f, err := client.FileReader(furi)
//	if err != nil {
//		log.Println("failed to get file reader: ", err)
//		return
//	}
//
//	//buf := make([]byte, 64*1024)
//	buf := make([]byte, 2)
//	for {
//		time.Sleep(time.Second)
//		r, err := f.Read(buf)
//		if r > 0 {
//			fmt.Print(string(buf[:r]))
//		}
//
//		if err == io.EOF {
//			break
//		}
//
//		if err != nil {
//			log.Print(err)
//			f.Close()
//			break
//		}
//	}
//	log.Println()
//	f.Close()
//}
//
//func askAuth(server_name, share_name string)(out_domain, out_username, out_password string) {
//	bio := bufio.NewReader(os.Stdin)
//	fmt.Printf("auth for %s %s\n", server_name, share_name)
//	// domain
//	fmt.Print("Domain: ")
//	domain, _, _ := bio.ReadLine()
//	// read username
//	fmt.Print("Username: ")
//	username, _, _ := bio.ReadLine()
//	// read pw from stdin
//	fmt.Print("Password: ")
//	setEcho(false)
//	password, _, _ := bio.ReadLine()
//	setEcho(true)
//	return strings.TrimSpace(string(domain)), strings.TrimSpace(string(username)), strings.TrimSpace(string(password))
//}
//
//func setEcho(terminal_echo_enabled bool) {
//	var cmd *exec.Cmd
//	if terminal_echo_enabled {
//		cmd = exec.Command("stty",  "-F", "/dev/tty", "echo")
//	} else  {
//		cmd = exec.Command("stty",  "-F", "/dev/tty", "-echo")
//	}
//	cmd.Run()
//}
//
//func multiThreadStressTest(client *libsmbclient.Smbc, uri string) {
//	fmt.Println("m: "+uri)
//	dh, err := client.OpenDir(uri)
//	if err != nil {
//		log.Print(err)
//		return
//	}
//	for {
//		dirent, err := dh.ReadDir()
//		if err != nil {
//			break
//		}
//		newUri := uri + "/" + dirent.Name
//		switch (dirent.Type) {
//		case libsmbclient.SMBC_DIR, libsmbclient.SMBC_FILE_SHARE:
//			fmt.Println("d: "+newUri)
//			go multiThreadStressTest(client, newUri)
//		case libsmbclient.SMBC_FILE:
//			fmt.Println("f: "+newUri)
//			go openSmbfile(client, newUri)
//		}
//	}
//	dh.CloseDir()
//
//	// FIXME: instead of sleep, wait for all threads to exit
//	time.Sleep(10*time.Second)
//}


func main() {
	smbc.RootCMD.Execute()
}
