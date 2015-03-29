package bindata

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _bash_backup_bash = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x90\xb1\x4e\x33\x31\x10\x84\xfb\x7b\x8a\x51\xfe\xbf\x80\xc2\x76\x41\x05\x54\x01\x51\x20\x24\x40\x10\x44\x11\x45\x91\x73\xb7\x4e\x56\x38\xf6\x61\xef\x05\xa2\x28\xef\x8e\x73\xba\x88\x02\x89\x82\x6e\x35\xf2\xcc\x37\xe3\x4c\x02\x45\x11\x2d\xb7\xe4\x2c\xfb\xaa\x5a\x5b\x0e\x27\xa7\xd8\x55\x00\x3b\x4c\xa7\x50\x0e\x66\x63\x93\xf1\xbc\x30\x6d\xcc\xb2\x4c\x94\xdf\xbd\x39\xd7\x67\xe6\xf0\xd6\x24\xaa\xe3\x86\xd2\x56\xd7\x31\x38\xcc\x66\xb8\x84\xac\x28\x14\x3f\x40\xf5\x2a\x62\xd4\x58\xb1\x0b\x9b\xe9\x02\xc7\x0b\x9c\x51\x77\x29\x51\x10\xbf\xc5\x90\xc0\x61\x09\x97\xe2\x1a\x16\x0b\x5b\xbf\x75\xad\xc6\x2b\x7b\x0f\x49\x5b\xd8\x65\x41\x21\xd0\xa7\x40\x78\x4d\x5a\xeb\x51\x01\x90\xcf\xd4\x73\xfe\xa1\xa5\xe4\x62\xfa\xf6\xf6\x72\xee\x9a\x08\xd5\xe1\x58\x1b\x14\x36\x0d\x27\x18\x92\xda\x7c\x58\xaf\x48\x37\xa6\x68\xe8\xef\xc1\xa9\xda\x2e\xaf\x7e\x9d\x3c\x20\x63\xe8\xbb\xcb\xa1\x59\x59\x0c\x6f\x85\xb2\xe0\x6a\x7c\x7d\xf7\xf2\xf8\x3c\x9f\x3c\xcc\x9f\x6e\x26\xe3\xdb\xfb\x21\x37\xff\xa5\x52\x43\x9e\x84\xa0\xd4\xe1\x73\xb9\xec\x1b\x78\xff\x77\x3f\x30\xfb\x92\xef\xb8\xda\x7f\x05\x00\x00\xff\xff\x9e\xae\xa1\x43\xd3\x01\x00\x00")

func bash_backup_bash_bytes() ([]byte, error) {
	return bindata_read(
		_bash_backup_bash,
		"bash/backup.bash",
	)
}

func bash_backup_bash() (*asset, error) {
	bytes, err := bash_backup_bash_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "bash/backup.bash", size: 467, mode: os.FileMode(493), modTime: time.Unix(1427679730, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _bash_postgres_init_bash = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\x8f\xbd\x4e\xc4\x30\x10\x84\xfb\x3c\xc5\x20\x0a\x40\xc2\xb8\xa0\x02\x5a\x9e\x80\xf6\x74\xc5\x26\x5e\x73\x2b\xf9\xbc\xc1\x76\xf8\x15\xef\xce\x3a\xd2\x89\x26\xba\xce\x2b\xcd\xf7\xcd\xb8\x72\x83\x63\xc5\x2c\x33\x47\x92\x34\x0c\x47\x92\x7c\x7d\x83\x9f\x01\xb8\x84\x64\x69\x42\x49\xbe\x19\x81\x1a\x8d\x54\x19\x12\xa1\xd9\x6e\xe5\x9a\xaf\x1a\x28\x15\xa6\xf0\x05\xfe\x94\xda\x56\x28\x6a\xb1\x8b\x8e\x73\xe2\x5b\x33\xa0\x1d\x18\x53\x27\x35\x82\x56\x0f\x26\xcd\xcd\x7a\xb8\x18\x60\xbe\xdd\x0e\x17\x70\x01\xfe\x9d\x8a\x4f\x32\xfa\x59\x6b\x7b\x2d\x5c\xdf\x92\x7f\xb8\xbb\xf7\x7d\x13\xf6\xfb\xa7\xae\xca\x86\x00\xd3\x41\x3f\x32\xdc\x0b\x4e\xc9\xc7\xd3\x63\x4b\xb2\x22\x75\x09\x0a\xb7\xe0\x3f\xb8\xd4\xcd\xb6\x51\xb2\xef\x1f\x0f\x23\xdc\xf3\xd9\x4d\xe6\x8d\x32\xfc\xfe\x05\x00\x00\xff\xff\x28\xb0\x81\x45\x45\x01\x00\x00")

func bash_postgres_init_bash_bytes() ([]byte, error) {
	return bindata_read(
		_bash_postgres_init_bash,
		"bash/postgres-init.bash",
	)
}

func bash_postgres_init_bash() (*asset, error) {
	bytes, err := bash_postgres_init_bash_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "bash/postgres-init.bash", size: 325, mode: os.FileMode(493), modTime: time.Unix(1427686312, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _bash_postgres_bash = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x8c\x54\x4d\x73\xd4\x38\x10\xbd\xcf\xaf\x78\xeb\xec\x56\x76\xab\xe2\x51\x16\x0e\x14\xa1\x48\x55\xa0\x72\x22\x70\x20\xa1\x38\x24\xa9\x44\x96\xdb\x33\x62\x64\xc9\x48\xf2\x0c\x43\xc8\x7f\xa7\x65\x7b\x1c\x42\x2a\x43\x0e\xf3\x61\xab\x5f\xf7\xeb\xd7\xfd\x14\x28\x22\x27\x87\x46\x37\x54\x49\x6d\x26\x93\x5a\x6a\xfb\xef\x7f\xb8\x99\x00\x3b\x20\x1b\x5a\x4f\xf8\x7c\x74\x02\xe3\x66\x28\x5a\xb5\x60\x00\x7d\xd3\x21\x06\x0e\x20\xbb\x2c\xb5\x87\xa0\xa8\xc4\x4a\x9a\x9c\xa6\xa5\xe0\x77\x10\xb2\x69\x44\xa1\xad\x50\x9e\x64\xa4\xab\x01\xf7\xf7\xcd\x9b\x4f\x6f\xdf\x1d\x9f\x5d\x7d\x38\x7a\x7f\x7c\x3b\xe1\x04\xba\xc2\xf9\x39\xfe\x42\x5e\x41\x2c\xa5\x17\x46\x17\xa2\x71\x21\xce\x3c\x85\xaf\x46\xbc\x9c\x3e\x17\x89\x8f\xd0\x56\x47\x2d\x8d\xfe\x4e\x25\x2e\x2f\x5f\x21\xce\xc9\x32\x9c\x19\xa8\xb9\x43\x56\xca\x28\x0b\x19\xe8\x00\xd6\xf5\xec\xb4\x9d\x61\xf3\x16\x95\x6b\x6d\x39\xcd\x3a\xc0\x0e\xd4\x9c\xd4\x22\x55\xe6\x24\xdc\x9b\x4c\x1f\xbb\x46\x21\xd5\xa2\x6d\x02\xf2\x3c\x9d\x05\xb7\x07\x43\x71\x37\x80\x99\x44\xe7\x69\x00\xaf\x08\xca\xb5\xa6\x44\xe3\x5d\x21\x0b\xb3\x46\xe9\x50\x50\x8c\xe4\x39\x9f\xb4\xf8\xd2\x86\x88\x48\x3d\x03\xdb\xd6\x05\x1f\xb8\x0a\x46\x5b\xea\x72\x3b\x4b\xdd\x03\x74\xe8\x63\x25\xe6\x24\x4b\x8e\xde\x43\x4d\xd2\x76\x30\xb7\x61\xd3\x55\xed\x45\xba\x7e\x4c\xec\xee\x3f\xa7\x66\x0a\xdc\x6b\x0f\xcc\x0d\x6b\x80\x67\x87\x10\x25\x2d\x85\x6d\x8d\xc1\x0f\xac\x14\x72\x73\x8d\x7c\x16\x91\xfd\x9f\xdd\xd7\xf1\xa1\x92\x7d\xdf\x89\x4e\xe5\x5d\x3d\xe4\x9d\x4e\x07\x19\x01\x5f\x23\xf7\xdb\xc7\x36\x44\x86\x96\x35\xca\x5b\x6c\x22\x1e\x5d\x9b\xbe\x93\xa1\x83\x8a\x4f\xe7\x5b\xd3\xe3\xe4\xe8\xec\xf8\xf4\x6c\xa8\xc2\xf4\x57\x16\xf9\xc7\xb1\xcc\xc1\x58\xef\x09\x1c\xd5\xbc\x76\x25\xf6\x5f\xec\xef\x3f\x25\xba\xd7\x6a\xd8\x8c\x2b\xe5\xea\x5a\xda\x12\xaf\xb1\xbb\xbd\xb3\xf4\xdd\xb7\x75\x91\xfd\x53\x5d\x64\xe9\xa7\xb9\xc8\x76\x33\x1e\xce\x03\x91\x22\xd1\x76\x4f\x78\x52\x6e\x49\x7e\x3d\x55\xce\x56\x38\x1c\x27\xdd\xfb\xc2\x04\x7a\x64\xb0\xa3\x97\xd2\x6c\x25\x2c\xad\x46\xa7\x8c\xe3\xad\xf4\xb0\xef\xa4\x93\x4b\x98\xf9\x7a\x2f\xed\x7e\x2d\xfd\x22\x2d\xcd\x9d\xb7\x64\xc0\x2f\xe6\xec\x50\xd1\xb5\x7f\x98\x9c\xb8\x0f\x19\xc9\xfe\x4e\xf5\x9e\x95\xc1\xca\x92\x62\xc5\xd7\x83\xa1\x71\x1a\xa5\xef\x8e\x37\x25\x06\xfe\xcc\xbe\xbb\xbd\x1a\xf2\x95\xe3\x3d\x95\xc9\x75\x79\xd4\x35\xf1\x56\x1b\x27\x4b\xa6\xc8\x98\xa6\x35\x7c\x39\xdd\xb5\x42\x36\x7a\x4d\xc9\x73\xa2\x0d\xcc\xdd\x29\x69\xba\x5b\xac\x07\x4d\x6e\x7f\x06\x00\x00\xff\xff\x3b\x83\x73\xfa\x2d\x05\x00\x00")

func bash_postgres_bash_bytes() ([]byte, error) {
	return bindata_read(
		_bash_postgres_bash,
		"bash/postgres.bash",
	)
}

func bash_postgres_bash() (*asset, error) {
	bytes, err := bash_postgres_bash_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "bash/postgres.bash", size: 1325, mode: os.FileMode(493), modTime: time.Unix(1427686292, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"bash/backup.bash": bash_backup_bash,
	"bash/postgres-init.bash": bash_postgres_init_bash,
	"bash/postgres.bash": bash_postgres_bash,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"bash": &_bintree_t{nil, map[string]*_bintree_t{
		"backup.bash": &_bintree_t{bash_backup_bash, map[string]*_bintree_t{
		}},
		"postgres-init.bash": &_bintree_t{bash_postgres_init_bash, map[string]*_bintree_t{
		}},
		"postgres.bash": &_bintree_t{bash_postgres_bash, map[string]*_bintree_t{
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

