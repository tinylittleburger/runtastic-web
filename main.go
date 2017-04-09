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
		username := "imateapot@mailinator.com"
		password := "imateapot"

		ctx := context.Background()
		session, err := api.Login(ctx, username, password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		activities, err := session.GetActivities(ctx)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(activities) == 0 {
			http.Error(w, "There are no activities to backup", http.StatusBadRequest)
			return
		}

		filename := fmt.Sprintf("Runtastic %s.zip", formatTime(time.Now()))
		header := fmt.Sprintf("attachment; filename=\"%s\"", filename)

		w.Header().Set("Content-Disposition", header)
		w.Header().Set("Content-Type", "application/zip")

		export(activities, w)
	})

	http.ListenAndServe(":80", nil)
}
