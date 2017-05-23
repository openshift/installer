package terraform

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"strings"

	"github.com/kardianos/osext"
)

// sources lists the files/directories that compose the TerraForm sources.
//
// The separator enables us to make the distinction between finding a directory
// and a file.
//
// NOTE: Any changes to this list must be reflected in the release script.
// TODO: Ideally, we would have only one source of truth. We could maybe use
// a LDFLAG to fill this for example.
var sources = []string{
	"modules" + string(filepath.Separator),
	"platforms" + string(filepath.Separator),
	"config.tf",
}

// RestoreSources locates the TerraForm sources for Tectonic and copy them to
// the given directory.
//
// It recursively searches the directory trees upward, for all the defined
// sources' files and directories, from to the installer's executable first and
// then from cwd if it didn't already find them.
func RestoreSources(dst string) error {
	// Look into the executable's folder.
	if execFolderPath, err := osext.ExecutableFolder(); err == nil {
		if src, _ := findSources(execFolderPath); src != "" {
			return copySources(src, dst)
		}
	}

	// Look into cwd.
	if workingDirectory, err := os.Getwd(); err == nil {
		if src, _ := findSources(workingDirectory); src != "" {
			return copySources(src, dst)
		}
	}

	return errors.New("could not find TerraForm sources")
}

// findSources recursively looks for the `platforms` and `modules` directory,
// browsing the tree upward from the given root.
func findSources(root string) (string, error) {
	// Check if all the sources are here. Return that root if it's the case.
	found := true
	for _, src := range sources {
		fi, err := os.Stat(filepath.Join(root, src))
		if err != nil || (os.IsPathSeparator(src[len(src)-1]) && !fi.IsDir()) {
			// The source does not exist, or we were expecting a directory and got
			// something else.
			found = false
			break
		}
	}

	if found {
		return root, nil
	}

	// Recurse.
	newRoot, err := filepath.Abs(filepath.Join(root, ".."))
	if err != nil || newRoot == root {
		return "", nil
	}
	return findSources(newRoot)
}

// copySources copies all the sources present in the root directory to the
// destination folder.
func copySources(root, dst string) error {
	for _, src := range sources {
		// We copy the source directory itself to the destination, rather than its
		// content.
		src = strings.TrimSuffix(src, string(filepath.Separator))

		if err := copyDir(filepath.Join(root, src), dst); err != nil {
			return err
		}
	}
	return nil
}

// copyDir recursively copies a directory, located at src, to dst.
//
// It preserves the permissions, but not the extra metadata.
// It supports regular files, directories and symlinks, but not devices,
// sockets, named pipes, etc.
// The destination can exist, everything will be overwritten as necessary.
//
// If dst is an existing directory, and src does not end with a file separator,
// then the src directory itself is copied into dst, rather than its content.
// If the src is actually a regular file, the file is copied to dst instead. If
// dst is an existing directory, the file will be copied inside.
func copyDir(src, dst string) error {
	if fi, err := os.Stat(src); err == nil && fi.Mode().IsRegular() {
		// src is actually a file, not a directory, copy it and return.
		return copyFile(src, fi.Mode(), dst)
	}

	if fi, err := os.Stat(dst); err == nil {
		if fi.Mode().IsDir() && !os.IsPathSeparator(src[len(src)-1]) {
			// dst is a directory and src does not end with a file separator, copy
			// the src directory itself rather than its content.
			_, srcDirName := filepath.Split(src)
			dst = filepath.Join(dst, srcDirName)
		}
	}

	return filepath.Walk(src, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Build the path to the destination using the relative path of the current
		// file regarding to the src.
		relp, err := filepath.Rel(src, p)
		if err != nil {
			return fmt.Errorf("could not determine relative path of %s: %s", p, err)
		}
		dstp := filepath.Join(dst, relp)

		switch {
		case fi.Mode().IsDir():
			// Directory.
			if err := os.MkdirAll(dstp, fi.Mode()); err != nil {
				return fmt.Errorf("could not make directory %s: %v", p, err)
			}
			return nil
		case fi.Mode().IsRegular():
			// File.
			if fi.Mode().IsRegular() {
				if err := copyFile(p, fi.Mode(), dstp); err != nil {
					return fmt.Errorf("could not copy file %s: %v", p, err)
				}
			}
		case fi.Mode()&os.ModeSymlink != 0:
			// Symlink.
			ldstp, err := os.Readlink(p)
			if err != nil {
				return fmt.Errorf("could not determine symlink's destination: %v", err)
			}

			os.Remove(dstp)
			if err := os.Symlink(ldstp, dstp); err != nil {
				return fmt.Errorf("could not create symlink %s: %v", p, err)
			}
		default:
			// Device, Socket, NamedPipe, ...
			return fmt.Errorf("could not copy file of type %v", fi.Mode())
		}

		return nil
	})
}

// copyFile copies a single file, located at src, to dst, with
// the provided permissions. If dst is an existing directory, the file is
// written inside it.
//
// The destination will be created or overwritten if necessary.
func copyFile(src string, srcMode os.FileMode, dst string) error {
	if fi, err := os.Stat(dst); err == nil && fi.IsDir() {
		// The destination is an existing folder, write the file into it.
		_, srcFileName := filepath.Split(src)
		dst = filepath.Join(dst, srcFileName)
	}

	// Open the source file.
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Open the destination file.
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcMode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Write the file content.
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
