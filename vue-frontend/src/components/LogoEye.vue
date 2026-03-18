<template>
  <div class="logo-circle" :style="{ width: size + 'px', height: size + 'px' }">
    <div class="logo-eye" :style="{ width: (size * 0.5) + 'px', height: (size * 0.5) + 'px' }">
      <div class="logo-pupil" :style="{ width: (size * 0.2) + 'px', height: (size * 0.2) + 'px', top: (size * 0.15) + 'px', left: (size * 0.15) + 'px' }"></div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LogoEye',
  props: {
    size: {
      type: Number,
      default: 40
    }
  }
}
</script>

<style scoped>
.logo-circle {
  background: #111827;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

.logo-eye {
  background: white;
  border-radius: 50%;
  position: relative;
}

.logo-pupil {
  background: #111827;
  border-radius: 50%;
  position: absolute;
  animation: lookAround 4s infinite;
}

@keyframes lookAround {
  0%, 100% { transform: translate(0, 0); }
  25% { transform: translate(3px, 0); }
  50% { transform: translate(0, 3px); }
  75% { transform: translate(-3px, -2px); }
}
</style>
