package share

import (
	"fmt"
	"os"
	"strconv"

	command "github.com/kosmos.io/netdoctor/pkg/command/share/remote-command"
	"github.com/kosmos.io/netdoctor/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"k8s.io/klog/v2"
)

type PrintCheckData struct {
	command.Result
	SrcNodeName string `json:"srcNodeName"`
	DstNodeName string `json:"dstNodeName"`
	TargetIP    string `json:"targetIP"`
}

func PrintResult(resultData []*PrintCheckData) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"S/N", "SRC_NODE_NAME", "DST_NODE_NAME", "TARGETP", "RESULT"})

	tableException := tablewriter.NewWriter(os.Stdout)
	tableException.SetHeader([]string{"S/N", "SRC_NODE_NAME", "DST_NODE_NAME", "TARGET", "RESULT", "LOG"})

	tableFailed := tablewriter.NewWriter(os.Stdout)
	tableFailed.SetHeader([]string{"S/N", "SRC_NODE_NAME", "DST_NODE_NAME", "TARGET", "RESULT", "LOG"})

	resumeData := []*PrintCheckData{}

	for index, r := range resultData {
		// klog.Infof(fmt.Sprintf("%s %s %v", r.SrcNodeName, r.DstNodeName, r.IsSucceed))
		row := []string{strconv.Itoa(index + 1), r.SrcNodeName, r.DstNodeName, r.TargetIP, command.PrintStatus(r.Status), r.ResultStr}
		if r.Status == command.CommandFailed {
			resumeData = append(resumeData, r)
			tableFailed.Rich(row, []tablewriter.Colors{
				{},
				{tablewriter.Bold, tablewriter.FgHiRedColor},
				{tablewriter.Bold, tablewriter.FgHiRedColor},
				{tablewriter.Bold, tablewriter.FgHiRedColor},
				{tablewriter.Bold, tablewriter.FgHiRedColor},
			})
		} else if r.Status == command.ExecError {
			resumeData = append(resumeData, r)
			tableException.Rich(row, []tablewriter.Colors{
				{},
				{tablewriter.Bold, tablewriter.FgCyanColor},
				{tablewriter.Bold, tablewriter.FgCyanColor},
				{tablewriter.Bold, tablewriter.FgCyanColor},
				{tablewriter.Bold, tablewriter.FgCyanColor},
			})
		} else {
			// resumeData = append(resumeData, r)
			table.Rich(row[:len(row)-1], []tablewriter.Colors{
				{},
				{tablewriter.Bold, tablewriter.FgGreenColor},
				{tablewriter.Bold, tablewriter.FgGreenColor},
				{tablewriter.Bold, tablewriter.FgGreenColor},
				{tablewriter.Bold, tablewriter.FgGreenColor},
			})
		}
	}
	if table.NumLines() > 0 {
		fmt.Println("")
		table.Render()
	}
	if tableFailed.NumLines() > 0 {
		fmt.Println("")
		tableFailed.Render()
	}
	if tableException.NumLines() > 0 {
		fmt.Println("")
		tableException.Render()
	}
	err := utils.WriteResume(resumeData)
	if err != nil {
		klog.Error("write resumeData error")
	}
}
