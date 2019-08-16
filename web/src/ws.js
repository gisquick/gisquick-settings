
export default function WebsocketMessenger (url) {
  let listeners = []
  const socket = new WebSocket(url)
  let timer = null
  const ws = {
    connected: false,
    pluginConnected: false,

    bind (type, callback) {
      listeners.push({ type, callback })
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
    if (e.data === 'PluginDisconnected') {
      ws.pluginConnected = false
      return
    }
    if (e.data === 'PluginConnected' || e.data === 'PongPlugin') {
      ws.pluginConnected = true
      return
    }
    const breakpoint = e.data.indexOf(':')
    const type = e.data.substring(0, breakpoint)
    const data = e.data.substring(breakpoint + 1)
    listeners.filter(l => l.type === type).forEach(l => {
      l.callback(e, data)
    })
  }
  timer = setInterval(() => {
    socket.send('Ping')
  }, 25 * 1000)
  return ws
}
