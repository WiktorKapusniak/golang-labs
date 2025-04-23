package main

import (
	"errors"
	"fmt"
	"time"
)

// INTERFEJSY

type FileSystemItem interface {
	Name() string
	Path() string
	Size() int64
	CreatedAt() time.Time
	ModifiedAt() time.Time
}

type Readable interface {
	Read(p []byte) (n int, err error)
}

type Writable interface {
	Write(p []byte) (n int, err error)
}

type Directory interface {
	FileSystemItem
	AddItem(item FileSystemItem) error
	RemoveItem(name string) error
	Items() []FileSystemItem
}

// BŁĘDY 

var (
	ErrItemExists       = errors.New("item already exists")
	ErrItemNotFound     = errors.New("item not found")
	ErrNotImplemented   = errors.New("operation not implemented")
	ErrPermissionDenied = errors.New("permission denied")
	ErrNotDirectory     = errors.New("not a directory")
	ErrIsDirectory      = errors.New("is a directory")
)

// STRUKTURA FILE 

type File struct {
	name       string
	path       string
	createdAt  time.Time
	modifiedAt time.Time
	content    []byte
}

func (f *File) Name() string          { return f.name }
func (f *File) Path() string          { return f.path }
func (f *File) Size() int64           { return int64(len(f.content)) }
func (f *File) CreatedAt() time.Time  { return f.createdAt }
func (f *File) ModifiedAt() time.Time { return f.modifiedAt }

func (f *File) Read(p []byte) (int, error) {
	n := copy(p, f.content)
	return n, nil
}

func (f *File) Write(p []byte) (int, error) {
	f.content = append(f.content, p...)
	f.modifiedAt = time.Now()
	return len(p), nil
}

// STRUKTURA READONLY FILE

type ReadOnlyFile struct {
	File
}

func (r *ReadOnlyFile) Write(p []byte) (int, error) {
	return 0, ErrPermissionDenied
}

// STRUKTURA DIRECTORY

type DirectoryImpl struct {
	name       string
	path       string
	createdAt  time.Time
	modifiedAt time.Time
	items      map[string]FileSystemItem
}

func NewDirectory(name, path string) *DirectoryImpl {
	return &DirectoryImpl{
		name:       name,
		path:       path,
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
		items:      make(map[string]FileSystemItem),
	}
}

func (d *DirectoryImpl) Name() string          { return d.name }
func (d *DirectoryImpl) Path() string          { return d.path }
func (d *DirectoryImpl) Size() int64           { return int64(len(d.items)) }
func (d *DirectoryImpl) CreatedAt() time.Time  { return d.createdAt }
func (d *DirectoryImpl) ModifiedAt() time.Time { return d.modifiedAt }

func (d *DirectoryImpl) AddItem(item FileSystemItem) error {
	if _, exists := d.items[item.Name()]; exists {
		return ErrItemExists
	}
	d.items[item.Name()] = item
	d.modifiedAt = time.Now()
	return nil
}

func (d *DirectoryImpl) RemoveItem(name string) error {
	if _, exists := d.items[name]; !exists {
		return ErrItemNotFound
	}
	delete(d.items, name)
	d.modifiedAt = time.Now()
	return nil
}

func (d *DirectoryImpl) Items() []FileSystemItem {
	var result []FileSystemItem
	for _, item := range d.items {
		result = append(result, item)
	}
	return result
}

// STRUKTURA SYMLINK

type SymLink struct {
	name       string
	path       string
	createdAt  time.Time
	modifiedAt time.Time
	target     FileSystemItem
}

func (s *SymLink) Name() string          { return s.name }
func (s *SymLink) Path() string          { return s.path }
func (s *SymLink) Size() int64           { return s.target.Size() }
func (s *SymLink) CreatedAt() time.Time  { return s.createdAt }
func (s *SymLink) ModifiedAt() time.Time { return s.modifiedAt }



// VIRTUAL FILE SYSTEM

type VirtualFileSystem struct {
	root *DirectoryImpl
}

func NewVirtualFileSystem() *VirtualFileSystem {
	return &VirtualFileSystem{
		root: NewDirectory("root", "/"),
	}
}

func (vfs *VirtualFileSystem) CreateFile(name string, path string, readOnly bool) (FileSystemItem, error) {
	var file FileSystemItem
	if readOnly {
		file = &ReadOnlyFile{File{name, path, time.Now(), time.Now(), []byte{}}}
	} else {
		file = &File{name, path, time.Now(), time.Now(), []byte{}}
	}
	return file, vfs.root.AddItem(file)
}

func (vfs *VirtualFileSystem) CreateDirectory(name string, path string) (*DirectoryImpl, error) {
	dir := NewDirectory(name, path)
	return dir, vfs.root.AddItem(dir)
}

func (vfs *VirtualFileSystem) CreateSymLink(name string, path string, target FileSystemItem) (*SymLink, error) {
	link := &SymLink{name, path, time.Now(), time.Now(), target}
	return link, vfs.root.AddItem(link)
}

func (vfs *VirtualFileSystem) FindItem(name string) (FileSystemItem, error) {
	for _, item := range vfs.root.Items() {
		if item.Name() == name {
			return item, nil
		}
	}
	return nil, ErrItemNotFound
}

func (vfs *VirtualFileSystem) DeleteItem(name string) error {
	return vfs.root.RemoveItem(name)
}

// PRZYKŁADOWE UŻYCIE

func main() {
	// Tworzymy katalog główny
	root := NewDirectory("root", "/")

	// Tworzymy plik i zapisujemy dane
	file := &File{
		name:       "plik.txt",
		path:       "/plik.txt",
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
	}
	file.Write([]byte("Hello, world!"))
	root.AddItem(file)

	// Tworzymy plik tylko do odczytu
	readonly := &ReadOnlyFile{
		File{
			name:       "readme.md",
			path:       "/readme.md",
			createdAt:  time.Now(),
			modifiedAt: time.Now(),
			content:    []byte("To jest tylko do odczytu."),
		},
	}
	root.AddItem(readonly)

	// Tworzymy symlink do pliku
	link := &SymLink{
		name:       "link-do-pliku",
		path:       "/link-do-pliku",
		createdAt:  time.Now(),
		modifiedAt: time.Now(),
		target:     file,
	}
	root.AddItem(link)

	// Wyświetlamy zawartość katalogu głównego
	fmt.Println("Zawartość katalogu root:")
	for _, item := range root.Items() {
		fmt.Printf("- %s (%s) - %d bajtów\n", item.Name(), item.Path(), item.Size())
	}
}

