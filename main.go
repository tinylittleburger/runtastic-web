package main

import (
	"archive/zip"
	"context"
	"fmt"
	"net/http"
	"time"

	"io"

	"github.com/metalnem/runtastic/api"
	"github.com/metalnem/runtastic/tcx"
)

func formatTime(t time.Time) string {
	return t.Local().Format("2006-01-02 15.04.05")
}

func export(activities []api.Activity, w io.Writer) error {
	zw := zip.NewWriter(w)
	defer zw.Close()

	for _, activity := range activities {
		t := formatTime(activity.EndTime)
		filename := fmt.Sprintf("Runtastic %s %s.tcx", t, activity.Type.DisplayName)

		header := zip.FileHeader{
			Name:   filename,
			Method: zip.Deflate,
		}

		header.SetModTime(time.Now())
		file, err := zw.CreateHeader(&header)

		if err != nil {
			return err
		}

		exp := tcx.NewExporter(file)

		if err = exp.Export(activity); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		api.Login(context.Background(), "", "")
	})

	http.ListenAndServe(":80", nil)
}
