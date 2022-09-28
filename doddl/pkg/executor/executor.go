package executor

import (
	"fmt"
	"os"

	"github.com/edjx/opsutil/doddl/pkg/getter"
)

var whiteList = []string{
	"development-nodes",
	"load-nodes",
	"performance-nodes",
	"staging-nodes",
	"production-nodes",
	"k8s",
	"Odoo",
	"west-nodes",
	"lab-nodes",
	"termination-protected",
}

func GetDOClient(token string) *getter.DOClient {
	if token == "" {
		token = os.Getenv("DO_TOKEN")
	}

	if token == "" {
		fmt.Println("Required token for DigitalOcean is missing")
		os.Exit(1)
	}

	return getter.Client(token)
}

func Run(c *getter.DOClient, list, tag string) []getter.Droplet {

	switch list {
	case "all-tagged":
		return c.DropletTags()
	case "no-tag":
		return c.DropletWithoutAnyTag()
	case "tag":
		return c.DropletWithSpecificTag(tag)
	case "no-whitelist":
		return c.DropletWithoutSpecificTag(whiteList)
	default:
		return c.DropletListAll()
	}

}

func RunStopped(c *getter.DOClient, ds []getter.Droplet) []getter.Droplet {
	return c.StoppedDroplets(ds)
}
