<template>
  <topics-editor
    :topics="config.topics"
    :layers="config.layers"
  />
</template>

<script>
import TopicsEditor from '@/components/TopicsEditor'

export default {
  name: 'Topics',
  components: { TopicsEditor },
  props: {
    config: Object
  }
}
</script>
