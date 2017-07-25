package spec

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/Songmu/retry"
)

func isEC2ForLinux() bool {
	data, err := ioutil.ReadFile("/sys/hypervisor/uuid")
	if err != nil {
		// Probably not EC2.
		return false
	}
	// Probably not EC2.
	if !strings.HasPrefix(string(data), "ec2") {
		return false
	}
	res := false
	cl := httpCli()
	err = retry.Retry(3, 2*time.Second, func() error {
		// '/ami-id` is probably an AWS specific URL
		resp, err := cl.Get(ec2BaseURL.String() + "/ami-id")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		res = resp.StatusCode == 200
		return nil
	})

	if err == nil {
		return res
	}
	return false
}
