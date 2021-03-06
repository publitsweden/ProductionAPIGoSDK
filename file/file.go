// Copyright 2017 Publit Sweden AB. All rights reserved.

// Handles Publit File data.
//
// This package has methods to show and index Files from the Publit API.
// It also conatins methods to manipule List of files such as downloading them.
// See the examples for implementational information.
package file

import (
	"errors"
	"fmt"
	"github.com/publitsweden/APIUtilityGoSDK/common"
	"github.com/publitsweden/ProductionAPIGoSDK"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// File attribute constants.
const (
	ID             = "id"
	TYPE           = "type"
	ORIGINAL_NAME  = "original_name"
	SIZE           = "size"
	EXTENSION      = "extension"
	MIME           = "mime"
	CHECKSUM       = "checksum"
	URL            = "url"
	AUTO_GENERATED = "auto_generated"
	CREATED_AT     = "created_at"
	UPDATED_AT     = "updated_at"
	DELETED_AT     = "deleted_at"
	PRESIGNED      = "presigned"
)

// Amount of workers for concurrent jobs.
var workerAmount int = 5

const (
	AUX_PRESIGNED = "presigned_url"
)

// Resource struct.
type Resource struct {
	Endpoint Endpoint
	Id       int
}

// Endpoint enumeration type.
type Endpoint int

// Endpoint enumeration constants.
const (
	INDEX Endpoint = 1 + iota
	SHOW
)

// Endpoints.
var endpoints map[Endpoint]string = map[Endpoint]string{
	INDEX: "files",
	SHOW:  "files/%v",
}

// FileList type for handling collections of files.
type FileList []*File

// Holds file information based on the Publit production APIs "files" resource response.
type File struct {
	ID            int               `json:"id,string,omitempty"`
	Type          string            `json:"type"`
	OriginalName  string            `json:"original_name"`
	Size          int               `json:"size,string"`
	Extension     string            `json:"extension"`
	Mime          string            `json:"mime_type"`
	Checksum      string            `json:"checksum"`
	URL           string            `json:"url"`
	AutoGenerated int               `json:"auto_generated,string"`
	CreatedAt     common.PublitTime `json:"created_at,omitempty"`
	UpdatedAt     common.PublitTime `json:"updatd_at,omitempty"`
	DeletedAt     common.PublitTime `json:"deleted_at,omitempty"`
	Presigned     string            `json:"presigned_url,omitempty"`
}

// ProductionAPIGetter defines how the client should perform GET calls.
type ProductionAPIGetter interface {
	Get(endpoint production.Endpointer, model interface{}, queryParams ...func(q url.Values)) error
}

// Index response object.
type IndexResponse struct {
	Count int      `json:"count"`
	Next  string   `json:"next"`
	Prev  string   `json:"prev"`
	Data  FileList `json:"data"`
}

// Returns File from Publit API.
func Show(c ProductionAPIGetter, id int, queryParams ...func(q url.Values)) (*File, error) {
	f := &File{}
	r := Resource{Endpoint: SHOW, Id: id}
	err := c.Get(r, f, queryParams...)
	return f, err
}

// Indexes Files from the Publit API.
func Index(c ProductionAPIGetter, queryParams ...func(q url.Values)) (*IndexResponse, error) {
	ir := &IndexResponse{}
	r := Resource{Endpoint: INDEX}
	err := c.Get(r, ir, queryParams...)
	return ir, err
}

// Retrieves presigned URLs for list of files. Uses worker concurrency pattern.
// Function returns a map indexed on FileId and any errors if they have occured.
func (fl FileList) GetPresigned(c ProductionAPIGetter) map[int]error {
	jobs := make(chan *File, len(fl))
	results := make(chan FileWorkerError, len(fl))

	// Create workers.
	for i := 0; i < workerAmount; i++ {
		go presignedWorker(c, jobs, results)
	}

	// Add jobs to channel.
	for _, v := range fl {
		jobs <- v
	}
	// Close channel to indicate all jobs have been pushed.
	close(jobs)

	// Create error array.
	errs := make(map[int]error, len(fl))

	// Range results to be sure to wait until all workers have completed their task.
	// Also store their results in an array.
	for range fl {
		err := <-results
		errs[err.FileId] = err.Error
	}

	return errs
}

// Holds file worker errors. Useful for communicating via channels.
type FileWorkerError struct {
	Error  error
	FileId int
}

// Worker designated for running presigned url queries.
func presignedWorker(c ProductionAPIGetter, files <-chan *File, results chan<- FileWorkerError) {
	for f := range files {
		err := f.GetPresignedUrl(c)
		fw := FileWorkerError{
			Error:  err,
			FileId: f.ID,
		}
		results <- fw
	}
}

// Retrieves presigned url for file.
// The presigned url is a download url valid for a certain amount of time.
func (f *File) GetPresignedUrl(c ProductionAPIGetter) error {
	r := Resource{Endpoint: SHOW, Id: f.ID}
	err := c.Get(r, f, GetPresignedAuxParamFunc())
	return err
}

// Creates presigned query param function.
func GetPresignedAuxParamFunc() func(q url.Values) {
	return func(q url.Values) { q.Set(common.QUERY_KEY_AUX, AUX_PRESIGNED) }
}

// Method to Resource that fullfils the Enpointer interface as stated in production.
func (r Resource) GetEndpoint() string {
	e := endpoints[r.Endpoint]

	end := e
	if strings.Contains(e, "%v") && r.Id != 0 {
		end = fmt.Sprintf(e, r.Id)
	}

	return end
}

// Sets number of workers for methods using concurrent workers.
func SetWorkerAmount(amount int) {
	workerAmount = amount
}

// Downloads file from FileList.
// Return map of errors indexed on fileID and potential error generated by the method.
func (fl FileList) DownloadFiles(c ProductionAPIGetter, outDir string) (map[int]error, error) {
	errs := make(map[int]error, len(fl))
	if stat, err := os.Stat(outDir); err != nil || !stat.IsDir() {
		return errs, errors.New("Output dir is not a directory.")
	}

	jobs := make(chan *File, len(fl))
	results := make(chan FileWorkerError, len(fl))

	var noPresigned FileList

	// Range FileList to load presigned_url for any files that doesn't have it.
	for _, v := range fl {
		if v.Presigned == "" {
			noPresigned = append(noPresigned, v)
		}
	}

	// Get presigned for files without presigned urls.
	if len(noPresigned) > 0 {
		preErrs := noPresigned.GetPresigned(c)
		// Range error list and see if any errors were returned.
		// Also "bake" the returned map into the map of this method to make sure map conforms to init length of fl.
		presignedError := false
		for k, v := range preErrs {
			if v != nil {
				errs[k] = v
				presignedError = true
			}
		}
		if presignedError {
			return preErrs, nil
		}
	}

	// Create workers.
	for i := 0; i < workerAmount; i++ {
		go downloadWorker(outDir, jobs, results)
	}

	// Range files and create a download job for each file.
	for _, v := range fl {
		jobs <- v
	}
	close(jobs)

	// Range results.
	for range fl {
		err := <-results
		errs[err.FileId] = err.Error
	}

	return errs, nil
}

// Plain getter method. Performs plain GET requests for file download from URL.
// Made as a variable for aiding testing.
var PlainGetter func(url string) (*http.Response, error) = func(url string) (*http.Response, error) {
	return http.Get(url)
}

// Download worker.
func downloadWorker(outDir string, files <-chan *File, results chan<- FileWorkerError) {
	for f := range files {
		fw := FileWorkerError{
			FileId: f.ID,
		}

		resp, err := PlainGetter(f.Presigned)
		if resp.StatusCode != http.StatusOK {
			fw.Error = errors.New(fmt.Sprintf(`Could not download file. Server responded with code: "%s"`, resp.StatusCode))
			results <- fw
			continue
		}

		defer resp.Body.Close()
		if err != nil {
			fw.Error = err
			results <- fw
			continue
		}
		f, err := os.Create(outDir + "/" + f.OriginalName)
		defer f.Close()
		if err != nil {
			fw.Error = err
			results <- fw
			continue
		}
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			fw.Error = err
			results <- fw
			continue
		}
		results <- fw
	}
}
