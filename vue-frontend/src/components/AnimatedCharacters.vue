<template>
  <div class="relative characters-wrapper">
    <!-- Purple tall rectangle character - Back layer -->
    <div 
      ref="purpleRef"
      class="character purple-char"
      :style="purpleStyle"
    >
      <!-- Eyes -->
      <div class="eyes-container" :style="purpleEyesStyle">
        <EyeBall 
          :size="18" 
          :pupil-size="7" 
          :max-distance="5" 
          eye-color="white" 
          pupil-color="#2D2D2D" 
          :is-blinking="isPurpleBlinking"
          :force-look-x="purpleForceLookX"
          :force-look-y="purpleForceLookY"
          :mouse-x="mouseX"
          :mouse-y="mouseY"
        />
        <EyeBall 
          :size="18" 
          :pupil-size="7" 
          :max-distance="5" 
          eye-color="white" 
          pupil-color="#2D2D2D" 
          :is-blinking="isPurpleBlinking"
          :force-look-x="purpleForceLookX"
          :force-look-y="purpleForceLookY"
          :mouse-x="mouseX"
          :mouse-y="mouseY"
        />
      </div>
    </div>

    <!-- Black tall rectangle character - Middle layer -->
    <div 
      ref="blackRef"
      class="character black-char"
      :style="blackStyle"
    >
      <!-- Eyes -->
      <div class="eyes-container gap-6" :style="blackEyesStyle">
        <EyeBall 
          :size="16" 
          :pupil-size="6" 
          :max-distance="4" 
          eye-color="white" 
          pupil-color="#2D2D2D" 
          :is-blinking="isBlackBlinking"
          :force-look-x="blackForceLookX"
          :force-look-y="blackForceLookY"
          :mouse-x="mouseX"
          :mouse-y="mouseY"
        />
        <EyeBall 
          :size="16" 
          :pupil-size="6" 
          :max-distance="4" 
          eye-color="white" 
          pupil-color="#2D2D2D" 
          :is-blinking="isBlackBlinking"
          :force-look-x="blackForceLookX"
          :force-look-y="blackForceLookY"
          :mouse-x="mouseX"
          :mouse-y="mouseY"
        />
      </div>
    </div>

    <!-- Orange semi-circle character - Front left -->
    <div 
      ref="orangeRef"
      class="character orange-char"
      :style="orangeStyle"
    >
      <!-- Eyes - just pupils, no white -->
      <div class="eyes-container gap-8 duration-200" :style="orangeEyesStyle">
        <CharacterPupil 
          :size="12" 
          :max-distance="5" 
          pupil-color="#2D2D2D" 
          :force-look-x="orangeForceLookX"
          :force-look-y="orangeForceLookY"
          :mouse-x="mouseX"
          :mouse-y="mouseY"
        />
        <CharacterPupil 
          :size="12" 
          :max-distance="5" 
          pupil-color="#2D2D2D" 
          :force-look-x="orangeForceLookX"
          :force-look-y="orangeForceLookY"
          :mouse-x="mouseX"
          :mouse-y="mouseY"
        />
      </div>
    </div>

    <!-- Yellow tall rectangle character - Front right -->
    <div 
      ref="yellowRef"
      class="character yellow-char"
      :style="yellowStyle"
    >
      <!-- Eyes - just pupils, no white -->
      <div class="eyes-container gap-6 duration-200" :style="yellowEyesStyle">
        <CharacterPupil 
          :size="12" 
          :max-distance="5" 
          pupil-color="#2D2D2D" 
          :force-look-x="yellowForceLookX"
          :force-look-y="yellowForceLookY"
          :mouse-x="mouseX"
          :mouse-y="mouseY"
        />
        <CharacterPupil 
          :size="12" 
          :max-distance="5" 
          pupil-color="#2D2D2D" 
          :force-look-x="yellowForceLookX"
          :force-look-y="yellowForceLookY"
          :mouse-x="mouseX"
          :mouse-y="mouseY"
        />
      </div>
      <!-- Horizontal line for mouth -->
      <div class="mouth" :style="yellowMouthStyle"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import EyeBall from './EyeBall.vue'
import CharacterPupil from './CharacterPupil.vue'

const props = defineProps({
  isTyping: {
    type: Boolean,
    default: false
  },
  showPassword: {
    type: Boolean,
    default: false
  },
  passwordLength: {
    type: Number,
    default: 0
  }
})

const mouseX = ref(0)
const mouseY = ref(0)
const isPurpleBlinking = ref(false)
const isBlackBlinking = ref(false)
const isLookingAtEachOther = ref(false)
const isPurplePeeking = ref(false)

const purpleRef = ref(null)
const blackRef = ref(null)
const yellowRef = ref(null)
const orangeRef = ref(null)

const handleMouseMove = (e) => {
  mouseX.value = e.clientX
  mouseY.value = e.clientY
}

let purpleBlinkTimeout, blackBlinkTimeout, lookingTimer, peekInterval

const schedulePurpleBlink = () => {
  const interval = Math.random() * 4000 + 3000
  purpleBlinkTimeout = setTimeout(() => {
    isPurpleBlinking.value = true
    setTimeout(() => {
      isPurpleBlinking.value = false
      schedulePurpleBlink()
    }, 150)
  }, interval)
}

const scheduleBlackBlink = () => {
  const interval = Math.random() * 4000 + 3000
  blackBlinkTimeout = setTimeout(() => {
    isBlackBlinking.value = true
    setTimeout(() => {
      isBlackBlinking.value = false
      scheduleBlackBlink()
    }, 150)
  }, interval)
}

onMounted(() => {
  window.addEventListener('mousemove', handleMouseMove)
  schedulePurpleBlink()
  scheduleBlackBlink()
})

onUnmounted(() => {
  window.removeEventListener('mousemove', handleMouseMove)
  clearTimeout(purpleBlinkTimeout)
  clearTimeout(blackBlinkTimeout)
  clearTimeout(lookingTimer)
  clearTimeout(peekInterval)
})

watch(() => props.isTyping, (newVal) => {
  if (newVal) {
    isLookingAtEachOther.value = true
    clearTimeout(lookingTimer)
    lookingTimer = setTimeout(() => {
      isLookingAtEachOther.value = false
    }, 800)
  } else {
    isLookingAtEachOther.value = false
  }
})

watch([() => props.passwordLength, () => props.showPassword], ([len, show]) => {
  if (len > 0 && show) {
    clearTimeout(peekInterval)
    const schedulePeek = () => {
      peekInterval = setTimeout(() => {
        isPurplePeeking.value = true
        setTimeout(() => {
          isPurplePeeking.value = false
        }, 800)
      }, Math.random() * 3000 + 2000)
    }
    schedulePeek()
  } else {
    isPurplePeeking.value = false
    clearTimeout(peekInterval)
  }
})

const calculatePosition = (elRef) => {
  if (!elRef.value) return { faceX: 0, faceY: 0, bodySkew: 0 }
  const rect = elRef.value.getBoundingClientRect()
  const centerX = rect.left + rect.width / 2
  const centerY = rect.top + rect.height / 3
  const deltaX = mouseX.value - centerX
  const deltaY = mouseY.value - centerY

  const faceX = Math.max(-15, Math.min(15, deltaX / 20))
  const faceY = Math.max(-10, Math.min(10, deltaY / 30))
  const bodySkew = Math.max(-6, Math.min(6, -deltaX / 120))

  return { faceX, faceY, bodySkew }
}

const isHidingPassword = computed(() => props.passwordLength > 0 && !props.showPassword)

// Computed styles and positions
const purplePos = computed(() => calculatePosition(purpleRef))
const blackPos = computed(() => calculatePosition(blackRef))
const yellowPos = computed(() => calculatePosition(yellowRef))
const orangePos = computed(() => calculatePosition(orangeRef))

const purpleStyle = computed(() => {
  const isShowPwd = props.passwordLength > 0 && props.showPassword
  const isTypingOrHiding = props.isTyping || isHidingPassword.value
  let transform = `skewX(${purplePos.value.bodySkew || 0}deg)`
  
  if (isShowPwd) {
    transform = 'skewX(0deg)'
  } else if (isTypingOrHiding) {
    transform = `skewX(${(purplePos.value.bodySkew || 0) - 12}deg) translateX(40px)`
  }

  return {
    height: isTypingOrHiding ? '440px' : '400px',
    transform
  }
})

const purpleEyesStyle = computed(() => {
  const isShowPwd = props.passwordLength > 0 && props.showPassword
  let left = `${45 + purplePos.value.faceX}px`
  let top = `${40 + purplePos.value.faceY}px`
  
  if (isShowPwd) {
    left = '20px'
    top = '35px'
  } else if (isLookingAtEachOther.value) {
    left = '55px'
    top = '65px'
  }
  return { left, top }
})

const purpleForceLookX = computed(() => {
  if (props.passwordLength > 0 && props.showPassword) return isPurplePeeking.value ? 4 : -4
  if (isLookingAtEachOther.value) return 3
  return undefined
})

const purpleForceLookY = computed(() => {
  if (props.passwordLength > 0 && props.showPassword) return isPurplePeeking.value ? 5 : -4
  if (isLookingAtEachOther.value) return 4
  return undefined
})

const blackStyle = computed(() => {
  const isShowPwd = props.passwordLength > 0 && props.showPassword
  const isTypingOrHiding = props.isTyping || isHidingPassword.value
  let transform = `skewX(${blackPos.value.bodySkew || 0}deg)`

  if (isShowPwd) {
    transform = 'skewX(0deg)'
  } else if (isLookingAtEachOther.value) {
    transform = `skewX(${(blackPos.value.bodySkew || 0) * 1.5 + 10}deg) translateX(20px)`
  } else if (isTypingOrHiding) {
    transform = `skewX(${(blackPos.value.bodySkew || 0) * 1.5}deg)`
  }
  return { transform }
})

const blackEyesStyle = computed(() => {
  const isShowPwd = props.passwordLength > 0 && props.showPassword
  let left = `${26 + blackPos.value.faceX}px`
  let top = `${32 + blackPos.value.faceY}px`

  if (isShowPwd) {
    left = '10px'
    top = '28px'
  } else if (isLookingAtEachOther.value) {
    left = '32px'
    top = '12px'
  }
  return { left, top }
})

const blackForceLookX = computed(() => {
  if (props.passwordLength > 0 && props.showPassword) return -4
  if (isLookingAtEachOther.value) return 0
  return undefined
})

const blackForceLookY = computed(() => {
  if (props.passwordLength > 0 && props.showPassword) return -4
  if (isLookingAtEachOther.value) return -4
  return undefined
})

const orangeStyle = computed(() => {
  const isShowPwd = props.passwordLength > 0 && props.showPassword
  const transform = isShowPwd ? 'skewX(0deg)' : `skewX(${orangePos.value.bodySkew || 0}deg)`
  return { transform }
})

const orangeEyesStyle = computed(() => {
  const isShowPwd = props.passwordLength > 0 && props.showPassword
  const left = isShowPwd ? '50px' : `${82 + (orangePos.value.faceX || 0)}px`
  const top = isShowPwd ? '85px' : `${90 + (orangePos.value.faceY || 0)}px`
  return { left, top }
})

const orangeForceLookX = computed(() => {
  if (props.passwordLength > 0 && props.showPassword) return -5
  return undefined
})

const orangeForceLookY = computed(() => {
  if (props.passwordLength > 0 && props.showPassword) return -4
  return undefined
})

const yellowStyle = computed(() => {
  const isShowPwd = props.passwordLength > 0 && props.showPassword
  const transform = isShowPwd ? 'skewX(0deg)' : `skewX(${yellowPos.value.bodySkew || 0}deg)`
  return { transform }
})

const yellowEyesStyle = computed(() => {
  const isShowPwd = props.passwordLength > 0 && props.showPassword
  const left = isShowPwd ? '20px' : `${52 + (yellowPos.value.faceX || 0)}px`
  const top = isShowPwd ? '35px' : `${40 + (yellowPos.value.faceY || 0)}px`
  return { left, top }
})

const yellowForceLookX = computed(() => {
  if (props.passwordLength > 0 && props.showPassword) return -5
  return undefined
})

const yellowForceLookY = computed(() => {
  if (props.passwordLength > 0 && props.showPassword) return -4
  return undefined
})

const yellowMouthStyle = computed(() => {
  const isShowPwd = props.passwordLength > 0 && props.showPassword
  const left = isShowPwd ? '10px' : `${40 + (yellowPos.value.faceX || 0)}px`
  const top = isShowPwd ? '88px' : `${88 + (yellowPos.value.faceY || 0)}px`
  return { left, top }
})

</script>

<style scoped>
.characters-wrapper {
  position: relative;
  width: 550px;
  height: 400px;
  transform: scale(0.9); /* scale down slightly to fit well in standard screens */
}

.character {
  position: absolute;
  bottom: 0;
  transition: all 0.7s ease-in-out;
  transform-origin: bottom center;
}

.eyes-container {
  position: absolute;
  display: flex;
  gap: 32px; /* fallback gap */
  transition: all 0.7s ease-in-out;
}

.duration-200 {
  transition-duration: 0.2s;
  transition-timing-function: ease-out;
}

.gap-6 { gap: 24px; }
.gap-8 { gap: 32px; }

.purple-char {
  left: 70px;
  width: 180px;
  background-color: #6C3FF5;
  border-radius: 10px 10px 0 0;
  z-index: 1;
}

.black-char {
  left: 240px;
  width: 120px;
  height: 310px;
  background-color: #2D2D2D;
  border-radius: 8px 8px 0 0;
  z-index: 2;
}

.orange-char {
  left: 0px;
  width: 240px;
  height: 200px;
  background-color: #FF9B6B;
  border-radius: 120px 120px 0 0;
  z-index: 3;
}

.yellow-char {
  left: 310px;
  width: 140px;
  height: 230px;
  background-color: #E8D754;
  border-radius: 70px 70px 0 0;
  z-index: 4;
}

.mouth {
  position: absolute;
  width: 80px;
  height: 4px;
  background-color: #2D2D2D;
  border-radius: 4px;
  transition: all 0.2s ease-out;
}
</style>
