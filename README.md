# Kindle Notes Grabber

A command-line app to check for `Amazon Kindle` emails in your account that contain attached CSV notes from a particular book, outputting these in a neat format to `$HOME/kindle-notes/<book-title>-notebook.txt`. Once you've finished a book, run this application to have your notes saved to your computer, do whatever you wish with them after!

This repository also offers a package for dealing with IMAP and Amazon Kindle emails with an attached CSV, namely in the `notes` package.

## Install

Build from the repository using Go

```bash
git clone https://github.com/jdockerty/kindle-notes-grabber.git
cd kindle-notes-grabber
go build -o kng
sudo mv kng /usr/local/bin
```

## Usage

This program uses either a configuration file, named `kng-config.yaml`, in your home directory or specified path via the `--config` flag. It can also check for the `KNG_EMAIL` and `KNG_PASSWORD` environment variables, if this file is not in use.

If you choose to use a configuration file, it should look something like this

```yaml
# $HOME/kng-config.yaml
email: your-email@gmail.com
password: <app or account password>
```

With environment variables, it would look something like

```bash
export KNG_EMAIL="youremail@gmail.com"
export KNG_PASSWORD="app_or_account_pass"
```

Ensure you have the generated binary saved into your `PATH`, this means that you can call it from the command line directly.

To save some time and remove any setup headaches, the `setup` sub-command can be used to create the relevant folder and file used when the main application is called. To do this, simply run

```bash
kng setup
```

After, you can simply do
```bash
kng run
```

To have the application parse your mailbox.

Run `--help` or `-h` on any command to see its explanation and possible flags or sub-commands.

### Note
It is recommended to create an [app password](https://support.google.com/accounts/answer/185833?hl=en) to use as your method of signing in using the app if you have 2FA enabled, when using Gmail. Otherwise using your regular password also works, but you may need to enable 'Insecure applications' to access your account, so it is recommended to have a generated app password which you can revoke at any time. The other providers which are currently supported are Outlook and Yahoo.
