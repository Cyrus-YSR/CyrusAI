<template>
  <button
    :type="type"
    class="group relative cursor-pointer overflow-hidden rounded-full border bg-background px-6 py-2 text-center font-semibold interactive-hover-button"
    :class="[className, { 'is-disabled': disabled }]"
    :disabled="disabled"
    @click="$emit('click', $event)"
  >
    <span class="inline-block transition-all duration-300 group-hover:translate-x-12 group-hover:opacity-0">
      {{ text }}
    </span>
    <div class="absolute inset-0 z-10 flex items-center justify-center gap-2 bg-primary text-primary-foreground opacity-0 transition-all duration-300 group-hover:opacity-100 rounded-full">
      <span>{{ text }}</span>
      <slot name="icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
      </slot>
    </div>
  </button>
</template>

<script setup>
defineProps({
  text: {
    type: String,
    default: 'Button'
  },
  type: {
    type: String,
    default: 'button'
  },
  className: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  }
})

defineEmits(['click'])
</script>

<style scoped>
.interactive-hover-button {
  width: 100%;
  height: 48px;
  background-color: #ffffff;
  border: 1px solid rgba(229, 231, 235, 0.6);
  color: #000000;
  border-radius: 9999px; /* full rounded */
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.3s ease;
  cursor: pointer;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  white-space: nowrap;
}

.interactive-hover-button.is-disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.inline-block {
  display: inline-block;
}

.group-hover\:translate-x-12 {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.group-hover\:opacity-0 {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.group:hover .group-hover\:translate-x-12 {
  transform: translateX(3rem);
}

.group:hover .group-hover\:opacity-0 {
  opacity: 0;
}

.absolute {
  position: absolute;
}

.inset-0 {
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
}

.z-10 {
  z-index: 10;
}

.flex {
  display: flex;
}

.items-center {
  align-items: center;
}

.justify-center {
  justify-content: center;
}

.gap-2 {
  gap: 0.5rem;
}

.bg-primary {
  background-color: #111827; /* Dark background on hover */
}

.text-primary-foreground {
  color: #ffffff; /* White text on hover */
}

.opacity-0 {
  opacity: 0;
}

.transition-all {
  transition-property: all;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 300ms;
}

.duration-300 {
  transition-duration: 300ms;
}

.group:hover .group-hover\:opacity-100 {
  opacity: 1;
}

.rounded-full {
  border-radius: 9999px;
}

.h-4 {
  height: 1rem;
}

.w-4 {
  width: 1rem;
}
</style>
