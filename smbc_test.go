package smbc

import (
	"path"
	"testing"
	_ "time"

	"github.com/smartystreets/goconvey/convey"
)

func TestSmbc(t *testing.T) {
	convey.Convey("Create directory", t, func() {
		var host = "192.168.0.5"
		var share = "/Downlaods"
		var testDir = "testing"

		client := New(host)
		s, err := client.Stat(share)
		convey.So(err, convey.ShouldBeNil)
		convey.So(s, convey.ShouldNotBeNil)

		//mkdir create tests dir
		err = client.MakeDir(path.Join(share, testDir), 0777)
		convey.So(err, convey.ShouldBeNil)

		// upload content
		var filemane = "upload.txt"
		w, err := client.FileWriter(path.Join(share, testDir, filemane), false)
		convey.So(err, convey.ShouldBeNil)

		bytes := []byte("uploaded text content")
		written, err := w.Write(bytes)
		convey.So(err, convey.ShouldBeNil)
		convey.So(written, convey.ShouldEqual, len(bytes))
		w.Close()

		s, err = client.Stat(path.Join(share, testDir, filemane))
		convey.So(err, convey.ShouldBeNil)
		convey.So(s, convey.ShouldNotBeNil)
		convey.So(int(s.Size), convey.ShouldEqual, len(bytes))

		var oldPath = path.Join(share, testDir, filemane)
		var newPath = path.Join(share, testDir, "renamed.txt")

		err = client.Rename(oldPath, newPath)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestStat(t *testing.T) {

}
