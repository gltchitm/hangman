import { readable, writable } from 'svelte/store'

const gameType = writable(null)
const socket = readable(new WebSocket('ws://' + window.location.host + '/ws'), () => {})

export { gameType, socket }
