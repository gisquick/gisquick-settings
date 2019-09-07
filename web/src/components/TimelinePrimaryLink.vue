<template>
  <v-layout column>
    <v-timeline-item
      large
      class="py-0 my-2"
      :color="color"
      :icon="icon"
    >
      <v-btn
        :to="to"
        :color="isActive ? color : ''"
        :dark="isActive && !disabled"
        active-class=""
        rounded class="step"
        :disabled="disabled"
        @click="$emit('click', $event)"
      >
        <slot>{{ label }}</slot>
      </v-btn>
    </v-timeline-item>
    <v-expand-transition v-if="links">
      <div v-if="showSublinks">
        <div
          v-for="(link, i) in displayedLinks"
          :key="link.page"
          class="expand-container"
        >
          <v-slide-y-transition>
            <v-layout
              v-if="renderLinks"
              :style="`transition-delay: ${0.1 * i}s`"
            >
              <timeline-secondary-link
                :label="link.text"
                :to="{name: to.name, params: {page: link.page}}"
                :color="color"
              />
            </v-layout>
          </v-slide-y-transition>
        </div>
      </div>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import TimelineSecondaryLink from '@/components/TimelineSecondaryLink'
import { setTimeout } from 'timers';

export default {
  components: { TimelineSecondaryLink },
  props: {
    label: String,
    to: Object,
    icon: String,
    color: String,
    disabled: Boolean,
    links: Array
  },
  data () {
    return {
      renderLinks: false
    }
  },
  computed: {
    linkPath () {
      return this.to && this.$router.resolve(this.to).route.fullPath
    },
    isActive () {
      return this.$route.fullPath === this.linkPath || this.$route.fullPath.startsWith(this.linkPath + '/')
    },
    showSublinks () {
      return !this.disabled && this.isActive
    },
    displayedLinks () {
      return this.showSublinks ? this.links : []
    }
  },
  watch: {
    showSublinks: {
      immediate: true,
      handler (show) {
        if (show) {
          setTimeout(() => {
            this.renderLinks = true
          }, 30)
        } else {
          this.renderLinks = false
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.expand-container {
  height: 2.5em;
}
.v-timeline-item {
  .v-btn {
    width: 100%;
    height: 2.5em;
    margin-top: 0.28em;
    text-transform: none;
    font-size: 17px;
    justify-content: flex-start;
    // border-top-right-radius: 6px;
    // border-bottom-right-radius: 6px;
    // border-bottom-left-radius: 6px;
    &:after {
      content: "";
      position: absolute;
      width: 2em;
      border-bottom: 2px dotted #ccc;
      left: -2em;
      box-shadow: none;
    }
  }
}
</style>
