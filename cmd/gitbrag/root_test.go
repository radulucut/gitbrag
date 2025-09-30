package gitbrag

import (
	"bytes"
	"os"
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

	testDir := createGitRepo(t)

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

func Test_AuthorFlag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	testDir := createGitRepo(t)

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	core := internal.NewCore(timeMock, printer)
	root, err := NewRoot("0.1.0", timeMock, printer, core)
	if err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"gitbrag", testDir, "--author", "John Doe"}

	err = root.Cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	// Only the second commit by John Doe should be counted (1 file changed, 0 insertions, 1 deletion)
	assert.Equal(t, "1 files changed\n0 insertions(+)\n1 deletions(-)\n", out.String())
}

func Test_AuthorFlagByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeMock := mocks.NewMockTime(ctrl)
	timeMock.EXPECT().Now().Return(defaultCurrentTime).AnyTimes()

	testDir := createGitRepo(t)

	out := new(bytes.Buffer)
	printer := internal.NewPrinter(nil, out, out)
	core := internal.NewCore(timeMock, printer)
	root, err := NewRoot("0.1.0", timeMock, printer, core)
	if err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"gitbrag", testDir, "--author", "john.doe@example.com"}

	err = root.Cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
	// Only the second commit by john.doe@example.com should be counted (1 file changed, 0 insertions, 1 deletion)
	assert.Equal(t, "1 files changed\n0 insertions(+)\n1 deletions(-)\n", out.String())
}
