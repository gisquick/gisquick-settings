<template>
  <div class="timeline my-2 mx-2">
    <v-layout class="desktop column align-center">
      <desktop-svg style="opacity:0.9"/>
      <svg height="32" viewBox="0 0 100 100" class="line">
        <path d="M 49 0 V 100" stroke="white" stroke-dasharray="4,4"/>
        <path
          :d="`M 49 0 V ${checkin.disabled ? 0 : 100}`"
          style="transition: all 0.5s"
          stroke="white"
          stroke-width="2"
        />
      </svg>
    </v-layout>

    <!-- Connection: Check-in -> Files -->
    <div style="grid-row: 2/5; grid-column:1/2; position: relative">
      <svg
        viewBox="0 0 100 100"
        preserveAspectRatio="none"
        width="100%"
        height="100%"
        style="padding:24px 0; position: absolute;"
        class="line"
      >
        <path
          d="M 50 0 Q 0 50 50 100"
          stroke="white"
          stroke-dasharray="4,4"
        />
        <!-- Progress line -->
        <path
          d="M 50 0 Q 0 50 50 100"
          class="progress-line"
          :stroke-dasharray="`${!visited[files.link.name] ? 0 : 120}, 500`"
        />
      </svg>

      <timeline-menu-lines
        v-if="activeItem === checkin"
        :x="81"
        :ymax="75"
        :number="2"
        :step="30"
        :duration="500"
        style="position: absolute; padding-top: 24px; padding-bottom:47px;"
        class="line secondary-link"
      />
    </div>

    <!-- Connection: Files -> Settings -->
    <div
      style="grid-row: 4/7; grid-column:1/2; position: relative"
      class="fill-height"
    >
      <svg
        viewBox="0 0 100 100"
        preserveAspectRatio="none"
        width="100%"
        height="100%"
        style="padding: 24px 0; position: absolute"
        class="line"
      >
        <path
          d="M 50 0 Q 100 50 50 100"
          stroke="white"
          stroke-dasharray="4,4"
        />
        <path
          d="M 50 0 Q 100 50 50 100"
          class="progress-line"
          :stroke-dasharray="`${!visited[settings.link.name] ? 0 : 120}, 500`"
        />
      </svg>
    </div>

    <checkin-svg
      class="bullet r-2"
      :class="{disabled: checkin.disabled, active: activeItem === checkin}"
    />
    <div class="space r-3"/>
    <files-svg
      class="bullet r-4"
      :class="{disabled: files.disabled, active: activeItem === files}"
    />
    <div class="space r-5"/>

    <!-- Connection: Settings -> Publish -->
    <div style="grid-row: 6/8; grid-column: 1/2; position: relative">
      <svg
        height="100%"
        width="100%"
        viewBox="0 0 100 100"
        preserveAspectRatio="none"
        style="position: absolute; padding-top: 24px"
        class="line"
      >
        <path d="M 49 0 V 100" stroke="white" stroke-dasharray="4,4"/>
        <path
          :d="`M 49 0 V ${publish.published ? 100 : 0}`"
          style="transition: all 0.5s"
          stroke="white"
          stroke-width="2"
        />
      </svg>
      <timeline-menu-lines
        v-if="activeItem === settings"
        :number="3"
        :duration="700"
        style="position: absolute; padding-top: 24px"
        class="line secondary-link"
      />
    </div>

    <settings-svg
      class="bullet r-6"
      :class="{disabled: settings.disabled, active: activeItem === settings}"
    />
    <div class="space r-7"/>

    <publish-svg
      class="bullet r-8"
      :class="{disabled: publish.disabled}"
    />

    <!-- Side server border -->
    <div style="grid-column: 1/2; grid-row: 5/8; position:relative" class="fill-height">
      <svg
        viewBox="0 0 100 100"
        height="100%"
        width="100%"
        style="max-height:100%"
        preserveAspectRatio="none"
        class="line"
      >
        <path
          d="M 0.5 0 V 100
             M 92.5 0 V 100"
          style="stroke: #00a6e3;
                 stroke-width: 2;
                 stroke-linecap: round;
                 stroke-linejoin: round;
                 stroke-dasharray: 2.5,3.8"
        />
      </svg>
    </div>

    <!-- Links -->
    <template v-for="(item, index) in [checkin, files, settings]">
      <v-btn
        :key="`link-${index}`"
        :class="`r-${2 + 2 * index}`"
        active-class="orange--text"
        :disabled="item.disabled"
        :to="item.link"
        text dark
      >
        {{ item.label }}
      </v-btn>
      <v-expand-transition
        v-if="item.sublinks"
        :key="`sublink-${index}`"
      >
        <v-layout
          v-if="activeItem === item"
          class="column submenu"
          :class="`r-${3 + 2 * index}`"
          style="grid-column: 2/3;"
        >
          <v-btn
            v-for="link in item.sublinks"
            :key="link.page"
            :to="{name: item.link.name, params: {page: link.page}}"
            active-class="orange--text"
            text small dark
          >
            <v-icon small class="mr-1">keyboard_arrow_right</v-icon>
            <span class="grow">{{ link.label }}</span>
            <v-icon
              v-if="link.status === 'error'"
              v-text="'report'"
              class="red--text text--darken-1"
              size="18"
            />
          </v-btn>
        </v-layout>
      </v-expand-transition>
    </template>

    <v-btn
      text dark
      class="r-8"
      :disabled="publish.disabled"
      @click="$emit('publish')"
    >
      {{ publish.label }}
    </v-btn>

    <browsers-svg style="opacity: 0.6"/>
  </div>
</template>

<script>
import CheckinSvg from '@/assets/inline/timeline-checkin.svg'
import FilesSvg from '@/assets/inline/timeline-files.svg'
import SettingsSvg from '@/assets/inline/timeline-settings.svg'
import PublishSvg from '@/assets/inline/timeline-publish.svg'
import DesktopSvg from '@/assets/inline/desktop-qgis.svg'
import BrowsersSvg from '@/assets/inline/timeline-browsers.svg'
import TimelineMenuLines from '@/components/TimelineMenuLines'

export default {
  name: 'timeline',
  components: { CheckinSvg, FilesSvg, SettingsSvg, PublishSvg, DesktopSvg, BrowsersSvg, TimelineMenuLines },
  props: {
    checkin: {
      type: Object,
      default: () => ({
        disabled: true,
        label: 'Check-in'
      })
    },
    files: {
      type: Object,
      default: () => ({
        label: 'Files'
      })
    },
    settings: {
      type: Object,
      default: () => ({
        label: 'Settings'
      })
    },
    publish : {
      type: Object,
      default: () => ({
        label: 'Publish'
      })
    },
    visited: {
      type: Object,
      default: () => ({})
    }
  },
  computed: {
    links () {
      return [this.checkin, this.files, this.settings]
    },
    activeItem () {
      return this.links.find(i => i && !i.disabled && i.link && this.$route.name === i.link.name)
    }
  }
}
</script>

<style lang="scss" scoped>
.timeline {
  display: grid;
  grid-template-columns: 128px 1fr;

  .bullet {
    grid-column: 1 / 2;
    position: relative;

    #bg {
      transition: all 0.5s;
      stroke-dasharray: none;
      stroke: transparent;
      stroke-width: 0;
      stroke: #ff9800;
      vector-effect: non-scaling-stroke;
    }
    &.active {
      #bg {
        // fill: #f5a630;
        stroke-width: 2px;
      }
    }
  }
  .v-btn {
    grid-column: 2 / 3;
    margin: 6px 12px;
    &.v-btn--active::before {
      // display: none;
      opacity: 0.12;
    }
    ::v-deep .v-btn__content {
      display: inline;
      text-align: left;
    }
  }
  .submenu {
    .v-btn {
      margin-top: 2px;
      margin-bottom: 2px;
      text-transform: none;
      &:before {
        display: none;
      }
      ::v-deep .v-btn__content {
        display: flex;
      }
    }
  }
  svg {
    &.line path {
      fill: none;
      vector-effect: non-scaling-stroke;
    }
    &.secondary-link {
      stroke: #ff9800;
      stroke-dasharray: 4,1;
    }
    .progress-line {
      transition: all 0.5s;
      stroke: white;
      stroke-width: 2;
    }
    #mask {
      fill: #444;
      fill-opacity: 0;
      transition: all 0.5s;
    }
    &.disabled #mask {
      fill-opacity: 0.6;
    }
  }
}
.r-1 {
  grid-row: 1 / 2;
}
.r-2 {
  grid-row: 2 / 3;
}
.r-3 {
  grid-row: 3 / 4;
}
.r-4 {
  grid-row: 4 / 5;
}
.r-5 {
  grid-row: 5 / 6;
}
.r-6 {
  grid-row: 6 / 7;
}
.r-7 {
  grid-row: 7 / 8;
}
.space {
  height: 20px;
  grid-column: 1 / 3;
}
.expand-transition-enter-active,
.expand-transition-leave-active {
  transition-duration: 0.6s
}
</style>
