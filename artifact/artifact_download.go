package artifact

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type stop struct {
	error
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			return s.error
		}
		if attempts--; attempts > 0 {
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}
	return nil
}

func DownloadBuildArtifact() {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "app-armeabi-v7a-release.apk")
	if err != nil {
		log.Fatalf("Cannot create temporary file %s\n", err)
	}
	defer tmpFile.Close()
	err = retry(3, time.Second, func() error {
		resp, err := http.Get("https://bitrise-prod-build-storage.s3.amazonaws.com/builds/df7d2a57d4272143/artifacts/28806521/app-armeabi-v7a-release.apk?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIV2YZWMVCNWNR2HA%2F20200803%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20200803T193627Z&X-Amz-Expires=43200&X-Amz-SignedHeaders=host&X-Amz-Signature=c8536204fa4df33ad257cd08fe927b6ece5d2aa51a9e42f9052c0628505f0ba0")
		if err != nil {
			return err
		}
		_, err = io.Copy(tmpFile, resp.Body)
		if err != nil {
			return stop{err}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Couldn't execute http request %s\n", err)
	}
	fmt.Println(tmpFile.Name())
}
