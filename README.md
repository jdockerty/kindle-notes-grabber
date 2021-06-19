# Kindle Notes Grabber

A service to check for `Amazon Kindle` emails in your account that contain attached notes from a particular book, outputting these in a neat format to `$HOME/kindle-notes/<book-title>-notes.txt`. Once you've finished a book, run this application to have your notes saved to your computer, do whatever you wish with them after!

This repository also offers various packages for dealing with IMAP and Amazon Kindle emails with an attached CSV.

## Install

Build from the repository using Go

```bash
git clone https://github.com/jdockerty/kindle-notes-grabber.git
cd kindle-notes-grabber
go build -o kng
sudo mv kn /usr/local/bin
```

## Usage

This program uses either a configuration file, named `kng-config.yaml`, in your home directory or utilises the `KNG_EMAIL` and `KNG_PASSWORD` environment variables.

If you choose to use a configuration file, it should look something like this

```yaml
# $HOME/kng-config.yaml
email: your-email@gmail.com
password: <app or account password>
```

Ensure you have the generated binary saved into your `PATH`, this means that you can call it from the command line directly.

```bash
kng run
```
*At present, this is setup to work with Gmail accounts (more will be added in the future). It is recommended to create an [app password](https://support.google.com/accounts/answer/185833?hl=en) to use as your method of signing in using the app, although using your regular password also works.*