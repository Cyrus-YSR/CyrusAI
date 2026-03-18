<template>
  <div
    ref="pupilRef"
    class="pupil-wrapper"
    :style="wrapperStyle"
  ></div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  size: { type: Number, default: 12 },
  maxDistance: { type: Number, default: 5 },
  pupilColor: { type: String, default: 'black' },
  forceLookX: { type: Number, default: undefined },
  forceLookY: { type: Number, default: undefined },
  mouseX: { type: Number, default: 0 },
  mouseY: { type: Number, default: 0 }
})

const pupilRef = ref(null)

const calculatePupilPosition = () => {
  if (!pupilRef.value) return { x: 0, y: 0 }

  if (props.forceLookX !== undefined && props.forceLookY !== undefined) {
    return { x: props.forceLookX, y: props.forceLookY }
  }

  const rect = pupilRef.value.getBoundingClientRect()
  const centerX = rect.left + rect.width / 2
  const centerY = rect.top + rect.height / 2

  const deltaX = props.mouseX - centerX
  const deltaY = props.mouseY - centerY
  const distance = Math.min(Math.sqrt(deltaX ** 2 + deltaY ** 2), props.maxDistance)

  const angle = Math.atan2(deltaY, deltaX)
  const x = Math.cos(angle) * distance
  const y = Math.sin(angle) * distance

  return { x, y }
}

const pupilPosition = computed(() => calculatePupilPosition())

const wrapperStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  backgroundColor: props.pupilColor,
  transform: `translate(${pupilPosition.value.x}px, ${pupilPosition.value.y}px)`
}))
</script>

<style scoped>
.pupil-wrapper {
  border-radius: 50%;
  transition: transform 0.1s ease-out;
}
</style>
