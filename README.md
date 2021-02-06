![Sqealer](squealer.png)

# Squealer

### Telling tales on you for leaking secrets!

Squealer scans a local git repository for secrets that are being leaked deep within the commit history. 

The built-in configuration has the following checks;

AWS
- access key id
- access secret key
- mws key

Github
- github token

Slack
- slack token
- webhook url

Other
- Asymmetric Private Key

Sometimes we have secrets committed to our projects, generally we can invalidate them and move on. If squealer is telling tales about a secret that you are aware of and has been mitigated, you can use the `exception` rule found in the output to register it as ignored.

## Installation

curl -s "https://raw.githubusercontent.com/owenrumney/squealer/main/scripts/install.sh" | bash


## Usage

Squealer is intended to be run either locally or as part of a CI process. 

### CLI Arguments

#### `--redacted`

The `--redacted` argument is used to redact the secrets in the report... no point exposing them further 

#### `--config-file`

Squealer will use the built in config unless you override it using the `--config-file` option. This config file should be laid out in the following format

```yaml
rules:
- rule: (A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}
  description: Check for AWS Access Key Id
- rule: (?i)aws(.{0,20})?(?-i)['\"][0-9a-zA-Z\/+]{40}['\"]
  description: Check for AWS Secret Access Key
- rule: amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}
  description: Check for AWS MWS Key
- rule: (?i)github[_\-\.]?token[\s:,="\]']+?(?-i)[0-9a-zA-Z]{35,40}
  description: Check for Github Token 
- rule: https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}
  description: Check for Slack webhook
- rule: xox[baprs]-([0-9a-zA-Z]{10,48})?
  description: Check for Slack token
- rule: '-----BEGIN ((EC|PGP|DSA|RSA|OPENSSH) )?PRIVATE KEY( BLOCK)?-----'
  description: Check for Private Asymetric Key
ignore_prefixes:
- vendor
- node_modules
ignore_extensions:
- .zip
- .png
- .jpg
- .pdf
- .xls
- .doc
- .docx
exceptions:
- exception: release/update.go:D2IDetI6aidl58GE6dv5uAaWmXM=
  reason: This is a webhook that we got rid of - can be ignored in this file
```

### Config breakdown

The config file is made up of the `rules`, `ignore_prefixes`, `ignore_extensions` and `exceptions`. 

#### rules

Rules define the regular expression that is used to detect the secret. Requires a description for posterity.

#### ignore_prefixes

Ignore prefixes are folders that you don't want to look ing - generally `vendor` and the like.

#### ignore_extensions

Ignore extensions have the file types that won't be scanned. Binaries are automatically ignored.

#### exceptions

Exceptions are the entries that you've already handled and don't want to be reported any more.

## Example Output

```shell
squealer /home/owen/go/src/github.com/owenrumney/go-github-pr-commenter     
                            
scanning /home/owen/go/src/github.com/owenrumney/go-github-pr-commenter

Process took 0.084771843

Processing:
  commits:      25
  commit files: 291

Transgressions:
  identified:   0
  ignored:      0
  reported:     0

```

## Credits

[Image by Derangedmisfit](https://derangedmisfit.newgrounds.com/)