/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package log

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// File describes a file that gets rotated daily
type File struct {
	sync.Mutex

	// pathGenerator and pathFormat are 2 ways for generating
	// a name of the file when we need to rotate
	// only one of them should be set
	pathGenerator func(time.Time) string
	pathFormat    string

	Location *time.Location

	// info about currently opened file
	day     int
	path    string
	file    *os.File
	onClose func(path string, didRotate bool)

	// position in the file of last Write or Write2, exposed for tests
	lastWritePos int64
}

func (f *File) close(didRotate bool) error {
	if f.file == nil {
		return nil
	}
	err := f.file.Close()
	f.file = nil
	if err == nil && f.onClose != nil {
		f.onClose(f.path, didRotate)
	}
	f.day = 0
	return err
}

// Path returns full path of the file we're currently writing to
func (f *File) Path() string {
	f.Lock()
	defer f.Unlock()
	return f.path
}

func (f *File) open() error {
	t := time.Now().In(f.Location)
	if f.pathGenerator != nil {
		f.path = f.pathGenerator(t)
	} else {
		f.path = t.Format(f.pathFormat)
	}
	f.day = t.YearDay()

	// we can't assume that the dir for the file already exists
	dir := filepath.Dir(f.path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// would be easier to open with os.O_APPEND but Seek() doesn't work in that case
	flag := os.O_CREATE | os.O_WRONLY
	f.file, err = os.OpenFile(f.path, flag, 0644)
	if err != nil {
		return err
	}
	_, err = f.file.Seek(0, io.SeekEnd)
	return err
}

// rotate on new day
func (f *File) reopenIfNeeded() error {
	t := time.Now().In(f.Location)
	if t.YearDay() == f.day {
		return nil
	}
	err := f.close(true)
	if err != nil {
		return err
	}
	return f.open()
}

// NewFile creates a new file that will be rotated daily (at midnight in specified location).
// pathFormat is file format accepted by time.Format that will be used to generate
// a name of the file. It should be unique in a given day e.g. 2006-01-02.txt.
// If you need more flexibility, use NewFileWithPathGenerator which accepts a
// function that generates a file path.
// onClose is an optional function that will be called every time existing file
// is closed, either as a result calling Close or due to being rotated.
// didRotate will be true if it was closed due to rotation.
// If onClose() takes a long time, you should do it in a background goroutine
// (it blocks all other operations, including writes)
// Warning: time.Format might format more than you expect e.g.
// time.Now().Format(`/logs/dir-2/2006-01-02.txt`) will change "-2" in "dir-2" to
// current day. For better control over path generation, use NewFileWithPathGenerator
func NewFile(pathFormat string, onClose func(path string, didRotate bool)) (*File, error) {
	return newFile(pathFormat, nil, onClose)
}

// NewFileWithPathGenerator creates a new file that will be rotated daily
// (at midnight in timezone specified by in specified location).
// pathGenerator is a function that will return a path for a daily log file.
// It should be unique in a given day e.g. time.Format of "2006-01-02.txt"
// creates a string unique for the day.
// onClose is an optional function that will be called every time existing file
// is closed, either as a result calling Close or due to being rotated.
// didRotate will be true if it was closed due to rotation.
// If onClose() takes a long time, you should do it in a background goroutine
// (it blocks all other operations, including writes)
func NewFileWithPathGenerator(pathGenerator func(time.Time) string, onClose func(path string, didRotate bool)) (*File, error) {
	return newFile("", pathGenerator, onClose)
}

func newFile(pathFormat string, pathGenerator func(time.Time) string, onClose func(path string, didRotate bool)) (*File, error) {
	f := &File{
		pathFormat:    pathFormat,
		pathGenerator: pathGenerator,
		Location:      time.UTC,
	}
	// force early failure if we can't open the file
	// note that we don't set onClose yet so that it won't get called due to
	// opening/closing the file
	err := f.reopenIfNeeded()
	if err != nil {
		return nil, err
	}
	err = f.close(false)
	if err != nil {
		return nil, err
	}
	f.onClose = onClose
	return f, nil
}

// Close closes the file
func (f *File) Close() error {
	f.Lock()
	defer f.Unlock()
	return f.close(false)
}

func (f *File) write(d []byte, flush bool) (int64, int, error) {
	err := f.reopenIfNeeded()
	if err != nil {
		return 0, 0, err
	}
	f.lastWritePos, err = f.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, 0, err
	}
	n, err := f.file.Write(d)
	if err != nil {
		return 0, n, err
	}
	if flush {
		err = f.file.Sync()
	}
	return f.lastWritePos, n, err
}

// Write writes data to a file
func (f *File) Write(d []byte) (int, error) {
	f.Lock()
	defer f.Unlock()
	_, n, err := f.write(d, false)
	return n, err
}

// Write2 writes data to a file, optionally flushes. To enable users to later
// seek to where the data was written, it returns name of the file where data
// was written, offset at which the data was written, number of bytes and error
func (f *File) Write2(d []byte, flush bool) (string, int64, int, error) {
	f.Lock()
	defer f.Unlock()
	writtenAtPos, n, err := f.write(d, flush)
	return f.path, writtenAtPos, n, err
}

// Flush flushes the file
func (f *File) Flush() error {
	f.Lock()
	defer f.Unlock()
	return f.file.Sync()
}
