package compress

import (
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestCompressPath(t *testing.T) {
	tarDir := t.TempDir()
	tempDir := t.TempDir()
	os.WriteFile(tempDir+"/test1.txt", []byte("test"), 0644)
	testBuf, err := os.OpenFile(tarDir+"/test1.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	err = CompressPath(tempDir+"/test1.txt", testBuf)
	if err != nil {
		t.Fatal(err)
	}
	testBuf.Close()
	cmd := exec.Command("tar", "-xvf", tarDir+"/test1.tar.gz", "-C", tarDir)
	cmd.Run()
	if _, err := os.Stat(tarDir + "/test1.txt"); os.IsNotExist(err) {
		t.Fatal(err)
	}
	os.Remove(tarDir + "/test1.txt")
	os.WriteFile(tempDir+"/test2.txt", []byte("test2"), 0644)
	testBuf, err = os.OpenFile(tarDir+"/test2.tar.gz", os.O_CREATE|os.O_WRONLY, 0644)
	err = CompressPath(tempDir, testBuf)
	if err != nil {
		t.Fatal(err)
	}
	testBuf.Close()
	exec.Command("tar", "-xvf", tarDir+"/test2.tar.gz", "-C", tarDir).Run()
	if _, err := os.Stat(tarDir + "/" + path.Base(tempDir) + "/test1.txt"); os.IsNotExist(err) {
		t.Fatal(err)
	}
	if _, err := os.Stat(tarDir + "/" + path.Base(tempDir) + "/test2.txt"); os.IsNotExist(err) {
		t.Fatal(err)
	}
}
