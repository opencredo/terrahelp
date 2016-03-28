package terrahelp

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testProject struct {
	t          *testing.T
	origWd     string
	baseSrcDir string
	baseTmpDir string
}

func newTempProject(t *testing.T) *testProject {
	currDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error trying to get currWd: %v", err)
	}
	tmpDir, err := ioutil.TempDir("", "terrahelp")
	if err != nil {
		t.Fatalf("error trying to create tmp Dir: %v", err)
	}
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("error trying to change to tmp Dir: %v", err)
	}
	return &testProject{t, currDir, baseSrcCodePackageDir(), tmpDir}
}

func (c *testProject) getSrcFile(f string) string {
	return path.Join(c.baseSrcDir, f)
}
func (c *testProject) getProjectFile(f string) string {
	return path.Join(c.baseTmpDir, f)
}
func (c *testProject) removeProjectFile(f string) {
	err := os.Remove(f)
	if err != nil {
		c.t.Fatalf("Unable to remove file %s : %s", f, err)
	}
}

func (c *testProject) restore() {
	e1 := os.Chdir(c.origWd)
	e2 := os.RemoveAll(c.baseTmpDir)
	if e1 != nil || e2 != nil {
		c.t.Fatalf("error trying to destroy temp project: %s or %s", e1, e2)
	}
}

func (c *testProject) copyExampleProject(ver string) {
	err := CopyDir(path.Join(c.baseSrcDir, "test-data", "example-project", ver),
		c.baseTmpDir)
	if err != nil {
		c.t.Fatalf("Issue trying to setup tmp dir for testing %s", err)
	}
}

func (c *testProject) assertExpectedFileContent(tmpFile, expFile string) {
	actual := getFileContents(c.t, c.baseTmpDir, tmpFile)
	expected := getFileContents(c.t, c.baseSrcDir, expFile)
	assert.Equal(c.t, expected, actual, "Expected file content for %s not as expected", tmpFile)

}

func assertFileDoesNotExist(t *testing.T, f string) {
	if _, err := os.Stat(f); err == nil {
		t.Fatalf("File %s exists and should not", f)
	}
}

func baseSrcCodePackageDir() string {
	_, file, _, _ := runtime.Caller(0)
	return path.Dir(file)
}

func getFileContents(t *testing.T, baseDir, file string) string {
	b, err := ioutil.ReadFile(path.Join(baseDir, file))
	if err != nil {
		t.Fatalf("Unable to read tmpFileContents %s", err)
	}
	return string(b)
}
