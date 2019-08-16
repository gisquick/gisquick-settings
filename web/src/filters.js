const SizeUnits = {
  GB: value => (value / 1073741824).toFixed(2),
  MB: value => (value / 1048576).toFixed(2),
  kB: value => (value / 1024).toFixed(2),
  B: value => Math.round(value)
}

export function formatFileSize (value, unit) {
  if (!value) {
    return value
  }
  if (!unit) {
    if (value > 1073741824) {
      unit = 'GB'
    } else if (value > 1048576) {
      unit = 'MB'
    } else if (value > 1024) {
      unit = 'kB'
    } else {
      unit = 'B'
    }
  }
  return `${SizeUnits[unit](value)} ${unit}`
}
