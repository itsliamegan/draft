# Draft

An experimental tool for quickly prototyping designs for web projects. Serves
HTML and CSS files in a directory to the browser, automatically updating the
page when those files change. Similar to many other development tools, but
distributed as a single binary with no dependencies, and built for markup-first
projects.

## Installation

1. Clone the repository and `cd` into the directory
2. Run `go build`
3. Optionally, move the newly created `draft` executable to a directory in your
`PATH` such as `/usr/local/bin/`

## Usage

1. Call the `draft` executable, optionally passing in a directory which contains
the HTML and CSS files (which defaults to the current directory). Run `draft
--help` for more information
2. Go to `localhost:4000/[my_html_file].html` in your browser
3. Make changes to your HTML and CSS files, and watch the browser automatically
update

## Architecture

Draft is composed of essentially four parts: the server, the watcher, the
announcer, and the listener. The server, the watcher, and the announcer are all
aspects of the program written in Go which run in the terminal, and the listener
is a script which runs in the browser.

- The **server** is a static file server, responding to web requests matching
the names of files it has been instructed to serve
- The **watcher** polls the filesystem every few-hundred milliseconds to
determine if there have been any changes, and passes them to the announcer
- The **announcer** publishes changes from the watcher to an
[EventSource][event_source] which is hooked into by the listener
- The **listener** subscribes to changes from the EventSource and updates the
page to reflect the changed HTML and CSS

## Caveats

Draft is still a prototype; it hasn't been extensively tested outside of
personal use. Some common scenarios haven't been implemented yet, and will
simply crash the program. CPU and filesystem resources, though they've been fine
in my own usage, might be improperly or excessively used. Security has been
almost completely ignored, since this is a tool for use on a personal
development machine.

Known causes of crashes:

- Having a directory present in the target directory
- Trying to visit the URL of a directory

Known limitations:
- Changes to the document `<head>`, though they are announced, are not processed
by the listener

[event_source]: https://developer.mozilla.org/en-US/docs/Web/API/EventSource
