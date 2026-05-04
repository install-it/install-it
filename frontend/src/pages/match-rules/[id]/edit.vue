<script setup lang="ts">
import MatchRuleFormComponent from '@/components/MatchRuleFormComponent.vue'
import { useRoute } from 'vue-router'

const $route = useRoute()

const params = $route.params as { id: string }

const id = Array.isArray(params.id) ? params.id[0] : params.id
</script>

<template>
  <MatchRuleFormComponent :id="id" />
</template>
