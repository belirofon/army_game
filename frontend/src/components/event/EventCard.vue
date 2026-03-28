<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  description: string
  templateId?: string
  location?: string
  showCharacter?: boolean
  characterImage?: string
}>()

const backgroundImage = computed(() => {
  if (!props.location) return null
  
  const locationMap: Record<string, string> = {
    'DINING_ROOM': '/dining_room.jpg',
    'BARRACKS': '/barracks.png',
  }
  
  return locationMap[props.location] || '/start_screen_main.png'
})

const characterImageUrl = computed(() => {
  if (!props.showCharacter || !props.characterImage) return null
  return `/${props.characterImage}`
})
</script>

<template>
  <div class="event-card">
    <div 
      class="event-card__illustration"
      :style="backgroundImage ? { backgroundImage: `url(${backgroundImage})`, backgroundSize: 'cover', backgroundPosition: 'center' } : {}"
    >
      <div v-if="characterImageUrl" class="event-card__character">
        <img :src="characterImageUrl" alt="Character" />
      </div>
      <div v-else-if="!backgroundImage" class="event-card__placeholder">
        <svg class="icon-lg" viewBox="0 0 24 24" fill="currentColor">
          <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
        </svg>
      </div>
    </div>
    <p class="event-card__description">{{ description }}</p>
  </div>
</template>

<style scoped>
.event-card {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.event-card__illustration {
  width: 100%;
  aspect-ratio: 16 / 9;
  background-color: #2D3527;
  border-radius: 0.5rem;
  overflow: hidden;
  position: relative;
}

.event-card__character {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  max-width: 70%;
  max-height: 80%;
}

.event-card__character img {
  width: auto;
  height: auto;
  max-width: 100%;
  max-height: 120px;
  object-fit: contain;
  filter: drop-shadow(0 4px 8px rgba(0,0,0,0.5));
}

.event-card__placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6B665E;
}

.icon-lg {
  width: 4rem;
  height: 4rem;
  opacity: 0.5;
}

.event-card__description {
  font-family: 'Roboto', sans-serif;
  font-size: 1rem;
  color: #E8E4DC;
  line-height: 1.625;
}
</style>