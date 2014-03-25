package util_test

import (
	"github.com/sellweek/TOGY/util"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"
)

var conf = quick.Config{
	Values: values,
}

func values(args []reflect.Value, r *rand.Rand) {
	args[0] = reflect.ValueOf(r.Intn(10000))
	args[1] = reflect.ValueOf(r.Intn(13))
	args[2] = reflect.ValueOf(r.Intn(32))
	args[3] = reflect.ValueOf(r.Intn(24))
	args[4] = reflect.ValueOf(r.Intn(60))
	args[5] = reflect.ValueOf(r.Intn(60))
	args[6] = reflect.ValueOf(r.Intn(1000000000))
}

func TestNormalizeDate(t *testing.T) {
	t.Parallel()
	testf := func(y, m, d, h, min, s, n int) bool {
		date := time.Date(y, time.Month(m), d, h, min, s, n, time.UTC)
		norm := util.NormalizeDate(date)
		hm := (norm.Hour() == 0) && (norm.Minute() == 0)
		sn := (norm.Second() == 0) && (norm.Nanosecond() == 0)
		tz := norm.Location() == util.Tz
		return hm && sn && tz
	}

	err := quick.Check(testf, &conf)
	if err != nil {
		t.Error(err)
	}
}

func TestNormalizeTime(t *testing.T) {
	t.Parallel()
	testf := func(y, m, d, h, min, s, n int) bool {
		date := time.Date(y, time.Month(m), d, h, min, s, n, time.UTC)
		norm := util.NormalizeTime(date)
		return (norm.Year() == 0) && (norm.Month() == 1) && (norm.Day() == 1) && (norm.Location() == util.Tz)
	}

	err := quick.Check(testf, &conf)
	if err != nil {
		t.Error(err)
	}
}

func TestFileType(t *testing.T) {
	t.Parallel()
	if util.GetFileType("hello.pptx") != "pptx" {
		t.Fatal("GetFileType can't get a simple file type.")
	}

	if util.GetFileType("hello.pptx.tar.gz") != "gz" {
		t.Fatal("GetFileType can't get the last file type.")
	}

	if util.GetFileType("README") != "" {
		t.Fatal("GetFile gets a file extension when there is none.")
	}
}
