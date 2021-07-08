import morphdom from "https://cdn.skypack.dev/morphdom"

class Client {
  constructor(events) {
    this.events = events
    this.listeners = []
  }

  start() {
    this.events.addEventListener("message", event => {
      let message = JSON.parse(event.data)
      this.notifyListenersOf(message)
    })
  }

  stop() {
    this.events.close()
  }

  addListener(listener) {
    this.listeners.push(listener)
  }

  notifyListenersOf(message) {
    for (let listener of this.listeners) {
      if (listener.wants(message)) {
        listener.receive(message)
      }
    }
  }
}

class DocumentUpdateListener {
  receive(message) {
    let newDocument = parseHTML(message.contents)
    morphdom(document.body, newDocument.body)
    document.head = newDocument.head
  }

  wants(message) {
    return (
      this.isDocumentUpdate(message) &&
      this.isUpdateToCurrentDocument(message)
    )
  }

  isDocumentUpdate(message) {
    return message.mimeType == "text/html"
  }

  isUpdateToCurrentDocument(message) {
    return message.filename == location.pathname.substring(1)
  }
}

class StylesheetUpdateListener {
  receive(message) {
    for (let oldLink of this.stylesheetLinks) {
      let parent = oldLink.parentNode
      let newLink = oldLink.cloneNode()

      newLink.href = oldLink.href.replace(/\?.*|$/, `?revision=${Date.now()}`)
      newLink.addEventListener("load", () => {
        oldLink.remove()
      })

      parent.appendChild(newLink)
    }
  }

  wants(message) {
    return this.isStylesheetUpdate(message)
  }

  isStylesheetUpdate(message) {
    return message.mimeType == "text/css"
  }

  get stylesheetLinks() {
    return document.querySelectorAll("link[rel=stylesheet]")
  }
}

function parseHTML(html) {
  let domParser = new DOMParser()

  return domParser.parseFromString(html, "text/html")
}

let eventsURL = new URL(location.origin)
eventsURL.port = Number.parseInt(eventsURL.port) + 1
let changeEvents = new EventSource(eventsURL, {withCredentials: true})

let client = new Client(changeEvents)
client.addListener(new DocumentUpdateListener())
client.addListener(new StylesheetUpdateListener())

client.start()
addEventListener("beforeunload", () => {
  client.stop()
})
