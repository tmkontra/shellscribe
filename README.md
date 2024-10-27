## shellscribe

### Installation

requirements:

- zsh
- git
- go1.23

```sh
git clone git@github.com:tmkontra/shellscribe.git ~/.shellscribe
# install shellscribe binary
# may need to export GOBIN=~/go/bin
cd ~/.shellscribe && go install ~/.shellscribe/cmd/shellscribe.go && cd -
# add shellscribe-on/shellscribe-off commands to your shell
echo "source ~/.shellscribe/bin/shellscribe-init.sh" >> ~/.profile

# start the web server to view logs
shellscribe server

# activate shellscribe to record commands
shellscribe-on

# run some commands
echo "this is captured by shellscribe!"

# stop recording commands
shellscribe-off
```

![image](/docs/images/shellscribe-server.png)
