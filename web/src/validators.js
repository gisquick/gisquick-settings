export function required (value) {
  if (!value) {
    return 'Field is required'
  }
  return true
}
