import morphdom from "https://cdn.skypack.dev/morphdom"

class Client {
  constructor(events) {
    this.events = events
    this.listeners = []
  }

  start() {
    this.events.addEventListener("message", event => {
      let message = JSON.parse(event.data)
      this.notifyApplicableOf(message)
    })
  }

  stop() {
    this.events.close()
  }

  addListener(listener) {
    this.listeners.push(listener)
  }

  notifyApplicableOf(message) {
    for (let listener of this.listeners) {
      if (listener.isApplicableTo(message)) {
        listener.onMessage.bind(listener)(message)
      }
    }
  }
}

class DocumentUpdateListener {
  onMessage(message) {
    let newDocument = parseHTML(message.contents)
    morphdom(document.body, newDocument.body)
  }

  isApplicableTo(message) {
    return (
      this.isDocumentUpdate(message) &&
      this.isUpdateToCurrentDocument(message)
    )
  }

  isDocumentUpdate(message) {
    return message.type === "document"
  }

  isUpdateToCurrentDocument(message) {
    return message.file == location.pathname.substring(1)
  }
}

class StylesheetUpdateListener {
  onMessage(message) {
    for (let oldLink of this.stylesheetLinks) {
      let parent = oldLink.parentNode
      let newLink = oldLink.cloneNode()

      newLink.href = oldLink.href.replace(/\?.*|$/, "?revision=" + Date.now())
      newLink.addEventListener("load", () => {
        oldLink.remove()
      })

      parent.appendChild(newLink)
    }
  }

  isApplicableTo(message) {
    return this.isStylesheetUpdate(message)
  }

  isStylesheetUpdate(message) {
    return message.type == "stylesheet"
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
