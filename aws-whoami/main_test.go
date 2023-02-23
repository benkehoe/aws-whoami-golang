package main

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func NewGetCallerIdentityOutput(account string, arn string, userId string) sts.GetCallerIdentityOutput {
	var getCallerIdentityOutput sts.GetCallerIdentityOutput
	getCallerIdentityOutput.Account = &account
	getCallerIdentityOutput.Arn = &arn
	getCallerIdentityOutput.UserId = &userId
	return getCallerIdentityOutput
}

func newWhoami(account string, arn string, userId string, accountAliases []string, region string) Whoami {
	whoami := Whoami{}
	whoami.AccountAliases = accountAliases
	whoami.Region = region

	getCallerIdentityOutput := NewGetCallerIdentityOutput(account, arn, userId)

	populateWhoamiFromGetCallerIdentityOutput(&whoami, getCallerIdentityOutput)

	return whoami
}

func TestBasic(t *testing.T) {
	account := "123456789012"
	arn := "arn:aws:iam::123456789012:assumed-role/ben/my-session"
	userId := "AIDAJQABLZS4A3QDU576Q"

	getCallerIdentityOutput := NewGetCallerIdentityOutput(account, arn, userId)

	whoami := Whoami{}

	populateWhoamiFromGetCallerIdentityOutput(&whoami, getCallerIdentityOutput)

	if whoami.Account != account {
		t.Fatalf("Account set incorrectly to '%v'", whoami.Account)
	}

	if whoami.Arn != arn {
		t.Fatalf("Arn set incorrectly to '%v'", whoami.Arn)
	}

	if whoami.UserId != userId {
		t.Fatalf("UserId set incorrectly to '%v'", whoami.UserId)
	}
}

func TestRootArn(t *testing.T) {
	account := "123456789012"
	arn := "arn:aws:iam::123456789012:root"
	userId := "AIDAJQABLZS4A3QDU576Q"

	getCallerIdentityOutput := NewGetCallerIdentityOutput(account, arn, userId)

	whoami := Whoami{}

	populateWhoamiFromGetCallerIdentityOutput(&whoami, getCallerIdentityOutput)

	if whoami.Type != "root" {
		t.Fatalf("Type is %v (should be root)", whoami.Type)
	}

	if whoami.Name != "root" {
		t.Fatalf("Name is %v (should be root)", whoami.Name)
	}

	if whoami.RoleSessionName != nil {
		t.Fatalf("RoleSessionName is set (to %v)", *whoami.RoleSessionName)
	}

	if whoami.SSOPermissionSet != nil {
		t.Fatalf("SSOPermissionSet is set (to %v)", *whoami.SSOPermissionSet)
	}
}

func TestUserArn(t *testing.T) {
	account := "123456789012"
	arn := "arn:aws:iam::123456789012:user/some/path/ben"
	userId := "AIDAJQABLZS4A3QDU576Q"

	getCallerIdentityOutput := NewGetCallerIdentityOutput(account, arn, userId)

	whoami := Whoami{}

	populateWhoamiFromGetCallerIdentityOutput(&whoami, getCallerIdentityOutput)

	if whoami.Type != "user" {
		t.Fatalf("Type is %v (should be user)", whoami.Type)
	}

	if whoami.Name != "ben" {
		t.Fatalf("Name is %v (should be ben)", whoami.Name)
	}

	if whoami.RoleSessionName != nil {
		t.Fatalf("RoleSessionName is set (to %v)", *whoami.RoleSessionName)
	}

	if whoami.SSOPermissionSet != nil {
		t.Fatalf("SSOPermissionSet is set (to %v)", *whoami.SSOPermissionSet)
	}
}

func TestAssumedRoleArn(t *testing.T) {
	account := "123456789012"
	arn := "arn:aws:iam::123456789012:assumed-role/ben/my-session"
	userId := "AROAJQABLZS4A3QDU576Q"

	getCallerIdentityOutput := NewGetCallerIdentityOutput(account, arn, userId)

	whoami := Whoami{}

	populateWhoamiFromGetCallerIdentityOutput(&whoami, getCallerIdentityOutput)

	if whoami.Type != "assumed-role" {
		t.Fatalf("Type is %v (should be assumed-role)", whoami.Type)
	}

	if whoami.Name != "ben" {
		t.Fatalf("Name is %v (should be ben)", whoami.Name)
	}

	if whoami.RoleSessionName == nil {
		t.Fatalf("RoleSessionName is not set")
	}

	if *whoami.RoleSessionName != "my-session" {
		t.Fatalf("RoleSessionName is %v (should be my-session)", whoami.RoleSessionName)
	}

	if whoami.SSOPermissionSet != nil {
		t.Fatalf("SSOPermissionSet is set (to %v)", *whoami.SSOPermissionSet)
	}
}

func TestSSOArn(t *testing.T) {
	account := "123456789012"
	arn := "arn:aws:iam::123456789012:assumed-role/AWSReservedSSO_SsoRole_abc123/ben"
	userId := "AROAJQABLZS4A3QDU576Q"

	getCallerIdentityOutput := NewGetCallerIdentityOutput(account, arn, userId)

	whoami := Whoami{}

	populateWhoamiFromGetCallerIdentityOutput(&whoami, getCallerIdentityOutput)

	if whoami.Type != "assumed-role" {
		t.Fatalf("Type is %v (should be assumed-role)", whoami.Type)
	}

	if whoami.Name != "AWSReservedSSO_SsoRole_abc123" {
		t.Fatalf("Name is %v (should be AWSReservedSSO_SsoRole_abc123)", whoami.Name)
	}

	if whoami.RoleSessionName == nil {
		t.Fatalf("RoleSessionName is not set")
	}

	if *whoami.RoleSessionName != "ben" {
		t.Fatalf("RoleSessionName is %v (should be ben)", whoami.RoleSessionName)
	}

	if whoami.SSOPermissionSet == nil {
		t.Fatalf("SSOPermissionSet is not set")
	}

	if *whoami.SSOPermissionSet != "SsoRole" {
		t.Fatalf("SSOPermissionSet is %v (should be SsoRole)", *whoami.SSOPermissionSet)
	}
}

func TestFederatedUser(t *testing.T) {
	account := "123456789012"
	arn := "arn:aws:iam::123456789012:federated-user/ben"
	userId := "AIDAJQABLZS4A3QDU576Q"

	getCallerIdentityOutput := NewGetCallerIdentityOutput(account, arn, userId)

	whoami := Whoami{}

	populateWhoamiFromGetCallerIdentityOutput(&whoami, getCallerIdentityOutput)

	if whoami.Type != "federated-user" {
		t.Fatalf("Type is %v (should be federated-user)", whoami.Type)
	}

	if whoami.Name != "ben" {
		t.Fatalf("Name is %v (should be ben)", whoami.Name)
	}

	if whoami.RoleSessionName != nil {
		t.Fatalf("RoleSessionName is set (to %v)", *whoami.RoleSessionName)
	}

	if whoami.SSOPermissionSet != nil {
		t.Fatalf("SSOPermissionSet is set (to %v)", *whoami.SSOPermissionSet)
	}
}

func TestDisableAccount(t *testing.T) {
	whoami := newWhoami(
		"123456789012", "arn:aws:iam::123456789012:assumed-role/ben/my-session", "AROAJQABLZS4A3QDU576Q", nil, "us-east-1")

	params := WhoamiParams{false, nil}

	if params.GetDisableAccountAlias(whoami) {
		t.Fatalf("Disabled when it shouldn't be")
	}

	params = WhoamiParams{true, nil}

	if !params.GetDisableAccountAlias(whoami) {
		t.Fatalf("Not disabled when it should be")
	}

	params = WhoamiParams{true, []string{"4444", "1234"}}

	// * The beginning or end of the account number
	// * The principal name or ARN
	// * The role session name
	if !params.GetDisableAccountAlias(whoami) {
		t.Fatalf("Should have been disabled by matching account prefix")
	}

	params = WhoamiParams{true, []string{"4444", "9012"}}
	if !params.GetDisableAccountAlias(whoami) {
		t.Fatalf("Should have been disabled by matching account suffix")
	}

	params = WhoamiParams{true, []string{"foo", "ben", "bar"}}
	if !params.GetDisableAccountAlias(whoami) {
		t.Fatalf("Should have been disabled by matching role name")
	}

	params = WhoamiParams{true, []string{"foo", "my-session", "bar"}}
	if !params.GetDisableAccountAlias(whoami) {
		t.Fatalf("Should have been disabled by matching role session name")
	}

	params = WhoamiParams{true, []string{"foo", "arn:aws:iam::123456789012:assumed-role/ben/my-session", "bar"}}
	if !params.GetDisableAccountAlias(whoami) {
		t.Fatalf("Should have been disabled by matching role ARN")
	}

	whoami = newWhoami(
		"123456789012", "arn:aws:iam::123456789012:assumed-role/AWSReservedSSO_SsoRole_abc123/benjamin", "AROAJQABLZS4A3QDU576Q", []string{"account-alias"}, "us-east-1")

	params = WhoamiParams{true, []string{"foo", "SsoRole", "bar"}}
	if !params.GetDisableAccountAlias(whoami) {
		t.Fatalf("Should have been disabled by matching SSO permission set")
	}

	params = WhoamiParams{true, []string{"foo", "benjamin", "bar"}}
	if !params.GetDisableAccountAlias(whoami) {
		t.Fatalf("Should have been disabled by matching role session name")
	}
}
