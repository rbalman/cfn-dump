package cmd

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:     "dump",
	Short:   "Dump the cfn stacks dependencies",
	Aliases: []string{"d"},
	Long:    `Dump the cfn stacks export and import dependencies`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return dump(pattern)
	},
}

type CFNManager struct {
}

type Stack struct {
	Exports map[string][]string
}

func dump(stackPrefix string) error {
	cm := CFNManager{}
	exports, err := cm.ListExports()
	if err != nil {
		return err
	}

	dependencies := map[string]Stack{}
	for _, export := range exports {
		var exportingStack string
		stackId := *export.ExportingStackId
		exportName := *export.Name

		re := regexp.MustCompile("\\/.*\\/")
		match := re.FindStringSubmatch(stackId)
		if len(match) <= 0 {
			continue
		}
		exportingStack = strings.Trim(match[0], "/")
		hasStackPrefix := strings.Contains(exportingStack, stackPrefix)
		if !hasStackPrefix {
			log.Printf("[INFO] Skipping... as %s stack does not have %s prefix", exportingStack, stackPrefix)
			continue
		}

		log.Printf("[INFO] Exporting Stack: %s", exportingStack)

		stacks, err := cm.ListImports(exportName)
		if err != nil {
			log.Printf("[WARN] Error while listing imports of export %s. Error: %s", exportName, err)
			continue
		}

		if stack, ok := dependencies[exportingStack]; ok {
			existingExports := stack.Exports
			existingExports[exportName] = stacks
			dependencies[exportingStack] = Stack{
				Exports: existingExports,
			}
		} else {
			dependencies[exportingStack] = Stack{
				Exports: map[string][]string{
					exportName: stacks,
				},
			}
		}
	}

	// if stackPrefix == "" {
	// 	stackPrefix = "global"
	// }

	// For more granular writes, open a file for writing.
	f, err := os.Create("./cfn-dependency.json")
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(dependencies)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (cm CFNManager) ListExports() ([]*cloudformation.Export, error) {
	lei := &cloudformation.ListExportsInput{}
	var exports []*cloudformation.Export

	cfn, err := cm.Session()
	if err != nil {
		return exports, err
	}

	for {
		res, err := cfn.ListExports(lei)
		if err != nil {
			return exports, err
		}

		exports = append(exports, res.Exports...)

		if res.NextToken == nil {
			break
		}
		lei.NextToken = res.NextToken
	}

	return exports, nil
}

func (cm CFNManager) ListImports(name string) ([]string, error) {
	var stacks []string
	cfn, err := cm.Session()
	if err != nil {
		return stacks, err
	}

	lii := &cloudformation.ListImportsInput{ExportName: &name}
	res, err := cfn.ListImports(lii)

	if err != nil {
		if strings.Contains(err.Error(), "not imported by any stack") {
			return stacks, nil
		} else {
			return stacks, err
		}
	}

	for _, i := range res.Imports {
		stacks = append(stacks, *i)
	}

	return stacks, nil
}

func (cm CFNManager) Session() (*cloudformation.CloudFormation, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
		SharedConfigState: session.SharedConfigEnable,
		Profile:           os.Getenv("AWS_PROFILE"),
	}))

	return cloudformation.New(sess), nil
}
