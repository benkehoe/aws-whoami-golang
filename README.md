# aws-whoami
**Find what AWS account and identity you're using**

> :warning: This is the successor to [the python implementation](https://github.com/benkehoe/aws-whoami) as a CLI tool. The other is still useful as a Python library.

You should know about [`aws sts get-caller-identity`](https://docs.aws.amazon.com/cli/latest/reference/sts/get-caller-identity.html), which sensibly returns the identity of the caller.
But even with `--output table`, I find this a bit lacking.
That ARN is a lot to visually parse, it doesn't tell you what region your credentials are configured for, and I am not very good at remembering AWS account numbers. `aws-whoami` makes it better.

```
$ aws-whoami
Account:         123456789012
                 my-account-alias
Region:          us-east-2
AssumedRole:     MyRole
RoleSessionName: ben
UserId:          AROASOMEOPAQUEID:ben
Arn:             arn:aws:sts::123456789012:assumed-role/MyRole/ben
```

Note: if you don't have permissions to [iam:ListAccountAliases](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListAccountAliases.html), your account alias won't appear.
See below for disabling this check if getting a permission denied on this call raises flags in your organization.

## Install

```
go install github.com/benkehoe/aws-whoami-golang/v2/aws-whoami@latest
```

[Or download the latest release for your platform](https://github.com/benkehoe/aws-whoami-golang/releases/latest).


## Options

`aws-whoami` uses [`the AWS Go SDK v2`](https://aws.amazon.com/sdk-for-go/), so it'll pick up your credentials in [the normal ways](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html#config-settings-and-precedence), including with the `--profile` parameter.

If you'd like the output as a JSON object, use the `--json` flag.
See below for field names.

The `--disable-account-alias` flag disables account alias checking (see below).

Use `--version` to output the version.

## Account alias checking

By default, `aws-whoami` calls the [IAM `ListAccountAliases` API](https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListAccountAliases.html) to find the account name, if set.
If you don't have access to this API (the `iam:ListAccountAliases` IAM action), it swallows that error.
In general this is fine, but if it causes trouble (e.g., raising security alerts in your organization), you can disable it.

There are two ways to disable account alias checking.
The first is the `--disable-account-alias` flag.
The second, setting the environment variable `AWS_WHOAMI_DISABLE_ACCOUNT_ALIAS`, allows for persistent and selective control.

To fully disable account alias checking, set `AWS_WHOAMI_DISABLE_ACCOUNT_ALIAS` to `true`.
To selectively disable it, you can also set the value to a comma-separated list where each item will be matched against the following:
* The beginning or end of the account number
* The principal name or ARN
* The role session name
* The SSO role (permission set) name

## JSON output

The JSON object that is printed when using the `--json` flag always includes the following fields:
* `Account`
* `AccountAliases` (NOTE: this is a list)
* `Arn`
* `Type`
* `Name`
* `RoleSessionName`
* `UserId`
* `Region`
* `SSOPermissionSet`

`Type`, `Name`, and `RoleSessionName` (and `SSOPermissionSet`) are split from the ARN for convenience.
`RoleSessionName` is `null` for IAM users.
For the account root, both the `Type` and `Name` are `"root"`.

`SSOPermissionSet` is set if the assumed role name conforms to the format `AWSReservedSSO_{permission-set}_{random-tag}`, otherwise it is `null`.

Note that the `AccountAliases` field is an empty list when account alias checking is disabled, not `null`.
