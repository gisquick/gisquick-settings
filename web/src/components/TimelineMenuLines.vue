<template>
  <svg
    height="100%"
    width="100%"
    viewBox="0 0 100 100"
    preserveAspectRatio="none"
  >
    <path
      v-for="(line, index) in additionalLines"
      :key="index"
    >
      <animate
        :begin="`${delay + line.delay}ms`"
        :dur="`${stepTime}ms`"
        attributeName="d"
        fill="freeze"
        :values="line.values"
      />
    </path>
    <path>
      <animate
        :begin="`${delay}ms`"
        :dur="`${duration}ms`"
        attributeName="d"
        fill="freeze"
        :keyTimes="mainLineTimes"
        :values="mainLineValues"
      />
    </path>
  </svg>
</template>

<script>

export default {
  props: {
    number: Number,
    x: {
      type: Number,
      default: 80
    },
    xmax: {
      type: Number,
      default: 100
    },
    ymax: {
      type: Number,
      default: 85
    },
    step: {
      type: Number,
      default: 25
    },
    duration: {
      type: Number,
      default: 1000
    },
    delay: {
      type: Number,
      default: 100
    }
  },
  computed: {
    mainLineValues () {
      const { x, xmax, ymax } = this
      return [
        `M ${xmax} ${ymax} h 0 v 0 l 0 0`,
        `M ${xmax} ${ymax} H ${x} v 0 l 0 0`,
        `M ${xmax} ${ymax} H ${x} V 22 l 0 0`,
        `M ${xmax} ${ymax} H ${x} V 22 L 65 8`
      ].join(';')
    },
    stepTime () {
      return this.duration / (this.number + 1)
    },
    mainLineTimes () {
      const t1 = this.stepTime / this.duration
      const t2 = 1 - t1
      return `0; ${t1}; ${t2}; 1`
    },
    additionalLines () {
      const { x, xmax, ymax, step } = this
      const lines = new Array(this.number - 1).fill(0)
        .map((v, i) => {
          const y = ymax - ((i + 1) * step)
          const values = [
            `M ${xmax} ${y} H ${xmax}`,
            `M ${xmax} ${y} H ${x}`,
          ].join(';')
          return {
            delay: (this.stepTime * (i + 1)),
            values
          }
        })
      return lines
    }
  }
}
</script>

<style lang="scss" scoped>
path {
  fill: none;
  vector-effect: non-scaling-stroke;
}
</style>
