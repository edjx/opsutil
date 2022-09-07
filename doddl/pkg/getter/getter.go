package getter

import (
	ctx "context"

	"github.com/digitalocean/godo"
	"golang.org/x/exp/slices"
)

type DOClient struct {
	c *godo.Client
}

type Droplet struct {
	ID     int      `json:"id,float64,omitempty"`
	Name   string   `json:"name,omitempty"`
	Status string   `json:"status,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	IPv4   string   `json:"ip,omitempty"`
}

// Get a new DO Client
func Client(t string) *DOClient {
	cl := new(DOClient)
	cl.c = godo.NewFromToken(t)
	return cl
}

func (do *DOClient) dropletList() []godo.Droplet {
	// create a list to hold our droplets
	list := []godo.Droplet{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}

	for {
		droplets, resp, err := do.c.Droplets.List(ctx.TODO(), opt)
		if err != nil {
			panic(err)
		}

		// append the current page's droplets to our list
		list = append(list, droplets...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			panic(err)
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list
}

// List all the droplets available for the given DO account
func (do *DOClient) DropletListAll() []Droplet {
	dl := do.dropletList()
	list := []Droplet{}

	for _, v := range dl {
		d := Droplet{
			ID:     v.ID,
			Name:   v.Name,
			Status: v.Status,
			Tags:   v.Tags,
		}

		for _, network := range v.Networks.V4 {
			if network.Type == "public" {
				d.IPv4 = network.IPAddress
			}
		}

		list = append(list, d)
	}

	return list
}

// List all the Droplets with their respective tags
func (do *DOClient) DropletTags() []Droplet {
	dl := do.dropletList()
	list := []Droplet{}

	for _, v := range dl {
		if len(v.Tags) != 0 {
			d := Droplet{
				ID:     v.ID,
				Name:   v.Name,
				Status: v.Status,
				Tags:   v.Tags,
			}

			for _, network := range v.Networks.V4 {
				if network.Type == "public" {
					d.IPv4 = network.IPAddress
				}
			}

			list = append(list, d)
		}
	}

	return list
}

// List all the Droplets that does not have any tag
func (do *DOClient) DropletWithoutAnyTag() []Droplet {
	dl := do.dropletList()
	list := []Droplet{}

	for _, v := range dl {
		if len(v.Tags) == 0 {
			d := Droplet{
				ID:     v.ID,
				Name:   v.Name,
				Status: v.Status,
				Tags:   v.Tags,
			}

			for _, network := range v.Networks.V4 {
				if network.Type == "public" {
					d.IPv4 = network.IPAddress
				}
			}

			list = append(list, d)
		}
	}

	return list
}

// List all the Droplets that has a specific tag
func (do *DOClient) DropletWithSpecificTag(t string) []Droplet {
	dl := do.dropletList()
	list := []Droplet{}

	for _, v := range dl {
		for _, tag := range v.Tags {
			if tag == t {
				d := Droplet{
					ID:     v.ID,
					Name:   v.Name,
					Status: v.Status,
					Tags:   v.Tags,
				}

				for _, network := range v.Networks.V4 {
					if network.Type == "public" {
						d.IPv4 = network.IPAddress
					}
				}

				list = append(list, d)
			}
		}
	}

	return list
}

// List all the Droplets that are in stopped state
func (do *DOClient) StoppedDroplets(ds []Droplet) []Droplet {
	list := []Droplet{}

	for _, v := range ds {
		if v.Status == "off" {
			d := Droplet{
				ID:     v.ID,
				Name:   v.Name,
				Status: v.Status,
				Tags:   v.Tags,
				IPv4:   v.IPv4,
			}

			list = append(list, d)
		}
	}

	return list
}

// List all the Droplets that does not have the specific tag
func (do *DOClient) DropletWithoutSpecificTag(tags []string) []Droplet {
	dl := do.dropletList()
	list := []Droplet{}

	for _, v := range dl {
		var notFound = true

		for _, tag := range tags {
			// GENERICS Baby!!
			// Would require go v1.18+
			// slices.Index returns the index of t if found in []T
			// and would return -1 if not found
			if slices.Index(v.Tags, tag) != -1 {
				notFound = false
			}
		}

		if notFound {
			d := Droplet{
				ID:     v.ID,
				Name:   v.Name,
				Status: v.Status,
				Tags:   v.Tags,
			}

			for _, network := range v.Networks.V4 {
				if network.Type == "public" {
					d.IPv4 = network.IPAddress
				}
			}

			list = append(list, d)
		}
	}

	return list
}
