// Copyright (c) 2012 - Cloud Instruments Co. Ltd.
// 
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met: 
// 
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer. 
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution. 
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package seelog

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const (
	shortPath = "common_context_test.go"
)

var testFullPath string

func fullPath(t *testing.T) string {
	if testFullPath == "" {
		wd, err := os.Getwd()

		if err != nil {
			t.Fatalf("Cannot get working directory: %s", err.Error())
		}

		testFullPath = filepath.Join(wd, shortPath)
	}

	return testFullPath
}

type s struct{}

func getCommonPrefix() string {
	var v1 = reflect.TypeOf((*s)(nil)).Elem()
	sv1 := v1.String()[:strings.LastIndex(v1.String(), ".")+1]
	return sv1
}

func TestContext(t *testing.T) {
	context, err := currentContext()

	nameFunc := getCommonPrefix() + "TestContext"

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	if context == nil {
		t.Fatalf("Expected: context != nil")
	}

	if context.Func() != nameFunc {
		t.Errorf("Expected context.Func == %s ; got %s", nameFunc, context.Func())
	}

	if context.ShortPath() != shortPath {
		t.Errorf("Expected context.ShortPath == %s ; got %s", shortPath, context.ShortPath())
	}

	fp := fullPath(t)

	if context.FullPath() != fp {
		t.Errorf("Expected context.FullPath == %s ; got %s", fp, context.FullPath())
	}
}

func innerContext() (context logContextInterface, err error) {
	return currentContext()
}

func TestInnerContext(t *testing.T) {
	context, err := innerContext()

	nameFunc := getCommonPrefix() + "innerContext"

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	if context == nil {
		t.Fatalf("Expected: context != nil")
	}

	if context.Func() != nameFunc {
		t.Errorf("Expected context.Func == %s ; got %s", nameFunc, context.Func())
	}

	if context.ShortPath() != shortPath {
		t.Errorf("Expected context.ShortPath == %s ; got %s", shortPath, context.ShortPath())
	}

	fp := fullPath(t)

	if context.FullPath() != fp {
		t.Errorf("Expected context.FullPath == %s ; got %s", fp, context.FullPath())
	}
}
