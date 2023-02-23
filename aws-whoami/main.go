// Copyright 2023 Ben Kehoe
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/aws/smithy-go"
)

var Version string = "2.5"
var DisableAccountAliasEnvVarName = "AWS_WHOAMI_DISABLE_ACCOUNT_ALIAS"

type Whoami struct {
	Account          string
	AccountAliases   []string
	Arn              string
	Type             string
	Name             string
	RoleSessionName  *string
	UserId           string
	Region           string
	SSOPermissionSet *string
}

type WhoamiParams struct {
	DisableAccountAlias       bool
	DisableAccountAliasValues []string
}

func NewWhoamiParams() WhoamiParams {
	var params WhoamiParams
	disableAccountAliasValue := os.Getenv(DisableAccountAliasEnvVarName)
	populateDisableAccountAlias(&params, disableAccountAliasValue)
	return params
}

func populateDisableAccountAlias(params *WhoamiParams, disableAccountAliasValue string) {
	switch strings.ToLower(disableAccountAliasValue) {
	case "":
		fallthrough
	case "0":
		fallthrough
	case "false":
		params.DisableAccountAlias = false
		params.DisableAccountAliasValues = nil
		return
	case "1":
		fallthrough
	case "true":
		params.DisableAccountAlias = true
		params.DisableAccountAliasValues = nil
		return
	default:
		accounts := strings.Split(disableAccountAliasValue, ",")
		if len(accounts) > 0 {
			params.DisableAccountAlias = true
			params.DisableAccountAliasValues = accounts
			return
		} else {
			params.DisableAccountAlias = false
			params.DisableAccountAliasValues = nil
			return
		}
	}
}

func (params WhoamiParams) GetDisableAccountAlias(whoami Whoami) bool {
	if !params.DisableAccountAlias {
		return false
	}
	if params.DisableAccountAliasValues == nil {
		return true
	}
	for _, disabledValue := range params.DisableAccountAliasValues {
		if strings.HasPrefix(whoami.Account, disabledValue) || strings.HasSuffix(whoami.Account, disabledValue) {
			return true
		}
		if whoami.Arn == disabledValue || whoami.Name == disabledValue {
			return true
		}
		if whoami.RoleSessionName != nil && *whoami.RoleSessionName == disabledValue {
			return true
		}
		if whoami.SSOPermissionSet != nil && *whoami.SSOPermissionSet == disabledValue {
			return true
		}
	}
	return false
}

func populateWhoamiFromGetCallerIdentityOutput(whoami *Whoami, getCallerIdentityOutput sts.GetCallerIdentityOutput) error {
	whoami.Account = *getCallerIdentityOutput.Account
	whoami.Arn = *getCallerIdentityOutput.Arn
	whoami.UserId = *getCallerIdentityOutput.UserId

	arnFields := strings.Split(whoami.Arn, ":")

	var arnResourceFields []string
	if arnFields[len(arnFields)-1] == "root" {
		arnResourceFields = []string{"root", "root"}
	} else {
		arnResourceFields = strings.SplitN(arnFields[len(arnFields)-1], "/", 2)
		if len(arnResourceFields) < 2 {
			return fmt.Errorf("arn %v has an unknown format", whoami.Arn)
		}
	}

	whoami.Type = arnResourceFields[0]
	if whoami.Type == "assumed-role" {
		nameFields := strings.SplitN(arnResourceFields[1], "/", 2)
		if len(arnResourceFields) < 2 {
			return fmt.Errorf("arn %v has an unknown format", whoami.Arn)
		}
		whoami.Name = nameFields[0]
		whoami.RoleSessionName = &nameFields[1]
	} else if whoami.Type == "user" {
		nameFields := strings.Split(arnResourceFields[1], "/")
		whoami.Name = nameFields[len(nameFields)-1]
	} else {
		whoami.Name = arnResourceFields[1]
	}

	if whoami.Type == "assumed-role" && strings.HasPrefix(whoami.Name, "AWSReservedSSO") {
		nameFields := strings.Split(whoami.Name, "_")
		if len(nameFields) >= 3 {
			permSetStr := strings.Join(nameFields[1:len(nameFields)-1], "_")
			whoami.SSOPermissionSet = &permSetStr
		}
	}

	return nil
}

func NewWhoami(awsConfig aws.Config, params WhoamiParams) (Whoami, error) {
	stsClient := sts.NewFromConfig(awsConfig)

	getCallerIdentityOutput, err := stsClient.GetCallerIdentity(context.TODO(), nil)

	if err != nil {
		return Whoami{}, err
	}

	var whoami Whoami
	whoami.AccountAliases = make([]string, 0, 1)

	whoami.Region = awsConfig.Region

	err = populateWhoamiFromGetCallerIdentityOutput(&whoami, *getCallerIdentityOutput)

	if err != nil {
		return whoami, err
	}

	if !params.GetDisableAccountAlias(whoami) {
		iam_client := iam.NewFromConfig(awsConfig)

		// pedantry
		paginator := iam.NewListAccountAliasesPaginator(iam_client, nil)

		for paginator.HasMorePages() {
			output, err := paginator.NextPage(context.TODO())
			if err != nil {
				var apiErr smithy.APIError
				if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
					break
				} else {
					return whoami, err
				}
			}
			whoami.AccountAliases = append(whoami.AccountAliases, output.AccountAliases...)
		}
	}

	return whoami, nil
}

type record struct {
	field string
	value string
}

func getTypeNameRecord(whoami Whoami) record {
	if whoami.Type == "root" {
		return record{"Type: ", "root"}
	}
	fields := strings.Split(whoami.Type, "-")
	typeParts := make([]string, 0, 3)
	for _, field := range fields {
		s := strings.ToUpper(field[:1]) + field[1:] // ok because always ASCII
		typeParts = append(typeParts, s)
	}
	typeParts = append(typeParts, ": ")
	return record{strings.Join(typeParts, ""), whoami.Name}
}

func (whoami Whoami) Format() string {
	records := make([]record, 0, 7)
	records = append(records, record{"Account: ", whoami.Account})
	for _, alias := range whoami.AccountAliases {
		records = append(records, record{"", alias})
	}
	records = append(records, record{"Region: ", whoami.Region})
	if whoami.SSOPermissionSet != nil {
		records = append(records, record{"AWS SSO: ", *whoami.SSOPermissionSet})
	} else {
		records = append(records, getTypeNameRecord(whoami))
	}
	if whoami.RoleSessionName != nil {
		records = append(records, record{"RoleSessionName: ", *whoami.RoleSessionName})
	}
	records = append(records, record{"UserId: ", whoami.UserId})
	records = append(records, record{"Arn: ", whoami.Arn})

	var maxLen int = 0
	for _, rec := range records {
		if len(rec.field) > maxLen {
			maxLen = len(rec.field)
		}
	}

	lines := make([]string, 0, 7)
	for _, rec := range records {
		lines = append(lines, rec.field+strings.Repeat(" ", maxLen-len(rec.field))+rec.value)
	}

	return strings.Join(lines, "\n")
}

func main() {
	profile := flag.String("profile", "", "A config profile to use")
	useJson := flag.Bool("json", false, "Output as JSON")
	disableAccountAlias := flag.Bool("disable-account-alias", false, "Disable account alias check")
	showVersion := flag.Bool("version", false, "Display the version")
	flag.Parse()

	if *showVersion {
		fmt.Println(Version)
		return
	}

	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(*profile))
	if err != nil {
		log.Fatal(err) // exits
	}

	whoamiParams := NewWhoamiParams()

	if *disableAccountAlias {
		whoamiParams.DisableAccountAlias = true
		whoamiParams.DisableAccountAliasValues = nil
	}

	Whoami, err := NewWhoami(awsConfig, whoamiParams)

	if err != nil {
		log.Fatal(err) // exits
	}

	if *useJson {
		bytes, err := json.Marshal(Whoami)
		if err != nil {
			log.Fatal(err) // exits
		}
		fmt.Println(string(bytes))
	} else {
		fmt.Println(Whoami.Format())
	}
}
