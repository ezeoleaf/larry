# Larry 🐦
[![Go](https://github.com/ezeoleaf/larry/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ezeoleaf/larry/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/ezeoleaf/larry/badge.svg?branch=main)](https://coveralls.io/github/ezeoleaf/larry?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/ezeoleaf/larry)](https://goreportcard.com/report/github.com/ezeoleaf/larry)
[![MIT License](https://img.shields.io/github/license/ezeoleaf/larry?style=flat-square)](https://github.com/ezeoleaf/larry/blob/main/LICENSE)

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
- TWITTER_CONSUMER_KEY
- TWITTER_CONSUMER_SECRET
- TWITTER_ACCESS_TOKEN
- TWITTER_ACCESS_SECRET
```

For generating Github access token you can follow this [guide](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token)

For getting Twitter keys and secrets you can follow this [guide](https://developer.twitter.com/en/docs/twitter-api/getting-started/guide)

#### Note: You will to generate both consumer and access pair of keys/tokens and secrets

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
   Larry - Twitter bot that publishes random information from providers

USAGE:
   larry [global options] command [command options] [arguments...]

AUTHORS:
   @ezeoleaf <ezeoleaf@gmail.com>
   @beesaferoot <hikenike6@gmail.com>

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
   --help, -h                      show help (default: false)
```

For running the bot, the command will depend on whatever you want to tweet, but, for tweeting about React repositories every 30 minutes, you could use

&nbsp;&nbsp;`larry --topic react --time 30 --safe-mode`

For running the bot for Rust tweets every 15 minutes

&nbsp;&nbsp;`larry --lang rust --time 15`


## Have questions? Need help with the bot?

If you're having issues with or have questions about the bot, [file an issue](https://github.com/ezeoleaf/larry/issues) in this repository so anyone can get back to you.

Or feel free to contact me <ezeoleaf@gmail.com> :)

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/H2H47X7QW)
