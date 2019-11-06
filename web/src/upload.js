
export function createUpload (ws, files, project) {

  const info = {
    files: {},
    totalProgress: 0
  }
  let totalSize = 0
  files.forEach(f => {
    totalSize += f.size
    info.files[f.path] = {
      size: f.size,
      progress: 0
    }
  })

  let task
  function onProgressMessage(msg) {
    const data = msg.data

    Object.entries(data).forEach(([file, progress]) => {
      const fileUpload = info.files[file]
      if (fileUpload) {
        fileUpload.progress = 100 * (progress / fileUpload.size)
      }
    })

    const uploaded = Object.values(info.files).reduce((sum, info) => sum + info.progress * info.size, 0)
    info.totalProgress = totalSize !== 0 ? uploaded / totalSize : 100

    const finished = info.totalProgress === 100
    if (task && task.onProgress) {
      task.onProgress(info)
    }
    if (finished) {
      ws.unbind('UploadProgress', onProgressMessage)
      task.resolve()
      task = null
    }
  }

  function onErrorMessage (msg) {
    task.reject(msg.data.trim())
  }

  return {
    info,
    start (onProgress) {
      return new Promise((resolve, reject) => {
        ws.bind('UploadProgress', onProgressMessage)
        ws.bind('UploadError', onErrorMessage)
        ws.send('UploadFiles', { files, project })
        task = {
          resolve,
          reject,
          onProgress
        }
      })
    },
    abort () {
      if (task) {
        ws.unbind('UploadProgress', onProgressMessage)
        ws.unbind('UploadError', onErrorMessage)
        task.reject('aborted')
        ws.send('AbortUpload', { project })
        task = null
      }
    }
  }
}
