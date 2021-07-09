import morphdom from "https://cdn.skypack.dev/morphdom"

class Client {
  constructor(events) {
    this.events = events
    this.listeners = []
  }

  start() {
    this.events.addEventListener("message", event => {
      let change = JSON.parse(event.data)
      this.notifyListenersOf(change)
    })
  }

  stop() {
    this.events.close()
  }

  addListener(listener) {
    this.listeners.push(listener)
  }

  notifyListenersOf(change) {
    for (let listener of this.listeners) {
      if (listener.wants(change)) {
        listener.receive(change)
      }
    }
  }
}

class Browser {
  static STYLESHEET_LINKS_SELECTOR = "link[rel=stylesheet]"

  load(newDocument) {
    morphdom(document.body, newDocument.body)
    document.head = newDocument.head
  }

  get stylesheets() {
    let links = document.querySelectorAll(STYLESHEET_LINKS_SELECTOR)
    let stylesheets = links.map(link => new Stylesheet(link))

    return stylesheets
  }
}

class Stylesheet {
  static EXISTING_QUERY_OR_NO_QUERY_PATTERN = /\?.*|$/

  constructor(element) {
    this.element = element
  }

  refresh() {
    let query = `?revision=${Date.now()}`

    this.element.href = this.element.href.replace(
      Stylesheet.EXISTING_QUERY_OR_NO_QUERY_PATTERN,
      query
    )
  }

}

class DocumentUpdateListener {
  constructor(browser) {
    this.browser = browser
  }

  receive(change) {
    let newDocument = parseHTML(change.contents)
    this.browser.load(newDocument)
  }

  wants(change) {
    return (
      this.isDocumentUpdate(change) &&
      this.isUpdateToCurrentDocument(change)
    )
  }

  isDocumentUpdate(change) {
    return change.mimeType == "text/html"
  }

  isUpdateToCurrentDocument(change) {
    return change.filename == location.pathname.substring(1)
  }
}

class StylesheetUpdateListener {
  constructor(browser) {
    this.browser = browser
  }

  receive(change) {
    for (let stylesheet of this.browser.stylesheets) {
      stylesheet.refresh()
    }
  }

  wants(change) {
    return this.isStylesheetUpdate(change)
  }

  isStylesheetUpdate(change) {
    return change.mimeType == "text/css"
  }
}

function parseHTML(html) {
  let domParser = new DOMParser()

  return domParser.parseFromString(html, "text/html")
}

let eventsURL = new URL(location.origin)
eventsURL.port = Number.parseInt(eventsURL.port) + 1
let changeEvents = new EventSource(eventsURL, {withCredentials: true})

let browser = new Browser()
let client = new Client(changeEvents)
client.addListener(new DocumentUpdateListener(browser))
client.addListener(new StylesheetUpdateListener(browser))

client.start()
addEventListener("beforeunload", () => {
  client.stop()
})
