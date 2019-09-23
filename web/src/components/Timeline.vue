<template>
  <div class="timeline my-2 mx-2">
    <v-layout class="desktop column align-center">
      <icon
        name="desktop-qgis"
        color="lime darken-2"
        class="mb-2"
        size="32px"
      />
      <svg height="32" viewBox="0 0 100 100">
        <path d="M 49 0 V 100" stroke="white" stroke-dasharray="4,4"/>
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
      >
        <path d="M 100 0 Q 0 50 100 100" stroke="white" stroke-dasharray="4,4"/>
      </svg>
    </div>

    <!-- Connection: Files -> Settings -->
    <div
      style="grid-row: 4/7; grid-column:2/3; position: relative"
      class="fill-height"
    >
      <svg
        viewBox="0 0 100 100"
        preserveAspectRatio="none"
        width="100%"
        height="100%"
        style="padding: 24px 0; position: absolute"
      >
        <path d="M 0 0 Q 100 50 0 100" stroke="white" stroke-dasharray="4,4"/>
      </svg>
    </div>

    <icon
      name="timeline-checkin"
      class="bullet r-2"
      size="48px"
    />
    <v-btn class="r-2" text dark disabled>
      Check-In
    </v-btn>

    <div class="space r-3"/>

    <icon
      name="timeline-files"
      class="bullet r-4"
      size="56px"
    />
    <v-btn
      text dark
      class="r-4"
      :to="{name: 'files'}"
    >
      Files
    </v-btn>

    <div class="space r-5"/>

    <icon
      name="timeline-settings"
      class="bullet r-6"
      size="48px"
    />
    <v-btn
      text dark
      class="r-6"
      :to="{name: 'settings'}"
    >
      Settings
    </v-btn>

    <!-- <svg
      height="32"
      viewBox="0 0 100 100"
      style="grid-row: 7/8; grid-column: 2/3"
    >
      <path d="M 1 0 V 100" stroke="white" stroke-dasharray="4,4"/>
    </svg> -->

    <div class="space r-7"/>
    <div style="grid-row: 7/8; grid-column: 2/3; position: relative">
      <svg
        height="100%"
        viewBox="0 0 100 100"
        style="position: absolute;"
      >
        <path d="M 1 0 V 100" stroke="white" stroke-dasharray="4,4"/>
      </svg>
    </div>

    <v-expand-transition>
      <v-layout
        v-if="$route.name === 'settings'"
        class="r-7 column submenu"
        style="grid-column: 3/4;"
      >
        <v-btn
          text small dark
          :to="{name: 'settings', params: {page: 'project'}}"
          active-class="primary--text"
          >
            <v-icon small class="mr-1">keyboard_arrow_right</v-icon>Project
          </v-btn>
        <v-btn
          text small dark
          :to="{name: 'settings', params: {page: 'layers'}}"
          active-class="primary--text"
        >
          <v-icon small class="mr-1">keyboard_arrow_right</v-icon>Layers
        </v-btn>
        <v-btn
          text small dark
          :to="{name: 'settings', params: {page: 'topics'}}"
          active-class="primary--text"
        >
          <v-icon small class="mr-1">keyboard_arrow_right</v-icon>Topics
        </v-btn>
      </v-layout>
    </v-expand-transition>

    <icon
      name="timeline-publish"
      class="bullet r-8"
      size="48px"
    />
    <v-btn
      text dark disabled
      class="r-8"
    >
      Publish
    </v-btn>

    <!-- Left server border -->
    <div style="grid-column: 1/2; grid-row: 5/8; position:relative" class="fill-height">
      <svg
        viewBox="0 0 32 32"
        height="100%"
        width="100%"
        style="max-height:100%"
        preserveAspectRatio="none"
      >
        <path
          d="M 0.5 0 V 31"
          style="stroke:#00a6e3;stroke-width: 2px;stroke-linecap:round;stroke-linejoin:round;stroke-dasharray: 2.5px,3.8px"
        />
      </svg>
    </div>
    <!-- Right server border -->
    <div style="grid-column: 2/3; grid-row: 5/8; position:relative" class="fill-height">
      <svg
        viewBox="0 0 32 32"
        height="100%"
        width="100%"
        style="max-height:100%"
        preserveAspectRatio="none"
      >
        <path
          d="M 31.5 0 V 31"
          style="stroke:#00a6e3;stroke-width: 2px;stroke-linecap:round;stroke-linejoin:round;stroke-dasharray: 2.5px,3.8px"
        />
      </svg>
    </div>

    <icon
      name="timeline-browsers"
      class="browsers"
      size="64px"
    />
  </div>
</template>
<script>
export default {
  
}
</script>

<style lang="scss" scoped>
.timeline {
  display: grid;
  grid-template-columns: 64px 64px 1fr;

  .desktop {
    grid-column: 1 / 3;
  }
  .bullet {
    grid-column: 1 / 3;
    position: relative;
  }
  .v-btn {
    grid-column: 3 / 4;
    margin: 6px 12px;
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
    }
  }
  .browsers {
    // grid-row: 9 / 10;
    grid-column: 1 / 3;
  }
  path {
    fill: none;
    stroke-width: 1;
    vector-effect: non-scaling-stroke;
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
  grid-column: 1 / 4;
}
</style>