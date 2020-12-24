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

package azureblob

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
)

type AzureBlob struct {
	AccountName string
	AccountKey  string
}

func NewAzureBlob(accountName string, accountKey string) *AzureBlob {
	ab := &AzureBlob{AccountName: accountName, AccountKey: accountKey}
	return ab
}

func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}

// filterErrors
// ignore ServiceCodeContainerAlreadyExists
func filterErrors(err error) error {
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok { // This error is a Service-specific
			switch serr.ServiceCode() { // Compare serviceCode to ServiceCodeXxx constants
			case azblob.ServiceCodeContainerAlreadyExists:
				log.Error("Received 409. Container already exists")
				return nil
			}
		}
		return err
	}
	return nil
}

// UploadFile
// upload file to azure blob
func (a *AzureBlob) UploadFile(containerName, targetFilepath, filepath string) (azblob.CommonResponse, error) {
	if !common.FileExists(filepath) {
		return nil, fmt.Errorf("file %s not exists", filepath)
	}
	file, err := os.Open(filepath)
	err = filterErrors(err)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return a.UploadFileObject(containerName, targetFilepath, file)
}

// UploadFileObject
// upload file to azure blob
func (a *AzureBlob) UploadFileObject(containerName, targetFilepath string, file *os.File) (azblob.CommonResponse, error) {
	ctx := context.Background()
	containerURL, err := a.GetContainerURL(ctx, containerName)
	if err != nil {
		return nil, err
	}

	blobURL := containerURL.NewBlockBlobURL(path.Base(targetFilepath))
	// You can use the low-level PutBlob API to upload files. Low-level APIs are simple wrappers for the Azure Storage REST APIs.
	// Note that PutBlob can upload up to 256MB data in one shot. Details: https://docs.microsoft.com/en-us/rest/api/storageservices/put-blob
	// Following is commented out intentionally because we will instead use UploadFileToBlockBlob API to upload the blob
	// _, err = blobURL.PutBlob(ctx, file, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	// filterErrors(err)

	// The high-level API UploadFileToBlockBlob function uploads blocks in parallel for optimal performance, and can handle large files as well.
	// This function calls PutBlock/PutBlockList for files larger 256 MBs, and calls PutBlob for any file smaller
	log.Infof("Uploading the file with blob name: %s\n", path.Base(targetFilepath))
	opt := azblob.UploadToBlockBlobOptions{BlockSize: 4 * 1024 * 1024, Parallelism: 16}
	r, err := azblob.UploadFileToBlockBlob(ctx, file, blobURL, opt)
	err = filterErrors(err)
	if err != nil {
		return nil, err
	}
	return r, nil
}




type FileItem struct {
	Container        string
	Filename         string
	ContentType      string
	ContentLength    int64
	ContentEncoding  string
	CreationTime     time.Time
	LastModified     time.Time
	VersionID        string
	IsCurrentVersion bool
	Metadata         map[string]string
}


// ListFiles
// Query the list of all file objects
func (a *AzureBlob) ListFiles(containerName string) ([]FileItem, error) {
	result := make([]FileItem, 0)
	ctx := context.Background()
	containerURL, err := a.GetContainerURL(ctx, containerName)
	if err != nil {
		return nil, err
	}

	for marker := (azblob.Marker{}); marker.NotDone(); {
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		err = filterErrors(err)
		if err != nil {
			return nil, err
		}
		// ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker
		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			fitem := FileItem{
				Container:       containerName,
				Filename:        blobInfo.Name,
				ContentType:     common.GetPointString(blobInfo.Properties.ContentType),
				ContentLength:   common.GetPointInt64(blobInfo.Properties.ContentLength),
				ContentEncoding: common.GetPointString(blobInfo.Properties.ContentEncoding),
				CreationTime:    common.GetPointTime(blobInfo.Properties.CreationTime),
				VersionID:       common.GetPointString(blobInfo.VersionID),
				IsCurrentVersion: common.GetPointBool(blobInfo.IsCurrentVersion),
				LastModified:    blobInfo.Properties.LastModified,
				Metadata:        blobInfo.Metadata,
			}
			result = append(result, fitem)
		}
	}
	return result, nil

}

// GetContainerURL
// create or get container access url
func (a *AzureBlob) GetContainerURL(ctx context.Context, containerName string) (*azblob.ContainerURL, error) {
	if len(containerName) == 0 {
		return nil, errors.New("containerName is empty")
	}
	credential, err := azblob.NewSharedKeyCredential(a.AccountName, a.AccountKey)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials with error: %s", err.Error())
	}
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", a.AccountName, containerName))
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	containerURL := azblob.NewContainerURL(*URL, p)
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	err = filterErrors(err)
	if err != nil {
		return nil, err
	}
	return &containerURL, nil
}

