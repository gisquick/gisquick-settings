import axios from 'axios'
import https from 'https'

const MB = 1024 * 1024
const GB = 1024 * MB
const HTTP = axios.create({
  baseURL: '',
  withCredentials: true,
  httpsAgent: new https.Agent({
    rejectUnauthorized: false
  }),
  maxContentLength: 1 * GB
})

HTTP.appendParams = function (url, params) {
  const u = new URL(url)
  for (let key in params) {
    if (params[key] !== undefined) {
      u.searchParams.append(key, params[key])
    }
  }
  return u.href
}

export default HTTP

