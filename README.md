# Larry 🐦
[![Go](https://github.com/ezeoleaf/larry/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ezeoleaf/larry/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/ezeoleaf/larry/badge.svg?branch=main)](https://coveralls.io/github/ezeoleaf/larry?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/ezeoleaf/larry)](https://goreportcard.com/report/github.com/ezeoleaf/larry)
[![MIT License](https://img.shields.io/github/license/ezeoleaf/larry?style=flat-square)](https://github.com/ezeoleaf/larry/blob/main/LICENSE)
[![Contribute with Gitpod](https://img.shields.io/badge/Contribute%20with-Gitpod-908a85?logo=gitpod)](https://gitpod.io/#https://github.com/ezeoleaf/larry)

Larry is a Golang cli bot that tweets random Github repositories.

## Disclaimer

I hold no liability for what you do with this bot or what happens to you by using this bot. Abusing this bot *can* get you banned from Twitter, so make sure to read up on [proper usage](https://support.twitter.com/articles/76915-automation-rules-and-best-practices) of the Twitter API.

## Running bots

- [GolangRepos](https://twitter.com/GolangRepos): Tweets repositories from Github that contain the "golang" topic
- [RustRepos](https://twitter.com/RustRepos): Tweets repositories from Github that contain the "rust" topic
- [MLRepositories](https://twitter.com/MLRepositories): Tweets repositories from Github that contain the "machine-learning" topic
- [CryptoRepos](https://twitter.com/CryptoRepos): Tweets repositories from Github that contain the "crypto" topic

## Installation

You can install Larry by cloning the repo and using `go install`

```bash
git clone https://github.com/ezeoleaf/larry.git
cd larry/cmd/larry
go install
```

You can also use make for building the project and generating an executable:
```bash
git clone https://github.com/ezeoleaf/larry.git
cd larry
make build
```

Or you can just run it on the go
```bash
git clone https://github.com/ezeoleaf/larry.git
cd larry
go run . [options]
```

## Usage

### Configuring the bot

Before running the bot, you must first set it up so it can connect to Github and Twitter API.

To do this, you will need to setup the following environment variables:
```
- GITHUB_ACCESS_TOKEN
- TWITTER_CLIENT_ID
- TWITTER_CLIENT_SECRET

If you want the content to be publish in a README file on a repo, you also need these variables
- GITHUB_PUBLISH_REPO_OWNER (Your Github username)
- GITHUB_PUBLISH_REPO_NAME (The name of the repo where your README is. It has to be public)
- GITHUB_PUBLISH_REPO_FILE (By default is README)
```

For generating Github access token you can follow this [guide](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token)

For getting Twitter client id and secret you can follow this [guide](https://developer.twitter.com/en/docs/apps/overview)

#### Note: You will have to generate both consumer and access pair of keys/tokens and secrets.

### Providers and Publishers

For information on publishers and providers click [here](PublishersAndProviders.md)
#### Providers (where the information is coming from)

| Name         | Key    | Environment Variables |
|--------------|:------:|-----------------------|
| Github       | github |   GITHUB_ACCESS_TOKEN |

_NOTE: The key is used in the --provider or --pr option_

#### Publishers (where the information is going to be posted)

| Name         | Key     | Environment Variables | Observation |
|--------------|:-------:|----------------------| -------------|
| Twitter      | twitter | TWITTER_CLIENT_ID<br>TWITTER_CLIENT_SECRET | |
| Github       | github  | GITHUB_PUBLISH_REPO_OWNER<br>GITHUB_PUBLISH_REPO_NAME<br>GITHUB_PUBLISH_REPO_FILE | For now it is only going to be posted in the README file and the repository must be **public** |

_NOTE: The key is used in the --publisher or --pub option_

### Running the bot

To run the bot, you have two ways.

If you have installed the bot, you can run it using
  `larry [options]`

If you want to run it without installing it globally you can use
  `go run . [options]`

Example:

`larry -h`

As a response you will see the entire options available

```
NAME:
   Larry - Bot that publishes information from providers to different publishers

USAGE:
   larry [global options] command [command options] [arguments...]

AUTHORS:
   @ezeoleaf <ezeoleaf@gmail.com>
   @beesaferoot <hikenike6@gmail.com>
   @shubhcoder
   @kannav02
   @siddhant-k-code <siddhantkhare2694@gmail.com>
   @savagedev

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --topic value, -t value         topic for searching repos
   --lang value, -l value          language for searching repos
   --time value, -x value          periodicity of tweet in minutes (default: 15)
   --cache value, -r value         size of cache for no repeating repositories (default: 50)
   --hashtag value, --ht value     list of comma separated hashtags
   --tweet-language, --tl          bool for allowing twetting the language of the repo (default: false)
   --safe-mode, --sf               bool for safe mode. If safe mode is enabled, no repository is published (default: false)
   --provider value, --pr value    provider where publishable content comes from (default: "github")
   --publisher value, --pub value  list of comma separared publishers (default: "twitter")
   --blacklist value, --bl value   optional file containing blacklisted repository Ids
   --help, -h                      show help (default: false)
```

For running the bot, the command will depend on whatever you want to tweet, but, for tweeting about React repositories every 30 minutes, you could use

&nbsp;&nbsp;`larry --topic react --time 30 --safe-mode`

For running the bot for Rust tweets every 15 minutes

&nbsp;&nbsp;`larry --lang rust --time 15`

For running the bot for Golang every 15 minutes and specifying a blacklist file named blacklist.txt

&nbsp;&nbsp;`larry --topic golang --time 15 --blacklist ./blacklist.txt`


## Blacklist File

The optional blacklist file consists of numeric GitHub repository IDs to exclude from the publishing process. These IDs can be found on the GitHub repository page source in the meta tag `octolytics-dimension-repository_id`.

An example blacklist file containing GitHub repository IDs. The file can contain comments using the # character, everything on the line after this character is ignored.

```
# Blacklisted repositories
123 # description of the respository
456
```


## Have questions? Need help with the bot?

If you're having issues with or have questions about the bot, [file an issue](https://github.com/ezeoleaf/larry/issues) in this repository so anyone can get back to you.

Or feel free to contact me <ezeoleaf@gmail.com> :)

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=ezeoleaf/larry&type=Date)](https://star-history.com/#ezeoleaf/larry)

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/H2H47X7QW)
