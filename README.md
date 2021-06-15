# ssmbrowse

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/bnaydenov/ssmbrowse/goreleaser)
[![GitHub license](https://img.shields.io/github/license/bnaydenov/ssmbrowse)](https://github.com/bnaydenov/ssmbrowse/blob/master/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/bnaydenov/ssmbrowse)](https://github.com/bnaydenov/ssmbrowse/issues)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/bnaydenov/ssmbrowse)
![GitHub all releases](https://img.shields.io/github/downloads/bnaydenov/ssmbrowse/total)
![Go Report Card](https://goreportcard.com/badge/github.com/bnaydenov/ssmbrowse)


Simple and elegant cli AWS SSM parameter browser.

<img align="left" src="assets/ssmbrowse-logo-transparent.png" style="float:  unset;">
<br clear="left"/>


## Install on Macos with brew: 
```bash
brew tap bnaydenov/ssmbrowse
brew install ssmbrowse
```

`ssmbrowse` uses AWS Golang SDK to access AWS. The AWS SDK for Go requires credentials (an access key and secret access key) to sign requests to AWS. You can specify your credentials in several different locations, depending on your particular use case. 

The default provider chain looks for credentials in the following order:
1. Environment variables.

2. Shared credentials file.
3. If your application is running on an Amazon EC2 instance, IAM role for Amazon EC2.

The SDK detects and uses the built-in providers automatically, without requiring manual configurations. For example, if you use IAM roles for Amazon EC2 instances, `ssmbrowse` automatically use the instance’s credentials. You don’t need to manually configure credentials in your application.

### Shared Credentials File
A credential file is a plaintext file that contains your access keys. The file must be on the same machine on which you’re running your application. The file must be named credentials and located in the .aws/ folder in your home directory. The home directory can vary by operating system. In Windows, you can refer to your home directory by using the environment variable %UserProfile%. In Unix-like systems, you can use the environment variable $HOME or ~ (tilde).

If you already use this file for other SDKs and tools (like the AWS CLI), you don’t need to change anything to use the files in this SDK. If you use different credentials for different tools or applications, you can use profiles to configure multiple access keys in the same configuration file.

 [More information how to configure credential check AWS docs here](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)

 ## Demo 

 ### Find all SSM params:
 To find all params simply type `/` as param prefix to search for. 
 
 
 By default is display first 30 results and if there is more params available will display them once you are move to last one with help of keyboard arrows.

### Find only specific SSM params containing prefix:
To find specific params containing prefix simply type it into prefix to search for. 



<img align="left" src="assets/demo-monokai1.gif" style="float:  unset;">
<br clear="left"/>