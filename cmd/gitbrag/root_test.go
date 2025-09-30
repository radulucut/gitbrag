package gitbrag

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/radulucut/gitbrag/internal"
	"github.com/radulucut/gitbrag/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Default(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	testDir := path.Join(os.TempDir(), "test_gitbrag_"+t.Name())
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Fatal(err)
		}
	}()

	err = createGitRepo(testDir)
	if err != nil {
		t.Fatal(err)
	}

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	core := internal.NewCore(timeMock, printer)
	root, err := NewRoot("0.1.0", timeMock, printer, core)
	if err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"gitbrag", testDir}

	err = root.Cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, " 2 files changed\n11 insertions(+)\n 1 deletions(-)\n", out.String())
}
