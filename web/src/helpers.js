import Vue from 'vue'
import isFunction from 'lodash/isFunction'

export function createMap (list, key, value) {
  const map = {}
  const keyFn = isFunction(key)
  const valueFn = isFunction(value)
  list.forEach(i => {
    const attr = keyFn ? key(i) : i[key]
    map[attr] = valueFn ? value(i) : i[value]
  })
  return map
}

export function setNestedProperty (obj, key, value) {
  const props = key.split('.')
  const last = props.pop()
  let dest = obj
  for (const prop of props) {
    if (!dest.hasOwnProperty(prop)) {
      Vue.set(dest, prop, {})
    }
    dest = dest[prop]
  }
  Vue.set(dest, last, value)
}

export function deleteNestedProperty (obj, key) {
  const props = key.split('.')
  const last = props.pop()
  let dest = obj
  for (const prop of props) {
    if (!dest.hasOwnProperty(prop)) {
      return
    }
    dest = dest[prop]
  }
  Vue.delete(dest, last)
}