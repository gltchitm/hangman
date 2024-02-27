import terser from '@rollup/plugin-terser'
import commonjs from '@rollup/plugin-commonjs'
import resolve from '@rollup/plugin-node-resolve'
import typescript from '@rollup/plugin-typescript'

import css from 'rollup-plugin-css-only'
import svelte from 'rollup-plugin-svelte'

import sveltePreprocess from 'svelte-preprocess'

const production = process.env.MIX_ENV === 'prod'

export default {
    input: 'src/main.ts',
    output: {
        sourcemap: false,
        format: 'iife',
        name: 'hangman',
        file: '../static/app.js'
    },
    plugins: [
        svelte({
            preprocess: sveltePreprocess({ sourceMap: false }),
            compilerOptions: { dev: !production }
        }),
        css({ output: 'app.css' }),
        resolve({ browser: true, dedupe: ['svelte'] }),
        commonjs(),
        typescript({ sourceMap: false }),
        production && terser()
    ],
    watch: {
        clearScreen: false
    }
}
