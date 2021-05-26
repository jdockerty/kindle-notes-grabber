# Kindle Notes Grabber

A daemon service to periodically check for `Amazon Kindle` emails in your account that contain attached notes from a particular book, outputting these in a neat format to `$HOME/kn/<book-title>-notes.txt`.


## TODO

- Read environment/configuration file for email log in.
- Unit tests [on-going]
- Look into running app without need for "insecure apps" enabled on Gmail (others later)