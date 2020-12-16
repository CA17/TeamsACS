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
	"io/ioutil"
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
func (a *AzureBlob) UploadFile(containerName, filepath string) (azblob.CommonResponse, error) {
	ctx := context.Background()
	containerURL, err := a.GetContainerURL(ctx, containerName)
	if err != nil {
		return nil, err
	}
	// Create a file to test the upload and download.
	// log.Info("Creating a dummy file to test the upload and download\n")
	// data := []byte(text)
	// err = ioutil.WriteFile(filename, data, 0700)
	// return filterErrors(err)

	if !common.FileExists(filepath) {
		return nil, fmt.Errorf("file %s not exists", filepath)
	}

	blobURL := containerURL.NewBlockBlobURL(path.Base(filepath))
	file, err := os.Open(filepath)
	err = filterErrors(err)
	if err != nil {
		return nil, err
	}
	// You can use the low-level PutBlob API to upload files. Low-level APIs are simple wrappers for the Azure Storage REST APIs.
	// Note that PutBlob can upload up to 256MB data in one shot. Details: https://docs.microsoft.com/en-us/rest/api/storageservices/put-blob
	// Following is commented out intentionally because we will instead use UploadFileToBlockBlob API to upload the blob
	// _, err = blobURL.PutBlob(ctx, file, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	// filterErrors(err)

	// The high-level API UploadFileToBlockBlob function uploads blocks in parallel for optimal performance, and can handle large files as well.
	// This function calls PutBlock/PutBlockList for files larger 256 MBs, and calls PutBlob for any file smaller
	log.Infof("Uploading the file with blob name: %s\n", path.Base(filepath))
	opt := azblob.UploadToBlockBlobOptions{BlockSize: 4 * 1024 * 1024, Parallelism: 16}
	r, err := azblob.UploadFileToBlockBlob(ctx, file, blobURL, opt)
	err = filterErrors(err)
	if err != nil {
		return nil, err
	}
	return r, nil
}


// GetContainerURL
func (a *AzureBlob) GetContainerURL(ctx context.Context, containerName string) (*azblob.ContainerURL, error) {
	if len(containerName) == 0 {
		return nil,  errors.New("containerName is empty")
	}
	credential, err := azblob.NewSharedKeyCredential(a.AccountName, a.AccountKey)
	if err != nil {
		return nil,fmt.Errorf("invalid credentials with error: %s", err.Error())
	}
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", a.AccountName, containerName))
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	containerURL := azblob.NewContainerURL(*URL, p)
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	err = filterErrors(err)
	if err != nil {
		return  nil, err
	}
	return &containerURL, nil
}



func main() {
	fmt.Printf("Azure Blob storage quick start sample\n")

	// From the Azure portal, get your storage account name and key and set environment variables.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a random string for the quick start container
	containerName := fmt.Sprintf("quickstart-%s", randomString())

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)

	// Create the container
	fmt.Printf("Creating a container named %s\n", containerName)
	ctx := context.Background() // This example uses a never-expiring context
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	filterErrors(err)

	// Create a file to test the upload and download.
	fmt.Printf("Creating a dummy file to test the upload and download\n")
	data := []byte("hello world this is a blob\n")
	fileName := randomString()
	err = ioutil.WriteFile(fileName, data, 0700)
	filterErrors(err)

	// Here's how to upload a blob.
	blobURL := containerURL.NewBlockBlobURL(fileName)
	file, err := os.Open(fileName)
	filterErrors(err)

	// You can use the low-level PutBlob API to upload files. Low-level APIs are simple wrappers for the Azure Storage REST APIs.
	// Note that PutBlob can upload up to 256MB data in one shot. Details: https://docs.microsoft.com/en-us/rest/api/storageservices/put-blob
	// Following is commented out intentionally because we will instead use UploadFileToBlockBlob API to upload the blob
	// _, err = blobURL.PutBlob(ctx, file, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	// filterErrors(err)

	// The high-level API UploadFileToBlockBlob function uploads blocks in parallel for optimal performance, and can handle large files as well.
	// This function calls PutBlock/PutBlockList for files larger 256 MBs, and calls PutBlob for any file smaller
	fmt.Printf("Uploading the file with blob name: %s\n", fileName)
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	filterErrors(err)

	// List the container that we have created above
	fmt.Println("Listing the blobs in the container:")
	for marker := (azblob.Marker{}); marker.NotDone(); {
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		filterErrors(err)

		// ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			fmt.Print("	Blob name: " + blobInfo.Name + "\n")
		}
	}
	//
	// // Here's how to download the blob
	// downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	//
	// // NOTE: automatically retries are performed if the connection fails
	// bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	//
	// // read the body into a buffer
	// downloadedData := bytes.Buffer{}
	// _, err = downloadedData.ReadFrom(bodyStream)
	// filterErrors(err)
	//
	// // The downloaded blob data is in downloadData's buffer. :Let's print it
	// fmt.Printf("Downloaded the blob: " + downloadedData.String())
	//
	// // Cleaning up the quick start by deleting the container and the file created locally
	// fmt.Printf("Press enter key to delete the sample files, example container, and exit the application.\n")
	// bufio.NewReader(os.Stdin).ReadBytes('\n')
	// fmt.Printf("Cleaning up.\n")
	// containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
	// file.Close()
	// os.Remove(fileName)
}
