package print

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/arush-sal/doddl/pkg/getter"
)

func Printer(ds []getter.Droplet) {
	const tformat = "%d\t|%s\t|%s\t|%v\t|%v\t\n"
	const format = "%s\t%s\t%s\t%s\t|%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 4, 1, ' ', 0)

	fmt.Fprintf(tw, format, "ID", "Name", "Status", "Tags", "IP Address")
	fmt.Fprintf(tw, format, "--", "----", "------", "----", "----------")
	for _, v := range ds {
		fmt.Fprintf(tw, tformat, v.ID, v.Name, v.Status, v.Tags, v.IPv4)
	}
	tw.Flush()
}

func JSONPrinter(ds []getter.Droplet) {
	b, err := json.Marshal(ds)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	os.Stdout.Write(b)
}
