package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	cli "github.com/spf13/cobra"

	"github.com/solaegis/boilr/pkg/boilr"
	"github.com/solaegis/boilr/pkg/template"
	"github.com/solaegis/boilr/pkg/util/exit"
	"github.com/solaegis/boilr/pkg/util/osutil"
	"github.com/solaegis/boilr/pkg/util/templateinput"
	"github.com/solaegis/boilr/pkg/util/validate"
)

// TemplateInRegistry checks whether the given name exists in the template registry.
func TemplateInRegistry(name string) (bool, error) {
	names, err := ListTemplates()
	if err != nil {
		return false, err
	}

	_, ok := names[name]
	return ok, nil
}

// TODO add --use-cache flag to execute a template from previous answers to prompts
// Use contains the cli-command for using templates located in the local template registry.
var Use = &cli.Command{
	Use:   "use <template-tag> <target-dir>",
	Short: "Execute a project template in the given directory",
	Run: func(cmd *cli.Command, args []string) {
		MustValidateArgs(args, []validate.Argument{
			{"template-tag", validate.UnixPath},
			{"target-dir", validate.UnixPath},
		})

		MustValidateTemplateDir()

		tmplName := args[0]
		targetDir, err := filepath.Abs(args[1])
		if err != nil {
			exit.Fatal(fmt.Errorf("use: %s", err))
		}

		templateFound, err := TemplateInRegistry(tmplName)
		if err != nil {
			exit.Fatal(fmt.Errorf("use: %s", err))
		}

		if !templateFound {
			exit.Fatal(fmt.Errorf("Template %q couldn't be found in the template registry", tmplName))
		}

		tmplPath, err := boilr.TemplatePath(tmplName)
		if err != nil {
			exit.Fatal(fmt.Errorf("use: %s", err))
		}

		tmpl, err := template.Get(tmplPath)
		if err != nil {
			exit.Fatal(fmt.Errorf("use: %s", err))
		}

		contextFile := GetStringFlag(cmd, "use-file")
		if contextFile != "" {

			err := tmpl.CachaedValuesFromJson(contextFile)
			if err != nil {
				exit.Fatal(fmt.Errorf("error reading values value frrom %s", contextFile))
			}
			tmpl.UseDefaultValues()
		}

		if shouldUseDefaults := GetBoolFlag(cmd, "use-defaults"); shouldUseDefaults {
			tmpl.UseDefaultValues()
		}

		executeTemplate := func() error {
			parentDir := filepath.Dir(targetDir)

			exists, err := osutil.DirExists(parentDir)
			if err != nil {
				return err
			}

			if !exists {
				return fmt.Errorf("use: parent directory %q doesn't exist", parentDir)
			}

			tmpDir, err := ioutil.TempDir("", "boilr-use-template")
			if err != nil {
				return err
			}
			defer os.RemoveAll(tmpDir)

			if err := tmpl.Execute(tmpDir); err != nil {
				return err
			}

			// Complete the template execution transaction by copying the temporary dir to the target directory.
			return osutil.CopyRecursively(tmpDir, targetDir)
		}

		if err := executeTemplate(); err != nil {
			exit.Fatal(fmt.Errorf("use: %v", err))
		}

		// store promted inputs in a json file
		jsonFile := GetStringFlag(cmd, "json-file")
		if jsonFile != "" {
			file, err := json.MarshalIndent(templateinput.UserInput, "", "  ")
			if err != nil {
				exit.Fatal(fmt.Errorf("use: %v", err))
			}
			if err := ioutil.WriteFile(jsonFile, file, 0644); err != nil {
				exit.Fatal(fmt.Errorf("use: %v", err))
			}
		}

		if contextFile != "" {
			f, err := os.Open(contextFile)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("could not find ", contextFile, " Please fix the error and run the boilr again.")
					exit.Fatal(fmt.Errorf("use: %v", err))
				}
				fmt.Println("could not read ", contextFile, " Please fix the error and run the boilr again.")
				exit.Fatal(fmt.Errorf("use: %v", err))
			}
			defer f.Close()
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				fmt.Println("could not read ", contextFile, " Please fix the error and run the boilr again.")
				exit.Fatal(fmt.Errorf("use: %v", err))
			}
			var contextFileJson map[string]interface{}
			if err := json.Unmarshal(buf, &contextFileJson); err != nil {
				fmt.Println("could not read ", contextFile, " as a Joson. Please fix the error and run the boilr again.")
				exit.Fatal(fmt.Errorf("use: %v", err))
			}
			returnError := false
			for k, _ := range templateinput.UsedKeys {
				if _, ok := contextFileJson[k]; !ok {
					returnError = true
					fmt.Printf(`
********************************************************************************************************************
boilr used project.json value for key %s. Please define the key and value in %s 
********************************************************************************************************************`,
						k, contextFile)
				}
			}
			if returnError {
				fmt.Println()
				exit.Fatal(fmt.Errorf("Missing values in %s, please review the file", contextFile))
			}
		}

		exit.OK("Successfully executed the project template %v in %v", tmplName, targetDir)
	},
}
