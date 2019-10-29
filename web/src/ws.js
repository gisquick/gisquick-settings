
export default function WebsocketMessenger (url) {
  let listeners = []
  const activeRequests = {}
  const socket = new WebSocket(url)
  let timer = null
  const ws = {
    connected: false,
    pluginConnected: false,

    bind (type, callback) {
      listeners.push({ type, callback })
      return callback
    },
    unbind (type, callback) {
      listeners = listeners.filter(l => l.type !== type || l.callback !== callback)
    },
    send (msg) {
      socket.send(msg)
    },
    sendJSON (name, data) {
      socket.send(name + ":" + JSON.stringify(data))
    },
    request (name, data = null) {
      return new Promise((resolve, reject) => {
        activeRequests[name] = { resolve, reject }
        if (data) {
          this.sendJSON(name, data)
        } else {
          this.send(name)
        }
      })
    },
    close () {
      if (timer !== null) {
        clearInterval(timer)
        timer = null
      }
      socket.close()
    }
  }

  socket.onopen = () => {
    ws.connected = true
    socket.send('PingPlugin')
  }
  socket.onclose = () => {
    ws.connected = false
    
  }
  socket.onmessage = (e) => {
    if (e.data === 'PongPlugin') {
      ws.pluginConnected = true
      return
    }
    if (e.data === 'PluginConnected') {
      ws.pluginConnected = true
    }
    if (e.data === 'PluginDisconnected') {
      ws.pluginConnected = false
    }

    const breakpoint = e.data.indexOf(':')
    let type, data
    if (breakpoint !== -1) {
      type = e.data.substring(0, breakpoint)
      data = e.data.substring(breakpoint + 1)
    } else {
      type = e.data
      data = ''
    }
    if (activeRequests[type]) {
      activeRequests[type].resolve(data)
      delete activeRequests[type]
    }
    listeners.filter(l => l.type === type).forEach(l => {
      l.callback(e, data)
    })
  }
  timer = setInterval(() => {
    socket.send('Ping')
  }, 25 * 1000)
  return ws
}
