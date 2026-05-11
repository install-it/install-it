import { defineStore } from 'pinia'

export default defineStore('unsavedForm', () => {
  const show = ref(false)

  let answerHandler: ((allow: boolean) => void) | null = null

  return {
    show,
    confirmLeave: (answer: boolean) => {
      show.value = false
      if (answerHandler) {
        answerHandler(answer)
        answerHandler = null
      }
    },
    setAnswerHandler: (handler: (allow: boolean) => void) => {
      answerHandler = handler
    }
  }
})
