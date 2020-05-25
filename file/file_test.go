package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {

	t.Run("selfPath", func(t *testing.T) {
		path := SelfPath()
		assert.Equal(t, "C:\\Users\\28423\\AppData\\Local\\Temp\\___go_test_github_com_Muskchen_toolkits_file.exe", path)
	})

	t.Run("RealPath", func(t *testing.T) {
		path, err := RealPath("./")
		if err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, "E:\\git\\gopath\\src\\github.com\\Muskchen\\toolkits\\file", path)
	})

	t.Run("selfDir", func(t *testing.T) {
		dir := SelfDir()
		assert.Equal(t, "C:\\Users\\28423\\AppData\\Local\\Temp", dir)
	})

	t.Run("baseName", func(t *testing.T) {
		basename := Basename("./file.go")
		assert.Equal(t, "file.go", basename)
	})

	t.Run("dir", func(t *testing.T) {
		dir := Dir("/etc/test.go")
		assert.Equal(t, "/etc", dir)
	})
}
