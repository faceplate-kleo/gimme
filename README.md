# gimme
A CLI tool for jumping to project contexts, setting environment variables, and running commands
___

> [!CAUTION]
> This tool is essentially malware if used irresponsibly. It is capable of arbitrary command execution _**BY DESIGN**_.
> Half of the code in this project is designed to circumvent the very reasonable security features of `zsh` and `bash` in order to produce the desired effect.
> I highly recommend you review the source code before making the decision to install it.
>
> If you're cool with all of that, read on!

## Prerequisites

- Either `bash` or `zsh`. Other shells may be compatible, but they are as yet untested.
- Root access to your system
- [Go >=v1.22.2](https://go.dev/doc/install)

## Installation
Clone this repository, and cd into it:
```
git clone https://github.com/faceplate-kleo/gimme && cd gimme
```
Run the installation script `INSTALL.sh`. It will automatically attempt to elevate to root privileges.
```
./INSTALL.sh
```
Source your edited shell runtime configuration file as specified by the installation script. For instance, in a `zsh` shell:
```
source ~/.zshrc
```
Initialize gimme. This will create the configuration directory, at `$HOME/.gimme` by default.
```
gimme init
```

You're done - happy gimming!

## Usage
Gimme works by discovering `.gimme.yaml` configuration files and associating paths to those files with aliases.
When gimme is called with a recognized alias as an argument, it automatically moves the shell session to the configuration file's location.
Once there, it will set any environment variables specified in the configuration file, and then prompt you for permission to run any requested shell commands.

Below is a barebones example of a `.gimme.yaml`. It binds the parent directory to the alias `example`, sets environment variables, and performs generally-safe shell commands after warping.
```yaml
gimmeVersion: v1
gimme:
  alias: example
  env:
    HELLO: world
    KUBECONFIG: ~/.kube/kubeconfig-bogus.yml
  init:
    - "echo 'hello world'"
    - "git status"
```

Once the file is created, you'll first need to tell gimme to refresh its alias bindings like so:
```
gimme sync
```

If you'd like to see all the visible `.gimme.yaml` files on your system, run the following:
```
gimme discover
```
This performs the same walk as `gimme sync`, but does not update any bindings. It simply prints out the full paths to the configuration files.


Once your gimme aliases have been synced, simply call the alias from anywhere!
```
gimme example
```

To view all of your bound aliases, simply run:
```
gimme list
```

## Configuration

Gimme's global configuration can be found at `$HOME/.gimme/config.yaml`. At time of writing, there are only two configuration options that matter: `homeOverride` and `autoTrust`

An example config.yaml:
```yaml
---
homeOverride: /path/to/another/home
autoTrust: false
```

`homeOverride` changes the behavior of gimme's filesystem walk to start at the specified directory, rather than $HOME. 
This can be useful if you have a very cluttered home directory, or if you simply want to limit gimme's scope of vision.

An absolute path is required here.

`autoTrust` skips the user input step when running commands specified in `.gimme.yaml`, thus *automatically trusting* all gimme files on the system. Use this option with caution.