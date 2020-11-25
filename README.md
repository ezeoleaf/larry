# GobotTweet

A Golang cli bot that tweets random Github repositories.

You can check a fuctional version working [here](https://twitter.com/GolangRepos)

## Disclaimer

I hold no liability for what you do with this bot or what happens to you by using this bot. Abusing this bot *can* get you banned from Twitter, so make sure to read up on [proper usage](https://support.twitter.com/articles/76915-automation-rules-and-best-practices) of the Twitter API.

## Installation

You can install the GobotTweet bot by cloning the repo and using `go install`

```bash
git clone https://github.com/ezeoleaf/GobotTweet.git
cd GobotTweet
go install
```

Or you can run it on the go
```bash
git clone https://github.com/ezeoleaf/GobotTweet.git
cd GobotTweet
go run . [options]
```

## Usage

### Configuring the bot

Before running the bot, you must first set it up so it can connect to Github and Twitter API. Create a config.json (or rename the config.example.json) file and fill in the following information:
```json
{
    "github_access_token": "xxxxx",
    "twitter_consumer_key": "xxxxx",
    "twitter_consumer_secret": "xxxxx",
    "twitter_access_token": "xxxxx",
    "twitter_access_secret": "xxxxx"
}
```

For generating Github access token you can follow this [guide](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token)

For getting Twitter keys and secrets you can follow this [guide](https://developer.twitter.com/en/docs/twitter-api/getting-started/guide)

### Running the bot

To run the bot, you have two ways.

If you have installed the bot, you can run it using
  `GobotTweet [options]`

If you want to run it without installing it globally you can use
  `go run . [options]`

Example:

`GobotTweet -h`

As a response you will see the entire options available

```
NAME:
   GobotTweet - Twitter bot that tweets random repositories

USAGE:
   GobotTweet [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --topic value, -t value   topic for searching repos
   --lang value, -l value    language for searching repos
   --config value, -c value  path to config file (default: "./config.json")
   --time value, -x value    periodicity of tweet in minutes (default: 15)
   --help, -h                show help (default: false)
```

For running the bot, the command will depend on whatever you want to tweet, but, for tweeting about React repositories every 30 minutes, you could use

&nbsp;&nbsp;`GobotTweet --topic react --config "path_to_react_config.json" --time 30`

For running the bot for Rust tweets every 15 minutes

&nbsp;&nbsp;`GobotTweet --lang rust --config "path_to_rust_config.json" --time 15`


## Have questions? Need help with the bot?

If you're having issues with or have questions about the bot, [file an issue](https://github.com/ezeoleaf/GobotTweet/issues) in this repository so anyone can get back to you.

Or feel free to contact me <ezeoleaf@gmail.com> :)
