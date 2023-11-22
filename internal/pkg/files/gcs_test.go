package files

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/fullstorydev/emulators/storage/gcsemu"
)

func TestGcs_Accepts(t *testing.T) {
	b, err := NewGcs()
	assert.NoError(t, err)

	testCases := []struct {
		path    Path
		accepts bool
	}{
		{"gs://bucket/path/to/file", true},
		{"gs://bucket", false},
		{"gs:bucket/path/to/file", false},
		{"/path/to/file", false},
		{"path/to/file", false},
	}
	for _, tc := range testCases {
		t.Run(string(tc.path), func(t *testing.T) {
			assert.Equal(t, tc.accepts, b.Accepts(tc.path))
		})
	}
}

const testBucket = "testbucket"

func runFakeGcs() (*gcsemu.Server, error) {
	host := "127.0.0.1:7890"

	server, err := gcsemu.NewServer(host, gcsemu.Options{})
	if err != nil {
		return nil, err
	}

	err = os.Setenv("GCS_EMULATOR_HOST", server.Addr)
	if err != nil {
		server.Close()
		return nil, err
	}

	err = server.InitBucket(testBucket)
	if err != nil {
		server.Close()
		return nil, err
	}

	return server, nil
}

func addFakeFile(name string, content string) error {
	gcsClient, err := gcsemu.NewClient(context.Background())
	if err != nil {
		return err
	}
	defer func() { _ = gcsClient.Close() }()

	bucket := gcsClient.Bucket(testBucket)
	object := bucket.Object(name)

	writer := object.NewWriter(context.Background())
	defer func() { _ = writer.Close() }()

	_, err = writer.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}

func TestGcs_Load(t *testing.T) {
	server, err := runFakeGcs()
	assert.NoError(t, err)
	t.Cleanup(func() {
		server.Close()
	})

	gcsClient, err := gcsemu.NewClient(context.Background())
	assert.NoError(t, err)

	b, err := NewGcs(WithGcsClient(gcsClient))
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = b.Close()
	})

	t.Run("should load an existing file", func(t *testing.T) {
		actualContent := "abc"

		err = addFakeFile("foo", actualContent)
		assert.NoError(t, err)

		p := Path(fmt.Sprintf("gs://%s/foo", testBucket))
		f, err := b.Load(p)
		assert.NoError(t, err)
		t.Cleanup(func() {
			_ = f.Close()
		})

		content, err := io.ReadAll(f)
		assert.Equal(t, []byte(actualContent), content)
	})

	t.Run("should error on a missing file", func(t *testing.T) {
		p := Path(fmt.Sprintf("gs://%s/bar", testBucket))
		_, err = b.Load(p)
		assert.Error(t, err)
	})
}
